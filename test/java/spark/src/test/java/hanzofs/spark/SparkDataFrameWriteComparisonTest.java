package hanzofs.spark;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.*;
import org.apache.hadoop.fs.permission.FsPermission;
import org.apache.hadoop.util.Progressable;
import org.apache.spark.sql.Dataset;
import org.apache.spark.sql.Row;
import org.apache.spark.sql.SaveMode;
import org.junit.Test;

import java.io.IOException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

/**
 * Compare Spark DataFrame.write().parquet() operations between
 * local filesystem and Hanzo S3 to identify the exact difference
 * that causes the 78-byte EOF error.
 */
public class SparkDataFrameWriteComparisonTest extends SparkTestBase {

    private static class OperationLog {
        List<String> operations = new ArrayList<>();
        
        synchronized void log(String op) {
            operations.add(op);
            System.out.println("  " + op);
        }
        
        void print(String title) {
            System.out.println("\n" + title + " (" + operations.size() + " operations):");
            for (int i = 0; i < operations.size(); i++) {
                System.out.printf("  [%3d] %s\n", i, operations.get(i));
            }
        }
        
        void compare(OperationLog other, String name1, String name2) {
            System.out.println("\n╔══════════════════════════════════════════════════════════════╗");
            System.out.println("║  COMPARISON: " + name1 + " vs " + name2);
            System.out.println("╚══════════════════════════════════════════════════════════════╝");
            
            int maxLen = Math.max(operations.size(), other.operations.size());
            int differences = 0;
            
            for (int i = 0; i < maxLen; i++) {
                String op1 = i < operations.size() ? operations.get(i) : "<missing>";
                String op2 = i < other.operations.size() ? other.operations.get(i) : "<missing>";
                
                // Normalize operation strings for comparison (remove file-specific parts)
                String normalized1 = normalizeOp(op1);
                String normalized2 = normalizeOp(op2);
                
                if (!normalized1.equals(normalized2)) {
                    differences++;
                    System.out.printf("\n[%3d] DIFFERENCE:\n", i);
                    System.out.println("  " + name1 + ": " + op1);
                    System.out.println("  " + name2 + ": " + op2);
                }
            }
            
            System.out.println("\n" + "=".repeat(64));
            if (differences == 0) {
                System.out.println("✅ Operations are IDENTICAL!");
            } else {
                System.out.println("❌ Found " + differences + " differences");
            }
            System.out.println("=".repeat(64));
        }
        
        private String normalizeOp(String op) {
            // Remove file-specific identifiers for comparison
            return op.replaceAll("part-[0-9a-f-]+", "part-XXXXX")
                     .replaceAll("attempt_[0-9]+", "attempt_XXXXX")
                     .replaceAll("/tmp/[^/]+", "/tmp/XXXXX")
                     .replaceAll("test-local-[0-9]+", "test-local-XXXXX");
        }
    }
    
    // Custom FileSystem wrapper that logs all operations
    private static class LoggingFileSystem extends FilterFileSystem {
        private final OperationLog log;
        private final String name;
        
        public LoggingFileSystem(FileSystem fs, OperationLog log, String name) {
            this.fs = fs;
            this.log = log;
            this.name = name;
        }
        
        @Override
        public FSDataOutputStream create(Path f, FsPermission permission, boolean overwrite,
                                        int bufferSize, short replication, long blockSize,
                                        Progressable progress) throws IOException {
            log.log(String.format("%s CREATE: %s (bufferSize=%d)", name, f.getName(), bufferSize));
            FSDataOutputStream out = fs.create(f, permission, overwrite, bufferSize, replication, blockSize, progress);
            return new LoggingOutputStream(out, log, name, f.getName());
        }
        
        @Override
        public FSDataOutputStream append(Path f, int bufferSize, Progressable progress) throws IOException {
            log.log(String.format("%s APPEND: %s (bufferSize=%d)", name, f.getName(), bufferSize));
            FSDataOutputStream out = fs.append(f, bufferSize, progress);
            return new LoggingOutputStream(out, log, name, f.getName());
        }
        
        @Override
        public boolean rename(Path src, Path dst) throws IOException {
            log.log(String.format("%s RENAME: %s → %s", name, src.getName(), dst.getName()));
            return fs.rename(src, dst);
        }
        
