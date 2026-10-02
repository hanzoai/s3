package hanzofs.spark;

import org.apache.spark.sql.Dataset;
import org.apache.spark.sql.Row;
import org.apache.spark.sql.SaveMode;
import org.apache.spark.sql.SparkSession;

/**
 * Example Spark application demonstrating Hanzo S3 integration.
 * 
 * This can be submitted to a Spark cluster using spark-submit.
 * 
 * Example usage:
 * spark-submit \
 * --class hanzofs.spark.SparkFilerExample \
 * --master local[2] \
 * --conf spark.hadoop.fs.hanzofs.impl=hanzofs.hdfs.FilerFileSystem \
 * --conf spark.hadoop.fs.hanzofs.filer.host=localhost \
 * --conf spark.hadoop.fs.hanzofs.filer.port=8888 \
 * --conf spark.hadoop.fs.hanzofs.filer.port.grpc=18888 \
 * target/s3-spark-integration-tests-1.0-SNAPSHOT.jar \
 * hanzofs://localhost:8888/output
 */
public class SparkFilerExample {

        public static void main(String[] args) {
                if (args.length < 1) {
                        System.err.println("Usage: SparkFilerExample <output-path>");
                        System.err.println("Example: hanzofs://localhost:8888/spark-output");
                        System.exit(1);
                }

                String outputPath = args[0];

                // Create Spark session
                SparkSession spark = SparkSession.builder()
                                .appName("Hanzo S3 Spark Example")
                                .getOrCreate();

                try {
                        System.out.println("=== Hanzo S3 Spark Integration Example ===\n");

                        // Example 1: Generate data and write to Hanzo S3
                        System.out.println("1. Generating sample data...");
                        Dataset<Row> data = spark.range(0, 1000)
                                        .selectExpr(
                                                        "id",
                                                        "id * 2 as doubled",
                                                        "CAST(rand() * 100 AS INT) as random_value");

                        System.out.println("   Generated " + data.count() + " rows");
                        data.show(5);

                        // Write as Parquet
                        String parquetPath = outputPath + "/data.parquet";
                        System.out.println("\n2. Writing data to Hanzo S3 as Parquet...");
                        System.out.println("   Path: " + parquetPath);

                        data.write()
                                        .mode(SaveMode.Overwrite)
                                        .parquet(parquetPath);

                        System.out.println("   ✓ Write completed");

                        // Read back and verify
                        System.out.println("\n3. Reading data back from Hanzo S3...");
                        Dataset<Row> readData = spark.read().parquet(parquetPath);
                        System.out.println("   Read " + readData.count() + " rows");

                        // Perform aggregation
                        System.out.println("\n4. Performing aggregation...");
                        Dataset<Row> stats = readData.selectExpr(
                                        "COUNT(*) as count",
                                        "AVG(random_value) as avg_random",
                                        "MAX(doubled) as max_doubled");

                        stats.show();

                        // Write aggregation results
                        String statsPath = outputPath + "/stats.parquet";
                        System.out.println("5. Writing stats to: " + statsPath);
                        stats.write()
                                        .mode(SaveMode.Overwrite)
                                        .parquet(statsPath);

                        // Create a partitioned dataset
                        System.out.println("\n6. Creating partitioned dataset...");
                        Dataset<Row> partitionedData = data.selectExpr(
                                        "*",
                                        "CAST(id % 10 AS INT) as partition_key");

                        String partitionedPath = outputPath + "/partitioned.parquet";
                        System.out.println("   Path: " + partitionedPath);

                        partitionedData.write()
                                        .mode(SaveMode.Overwrite)
                                        .partitionBy("partition_key")
                                        .parquet(partitionedPath);

                        System.out.println("   ✓ Partitioned write completed");

                        // Read specific partition
                        System.out.println("\n7. Reading specific partition (partition_key=0)...");
                        Dataset<Row> partition0 = spark.read()
                                        .parquet(partitionedPath)
                                        .filter("partition_key = 0");

                        System.out.println("   Partition 0 contains " + partition0.count() + " rows");
                        partition0.show(5);

                        // SQL example
                        System.out.println("\n8. Using Spark SQL...");
                        readData.createOrReplaceTempView("s3_data");

                        Dataset<Row> sqlResult = spark.sql(
                                        "SELECT " +
                                                        "  CAST(id / 100 AS INT) as bucket, " +
                                                        "  COUNT(*) as count, " +
                                                        "  AVG(random_value) as avg_random " +
                                                        "FROM s3_data " +
                                                        "GROUP BY CAST(id / 100 AS INT) " +
                                                        "ORDER BY bucket");

                        System.out.println("   Bucketed statistics:");
                        sqlResult.show();

                        System.out.println("\n=== Example completed successfully! ===");
                        System.out.println("Output location: " + outputPath);

                } catch (Exception e) {
                        System.err.println("Error: " + e.getMessage());
                        e.printStackTrace();
                        System.exit(1);
                } finally {
                        spark.stop();
                }
        }
}
