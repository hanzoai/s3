package hanzofs.spark;

import org.apache.hadoop.conf.Configuration;
import org.apache.spark.SparkConf;
import org.apache.spark.sql.SparkSession;
import org.junit.After;
import org.junit.Before;

import java.io.IOException;

/**
 * Base class for Spark integration tests with Hanzo S3.
 * 
 * These tests require a running Hanzo S3 cluster.
 * Set environment variable S3_TEST_ENABLED=true to enable these tests.
 */
public abstract class SparkTestBase {

    protected SparkSession spark;
    protected static final String TEST_ROOT = "/test-spark";
    protected static final boolean TESTS_ENABLED = "true".equalsIgnoreCase(System.getenv("S3_TEST_ENABLED"));

    // Hanzo S3 connection settings
    protected static final String S3_HOST = System.getenv().getOrDefault("S3_FILER_HOST", "localhost");
    protected static final String S3_PORT = System.getenv().getOrDefault("S3_FILER_PORT", "8888");
    protected static final String S3_GRPC_PORT = System.getenv().getOrDefault("S3_FILER_GRPC_PORT",
            "18888");

    @Before
    public void setUpSpark() throws IOException {
        if (!TESTS_ENABLED) {
            return;
        }

        SparkConf sparkConf = new SparkConf()
                .setAppName("Hanzo S3 Integration Test")
                .setMaster("local[1]") // Single thread to avoid concurrent gRPC issues
                .set("spark.driver.host", "localhost")
                .set("spark.sql.warehouse.dir", getS3Path("/spark-warehouse"))
                // Hanzo S3 configuration
                .set("spark.hadoop.fs.defaultFS", String.format("hanzofs://%s:%s", S3_HOST, S3_PORT))
                .set("spark.hadoop.fs.hanzofs.impl", "hanzofs.hdfs.FilerFileSystem")
                .set("spark.hadoop.fs.hanzofs.filer.host", S3_HOST)
                .set("spark.hadoop.fs.hanzofs.filer.port", S3_PORT)
                .set("spark.hadoop.fs.hanzofs.filer.port.grpc", S3_GRPC_PORT)
                .set("spark.hadoop.fs.AbstractFileSystem.hanzofs.impl", "hanzofs.hdfs.FilerAbstractFileSystem")
                // Set replication to empty string to use filer default
                .set("spark.hadoop.fs.hanzofs.replication", "")
                // Smaller buffer to reduce load
                .set("spark.hadoop.fs.hanzofs.buffer.size", "1048576") // 1MB
                // Reduce parallelism
                .set("spark.default.parallelism", "1")
                .set("spark.sql.shuffle.partitions", "1")
                // Simpler output committer
                .set("spark.hadoop.mapreduce.fileoutputcommitter.algorithm.version", "2")
                .set("spark.sql.sources.commitProtocolClass",
                        "org.apache.spark.sql.execution.datasources.SQLHadoopMapReduceCommitProtocol")
                // Disable speculative execution to reduce load
                .set("spark.speculation", "false")
                // Increase task retry to handle transient consistency issues
                .set("spark.task.maxFailures", "4")
                // Wait longer before retrying failed tasks
                .set("spark.task.reaper.enabled", "true")
                .set("spark.task.reaper.pollingInterval", "1s");

        spark = SparkSession.builder()
                .config(sparkConf)
                .getOrCreate();

        // Clean up test directory
        cleanupTestDirectory();
    }

    @After
    public void tearDownSpark() {
        if (!TESTS_ENABLED || spark == null) {
            return;
        }

        try {
            // Try to cleanup but don't fail if it doesn't work
            cleanupTestDirectory();
        } catch (Exception e) {
            System.err.println("Cleanup failed: " + e.getMessage());
        } finally {
            try {
                spark.stop();
            } catch (Exception e) {
                System.err.println("Spark stop failed: " + e.getMessage());
            }
            spark = null;
        }
    }

    protected String getS3Path(String path) {
        return String.format("hanzofs://%s:%s%s", S3_HOST, S3_PORT, path);
    }

    protected String getTestPath(String subPath) {
        return getS3Path(TEST_ROOT + "/" + subPath);
    }

    private void cleanupTestDirectory() {
        if (spark != null) {
            try {
                Configuration conf = spark.sparkContext().hadoopConfiguration();
                org.apache.hadoop.fs.FileSystem fs = org.apache.hadoop.fs.FileSystem.get(
                        java.net.URI.create(getS3Path("/")), conf);
                org.apache.hadoop.fs.Path testPath = new org.apache.hadoop.fs.Path(TEST_ROOT);
                if (fs.exists(testPath)) {
                    fs.delete(testPath, true);
                }
            } catch (Exception e) {
                // Suppress cleanup errors - they shouldn't fail tests
                // Common in distributed systems with eventual consistency
                System.err.println("Warning: cleanup failed (non-critical): " + e.getMessage());
            }
        }
    }

    protected void skipIfTestsDisabled() {
        if (!TESTS_ENABLED) {
            System.out.println("Skipping test - S3_TEST_ENABLED not set");
            org.junit.Assume.assumeTrue("S3_TEST_ENABLED not set", false);
        }
    }
}