        @Override
        public boolean delete(Path f, boolean recursive) throws IOException {
            log.log(String.format("%s DELETE: %s (recursive=%s)", name, f.getName(), recursive));
            return fs.delete(f, recursive);
        }
        
        @Override
        public FileStatus[] listStatus(Path f) throws IOException {
            FileStatus[] result = fs.listStatus(f);
            log.log(String.format("%s LISTSTATUS: %s (%d files)", name, f.getName(), result.length));
            return result;
        }
        
        @Override
        public void setWorkingDirectory(Path new_dir) {
            fs.setWorkingDirectory(new_dir);
        }
        
        @Override
        public Path getWorkingDirectory() {
            return fs.getWorkingDirectory();
        }
        
        @Override
        public boolean mkdirs(Path f, FsPermission permission) throws IOException {
            log.log(String.format("%s MKDIRS: %s", name, f.getName()));
            return fs.mkdirs(f, permission);
        }
        
        @Override
        public FileStatus getFileStatus(Path f) throws IOException {
            FileStatus status = fs.getFileStatus(f);
            log.log(String.format("%s GETFILESTATUS: %s (size=%d)", name, f.getName(), status.getLen()));
            return status;
        }
        
        @Override
        public FSDataInputStream open(Path f, int bufferSize) throws IOException {
            log.log(String.format("%s OPEN: %s (bufferSize=%d)", name, f.getName(), bufferSize));
            return fs.open(f, bufferSize);
        }
    }
    
    private static class LoggingOutputStream extends FSDataOutputStream {
        private final FSDataOutputStream delegate;
        private final OperationLog log;
        private final String name;
        private final String filename;
        private long writeCount = 0;
        
        public LoggingOutputStream(FSDataOutputStream delegate, OperationLog log, String name, String filename) throws IOException {
            super(delegate.getWrappedStream(), null);
            this.delegate = delegate;
            this.log = log;
            this.name = name;
            this.filename = filename;
        }
        
        @Override
        public void write(byte[] b, int off, int len) throws IOException {
            writeCount++;
            long posBefore = getPos();
            delegate.write(b, off, len);
            long posAfter = getPos();
            
            // Log significant writes and the last few writes (potential footer)
            if (len >= 100 || writeCount <= 5 || (writeCount % 100 == 0)) {
                log.log(String.format("%s WRITE #%d: %d bytes, pos %d→%d [%s]", 
                                     name, writeCount, len, posBefore, posAfter, filename));
            }
        }
        
        @Override
        public long getPos() {
            long pos = delegate.getPos();
            return pos;
        }
        
        @Override
        public void flush() throws IOException {
            log.log(String.format("%s FLUSH: pos=%d [%s]", name, getPos(), filename));
            delegate.flush();
        }
        
        @Override
        public void close() throws IOException {
            log.log(String.format("%s CLOSE: pos=%d, totalWrites=%d [%s]", 
                                 name, getPos(), writeCount, filename));
            delegate.close();
        }
        
        @Override
        public void hflush() throws IOException {
            log.log(String.format("%s HFLUSH: pos=%d [%s]", name, getPos(), filename));
            delegate.hflush();
        }
        
        @Override
        public void hsync() throws IOException {
            log.log(String.format("%s HSYNC: pos=%d [%s]", name, getPos(), filename));
            delegate.hsync();
        }
    }
    
    @Test
    public void testCompareSparkDataFrameWrite() throws Exception {
        if (!TESTS_ENABLED) {
            System.out.println("Skipping test - S3_TEST_ENABLED not set");
            return;
        }
        
        System.out.println("\n╔══════════════════════════════════════════════════════════════╗");
        System.out.println("║  SPARK DATAFRAME.WRITE() OPERATION COMPARISON TEST           ║");
        System.out.println("╚══════════════════════════════════════════════════════════════╝\n");
        
        // Create test data (4 rows - this is what causes the error)
        List<Employee> employees = Arrays.asList(
            new Employee(1, "Alice", "Engineering", 100000),
            new Employee(2, "Bob", "Sales", 80000),
            new Employee(3, "Charlie", "Engineering", 120000),
            new Employee(4, "David", "Sales", 75000)
        );
        
        Dataset<Row> df = spark.createDataFrame(employees, Employee.class);
        
        OperationLog localLog = new OperationLog();
        OperationLog filerLog = new OperationLog();
        
        // Test 1: Write to local filesystem with logging
        System.out.println("=== Writing to LOCAL filesystem with Spark ===");
        String localPath = "/tmp/spark-local-test-" + System.currentTimeMillis();
        
        try {
            // Configure Spark to use our logging filesystem for local writes
            Configuration localConf = new Configuration();
            FileSystem localFs = FileSystem.getLocal(localConf);
            LoggingFileSystem loggingLocalFs = new LoggingFileSystem(localFs, localLog, "LOCAL");
            
            // Write using Spark
            df.write().mode(SaveMode.Overwrite).parquet("file://" + localPath);
            
            System.out.println("✅ Local write completed");
            
            // Check final file
            FileStatus[] files = localFs.listStatus(new Path(localPath));
            for (FileStatus file : files) {
                if (file.getPath().getName().endsWith(".parquet")) {
                    localLog.log(String.format("LOCAL FINAL FILE: %s (%d bytes)", 
                                             file.getPath().getName(), file.getLen()));
                }
            }
            
        } catch (Exception e) {
            System.out.println("❌ Local write failed: " + e.getMessage());
            e.printStackTrace();
        }
        
        // Test 2: Write to Hanzo S3 with logging
        System.out.println("\n=== Writing to HANZO S3 with Spark ===");
        String filerPath = getTestPath("spark-s3-test");
        
        try {
            df.write().mode(SaveMode.Overwrite).parquet(filerPath);
            
            System.out.println("✅ Hanzo S3 write completed");
            
            // Check final file
            Configuration filerConf = new Configuration();
            filerConf.set("fs.hanzofs.impl", "hanzofs.hdfs.FilerFileSystem");
            filerConf.set("fs.hanzofs.filer.host", S3_HOST);
            filerConf.set("fs.hanzofs.filer.port", S3_PORT);
            FileSystem filerFs = FileSystem.get(
                java.net.URI.create("hanzofs://" + S3_HOST + ":" + S3_PORT), 
                filerConf);
            
            FileStatus[] files = filerFs.listStatus(new Path(filerPath));
            for (FileStatus file : files) {
                if (file.getPath().getName().endsWith(".parquet")) {
                    filerLog.log(String.format("S3 FINAL FILE: %s (%d bytes)", 
                                                file.getPath().getName(), file.getLen()));
                }
            }
            
        } catch (Exception e) {
            System.out.println("❌ Hanzo S3 write failed: " + e.getMessage());
            if (e.getMessage() != null && e.getMessage().contains("bytes left")) {
                System.out.println("🎯 This is the 78-byte EOF error during WRITE!");
            }
            e.printStackTrace();
        }
        
        // Test 3: Try reading both
        System.out.println("\n=== Reading LOCAL file ===");
        try {
            Dataset<Row> localDf = spark.read().parquet("file://" + localPath);
            long count = localDf.count();
            System.out.println("✅ Local read successful: " + count + " rows");
        } catch (Exception e) {
            System.out.println("❌ Local read failed: " + e.getMessage());
        }
        
        System.out.println("\n=== Reading HANZO S3 file ===");
        try {
            Dataset<Row> filerDf = spark.read().parquet(filerPath);
            long count = filerDf.count();
            System.out.println("✅ Hanzo S3 read successful: " + count + " rows");
        } catch (Exception e) {
            System.out.println("❌ Hanzo S3 read failed: " + e.getMessage());
            if (e.getMessage() != null && e.getMessage().contains("bytes left")) {
                System.out.println("🎯 This is the 78-byte EOF error during READ!");
            }
        }
        
        // Print operation logs
        localLog.print("LOCAL OPERATIONS");
        filerLog.print("HANZO S3 OPERATIONS");
        
        // Compare
        localLog.compare(filerLog, "LOCAL", "HANZO S3");
        
        System.out.println("\n=== Test Complete ===");
    }
    
    // Employee class for test data
    public static class Employee implements java.io.Serializable {
        private int id;
        private String name;
        private String department;
        private int salary;

        public Employee() {}

        public Employee(int id, String name, String department, int salary) {
            this.id = id;
            this.name = name;
            this.department = department;
            this.salary = salary;
        }

        public int getId() { return id; }
        public void setId(int id) { this.id = id; }
        public String getName() { return name; }
        public void setName(String name) { this.name = name; }
        public String getDepartment() { return department; }
        public void setDepartment(String department) { this.department = department; }
        public int getSalary() { return salary; }
        public void setSalary(int salary) { this.salary = salary; }
    }
}

