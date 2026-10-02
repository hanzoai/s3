package hanzofs.spark;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.FSDataOutputStream;
import org.apache.hadoop.fs.FileSystem;
import org.apache.hadoop.fs.Path;
import org.apache.parquet.example.data.Group;
import org.apache.parquet.example.data.simple.SimpleGroupFactory;
import org.apache.parquet.hadoop.ParquetFileWriter;
import org.apache.parquet.hadoop.ParquetWriter;
import org.apache.parquet.hadoop.example.GroupWriteSupport;
import org.apache.parquet.hadoop.metadata.CompressionCodecName;
import org.apache.parquet.schema.MessageType;
import org.apache.parquet.schema.MessageTypeParser;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;

import java.io.IOException;
import java.io.OutputStream;
import java.net.URI;
import java.util.ArrayList;
import java.util.List;

import static org.junit.Assert.*;

/**
 * Compare OutputStream behavior between local disk and Hanzo S3
 * to understand why Parquet files written to Hanzo S3 have incorrect metadata.
 */
public class OutputStreamComparisonTest extends SparkTestBase {

    private static class WriteOperation {
        String source;
        String operation;
        long positionBefore;
        long positionAfter;
        int bytesWritten;
        long timestamp;
        String details;

        WriteOperation(String source, String operation, long positionBefore, long positionAfter, 
                      int bytesWritten, String details) {
            this.source = source;
            this.operation = operation;
            this.positionBefore = positionBefore;
            this.positionAfter = positionAfter;
            this.bytesWritten = bytesWritten;
            this.timestamp = System.nanoTime();
            this.details = details;
        }

        @Override
        public String toString() {
            return String.format("[%s] %s: posBefore=%d, posAfter=%d, written=%d %s",
                    source, operation, positionBefore, positionAfter, bytesWritten, 
                    details != null ? "(" + details + ")" : "");
        }
    }

    private static class LoggingOutputStream extends OutputStream {
        private final FSDataOutputStream wrapped;
        private final String source;
        private final List<WriteOperation> operations;

        LoggingOutputStream(FSDataOutputStream wrapped, String source, List<WriteOperation> operations) {
            this.wrapped = wrapped;
            this.source = source;
            this.operations = operations;
        }

        @Override
        public void write(int b) throws IOException {
            long posBefore = wrapped.getPos();
            wrapped.write(b);
            long posAfter = wrapped.getPos();
            operations.add(new WriteOperation(source, "write(int)", posBefore, posAfter, 1, null));
        }

        @Override
        public void write(byte[] b, int off, int len) throws IOException {
            long posBefore = wrapped.getPos();
            wrapped.write(b, off, len);
            long posAfter = wrapped.getPos();
            operations.add(new WriteOperation(source, "write(byte[])", posBefore, posAfter, len, 
                    "len=" + len));
        }

        @Override
        public void flush() throws IOException {
            long posBefore = wrapped.getPos();
            wrapped.flush();
            long posAfter = wrapped.getPos();
            operations.add(new WriteOperation(source, "flush()", posBefore, posAfter, 0, null));
        }

        @Override
        public void close() throws IOException {
            long posBefore = wrapped.getPos();
            wrapped.close();
            long posAfter = 0; // Can't call getPos() after close
            operations.add(new WriteOperation(source, "close()", posBefore, posAfter, 0, 
                    "finalPos=" + posBefore));
        }

        public long getPos() throws IOException {
            long pos = wrapped.getPos();
            operations.add(new WriteOperation(source, "getPos()", pos, pos, 0, "returned=" + pos));
            return pos;
        }

        public void hflush() throws IOException {
            long posBefore = wrapped.getPos();
            wrapped.hflush();
            long posAfter = wrapped.getPos();
            operations.add(new WriteOperation(source, "hflush()", posBefore, posAfter, 0, null));
        }

        public void hsync() throws IOException {
            long posBefore = wrapped.getPos();
            wrapped.hsync();
            long posAfter = wrapped.getPos();
            operations.add(new WriteOperation(source, "hsync()", posBefore, posAfter, 0, null));
        }
    }

    private static final MessageType SCHEMA = MessageTypeParser.parseMessageType(
            "message schema {"
                    + "required int32 id;"
                    + "required binary name;"
                    + "required int32 age;"
                    + "}"
    );

    @Before
    public void setUp() throws IOException {
        if (!TESTS_ENABLED) {
            return;
        }
        super.setUpSpark();
    }

    @After
    public void tearDown() throws IOException {
        if (!TESTS_ENABLED) {
            return;
        }
        super.tearDownSpark();
    }

    @Test
    public void testCompareOutputStreamBehavior() throws Exception {
        skipIfTestsDisabled();

        System.out.println("\n╔══════════════════════════════════════════════════════════════╗");
        System.out.println("║  REAL-TIME OUTPUTSTREAM COMPARISON: LOCAL vs HANZO S3      ║");
        System.out.println("╚══════════════════════════════════════════════════════════════╝");

        // Prepare file systems
        Configuration conf = new Configuration();
        FileSystem localFs = FileSystem.getLocal(conf);
        
        conf.set("fs.hanzofs.impl", "hanzofs.hdfs.FilerFileSystem");
        conf.set("fs.hanzofs.filer.host", S3_HOST);
        conf.set("fs.hanzofs.filer.port", String.valueOf(S3_PORT));
        FileSystem filerFs = FileSystem.get(URI.create(String.format("hanzofs://%s:%s", 
                S3_HOST, S3_PORT)), conf);

        // Prepare paths
        new java.io.File("/workspace/target/test-output").mkdirs();
        Path localPath = new Path("file:///workspace/target/test-output/write-comparison-local.parquet");
        Path filerPath = new Path(getTestPath("write-comparison-s3.parquet"));

        // Delete if exists
        localFs.delete(localPath, false);
        filerFs.delete(filerPath, false);

        List<WriteOperation> localOps = new ArrayList<>();
        List<WriteOperation> filerOps = new ArrayList<>();

        System.out.println("\n1. Writing Parquet files with synchronized operations...\n");

        // Write using ParquetWriter with custom OutputStreams
        GroupWriteSupport.setSchema(SCHEMA, conf);
        
        // Create data
        SimpleGroupFactory groupFactory = new SimpleGroupFactory(SCHEMA);
        List<Group> groups = new ArrayList<>();
        groups.add(groupFactory.newGroup().append("id", 1).append("name", "Alice").append("age", 30));
        groups.add(groupFactory.newGroup().append("id", 2).append("name", "Bob").append("age", 25));
        groups.add(groupFactory.newGroup().append("id", 3).append("name", "Charlie").append("age", 35));

        // Write to local disk
        System.out.println("   Writing to LOCAL DISK...");
        try (ParquetWriter<Group> localWriter = new ParquetWriter<>(
                localPath,
                new GroupWriteSupport(),
                CompressionCodecName.SNAPPY,
                1024 * 1024, // Block size
                1024, // Page size
                1024, // Dictionary page size
                true, // Enable dictionary
                false, // Don't validate
                ParquetWriter.DEFAULT_WRITER_VERSION,
                conf)) {
            for (Group group : groups) {
                localWriter.write(group);
            }
        }
        System.out.println("   ✅ Local write complete");

        // Write to Hanzo S3
        System.out.println("\n   Writing to HANZO S3...");
        try (ParquetWriter<Group> filerWriter = new ParquetWriter<>(
                filerPath,
                new GroupWriteSupport(),
                CompressionCodecName.SNAPPY,
                1024 * 1024, // Block size
                1024, // Page size
                1024, // Dictionary page size
                true, // Enable dictionary
                false, // Don't validate
                ParquetWriter.DEFAULT_WRITER_VERSION,
                conf)) {
            for (Group group : groups) {
                filerWriter.write(group);
            }
        }
        System.out.println("   ✅ Hanzo S3 write complete");

        // Compare file sizes
        System.out.println("\n2. Comparing final file sizes...");
        long localSize = localFs.getFileStatus(localPath).getLen();
        long filerSize = filerFs.getFileStatus(filerPath).getLen();
        System.out.println("   LOCAL:    " + localSize + " bytes");
        System.out.println("   S3:  " + filerSize + " bytes");

        if (localSize == filerSize) {
            System.out.println("   ✅ File sizes MATCH");
        } else {
            System.out.println("   ❌ File sizes DIFFER by " + Math.abs(localSize - filerSize) + " bytes");
        }

        // Now test reading both files
        System.out.println("\n3. Testing if both files can be read by Spark...");
        
        System.out.println("\n   Reading LOCAL file:");
        try {
            org.apache.spark.sql.Dataset<org.apache.spark.sql.Row> localDf = 
                    spark.read().parquet(localPath.toString());
            long localCount = localDf.count();
            System.out.println("   ✅ LOCAL read SUCCESS - " + localCount + " rows");
            localDf.show();
        } catch (Exception e) {
            System.out.println("   ❌ LOCAL read FAILED: " + e.getMessage());
            e.printStackTrace();
        }

        System.out.println("\n   Reading HANZO S3 file:");
        try {
            org.apache.spark.sql.Dataset<org.apache.spark.sql.Row> filerDf = 
                    spark.read().parquet(filerPath.toString());
            long filerCount = filerDf.count();
            System.out.println("   ✅ HANZO S3 read SUCCESS - " + filerCount + " rows");
            filerDf.show();
        } catch (Exception e) {
            System.out.println("   ❌ HANZO S3 read FAILED: " + e.getMessage());
            e.printStackTrace();
        }

        System.out.println("\n╔══════════════════════════════════════════════════════════════╗");
        System.out.println("║  COMPARISON COMPLETE                                         ║");
        System.out.println("╚══════════════════════════════════════════════════════════════╝");
    }

    @Test
    public void testCompareRawOutputStreamOperations() throws Exception {
        skipIfTestsDisabled();

        System.out.println("\n╔══════════════════════════════════════════════════════════════╗");
        System.out.println("║  RAW OUTPUTSTREAM COMPARISON: LOCAL vs HANZO S3            ║");
        System.out.println("╚══════════════════════════════════════════════════════════════╝");

        // Prepare file systems
        Configuration conf = new Configuration();
        FileSystem localFs = FileSystem.getLocal(conf);
        
        conf.set("fs.hanzofs.impl", "hanzofs.hdfs.FilerFileSystem");
        conf.set("fs.hanzofs.filer.host", S3_HOST);
        conf.set("fs.hanzofs.filer.port", String.valueOf(S3_PORT));
        FileSystem filerFs = FileSystem.get(URI.create(String.format("hanzofs://%s:%s", 
                S3_HOST, S3_PORT)), conf);

        // Prepare paths
        new java.io.File("/workspace/target/test-output").mkdirs();
        Path localPath = new Path("file:///workspace/target/test-output/raw-comparison-local.dat");
        Path filerPath = new Path(getTestPath("raw-comparison-s3.dat"));

        // Delete if exists
        localFs.delete(localPath, false);
        filerFs.delete(filerPath, false);

        List<WriteOperation> localOps = new ArrayList<>();
        List<WriteOperation> filerOps = new ArrayList<>();

        System.out.println("\n1. Performing synchronized write operations...\n");

        // Open both streams
        FSDataOutputStream localStream = localFs.create(localPath, true);
        FSDataOutputStream filerStream = filerFs.create(filerPath, true);

        LoggingOutputStream localLogging = new LoggingOutputStream(localStream, "LOCAL", localOps);
        LoggingOutputStream filerLogging = new LoggingOutputStream(filerStream, "S3", filerOps);

        int opCount = 0;
        boolean mismatchFound = false;

        // Operation 1: Write 4 bytes (magic)
        opCount++;
        System.out.println("   Op " + opCount + ": write(4 bytes) - Writing magic bytes");
        byte[] magic = "PAR1".getBytes();
        localLogging.write(magic, 0, 4);
        filerLogging.write(magic, 0, 4);
        long localPos1 = localLogging.getPos();
        long filerPos1 = filerLogging.getPos();
        System.out.println("      LOCAL:   getPos() = " + localPos1);
        System.out.println("      S3: getPos() = " + filerPos1);
        if (localPos1 != filerPos1) {
            System.out.println("      ❌ MISMATCH!");
            mismatchFound = true;
        } else {
            System.out.println("      ✅ Match");
        }

        // Operation 2: Write 100 bytes of data
        opCount++;
        System.out.println("\n   Op " + opCount + ": write(100 bytes) - Writing data");
        byte[] data = new byte[100];
        for (int i = 0; i < 100; i++) {
            data[i] = (byte) i;
        }
        localLogging.write(data, 0, 100);
        filerLogging.write(data, 0, 100);
        long localPos2 = localLogging.getPos();
        long filerPos2 = filerLogging.getPos();
        System.out.println("      LOCAL:   getPos() = " + localPos2);
        System.out.println("      S3: getPos() = " + filerPos2);
        if (localPos2 != filerPos2) {
            System.out.println("      ❌ MISMATCH!");
            mismatchFound = true;
        } else {
            System.out.println("      ✅ Match");
        }

        // Operation 3: Flush
        opCount++;
        System.out.println("\n   Op " + opCount + ": flush()");
        localLogging.flush();
        filerLogging.flush();
        long localPos3 = localLogging.getPos();
        long filerPos3 = filerLogging.getPos();
        System.out.println("      LOCAL:   getPos() after flush = " + localPos3);
        System.out.println("      S3: getPos() after flush = " + filerPos3);
        if (localPos3 != filerPos3) {
            System.out.println("      ❌ MISMATCH!");
            mismatchFound = true;
        } else {
            System.out.println("      ✅ Match");
        }

        // Operation 4: Write more data
        opCount++;
        System.out.println("\n   Op " + opCount + ": write(50 bytes) - Writing more data");
        byte[] moreData = new byte[50];
        for (int i = 0; i < 50; i++) {
            moreData[i] = (byte) (i + 100);
        }
        localLogging.write(moreData, 0, 50);
        filerLogging.write(moreData, 0, 50);
        long localPos4 = localLogging.getPos();
        long filerPos4 = filerLogging.getPos();
        System.out.println("      LOCAL:   getPos() = " + localPos4);
        System.out.println("      S3: getPos() = " + filerPos4);
        if (localPos4 != filerPos4) {
            System.out.println("      ❌ MISMATCH!");
            mismatchFound = true;
        } else {
            System.out.println("      ✅ Match");
        }

        // Operation 5: Write final bytes (simulating footer)
        opCount++;
        System.out.println("\n   Op " + opCount + ": write(8 bytes) - Writing footer");
        byte[] footer = new byte[]{0x6B, 0x03, 0x00, 0x00, 0x50, 0x41, 0x52, 0x31};
        localLogging.write(footer, 0, 8);
        filerLogging.write(footer, 0, 8);
        long localPos5 = localLogging.getPos();
        long filerPos5 = filerLogging.getPos();
        System.out.println("      LOCAL:   getPos() = " + localPos5);
        System.out.println("      S3: getPos() = " + filerPos5);
        if (localPos5 != filerPos5) {
            System.out.println("      ❌ MISMATCH!");
            mismatchFound = true;
        } else {
            System.out.println("      ✅ Match");
        }

        // Operation 6: Close
        opCount++;
        System.out.println("\n   Op " + opCount + ": close()");
        System.out.println("      LOCAL:   closing at position " + localPos5);
        System.out.println("      S3: closing at position " + filerPos5);
        localLogging.close();
        filerLogging.close();

        // Check final file sizes
        System.out.println("\n2. Comparing final file sizes...");
        long localSize = localFs.getFileStatus(localPath).getLen();
        long filerSize = filerFs.getFileStatus(filerPath).getLen();
        System.out.println("   LOCAL:    " + localSize + " bytes");
        System.out.println("   S3:  " + filerSize + " bytes");
        
        if (localSize != filerSize) {
            System.out.println("   ❌ File sizes DIFFER by " + Math.abs(localSize - filerSize) + " bytes");
            mismatchFound = true;
        } else {
            System.out.println("   ✅ File sizes MATCH");
        }

        System.out.println("\n╔══════════════════════════════════════════════════════════════╗");
        System.out.println("║  COMPARISON SUMMARY                                          ║");
        System.out.println("╚══════════════════════════════════════════════════════════════╝");
        System.out.println("   Total operations: " + opCount);
        System.out.println("   LOCAL operations:  " + localOps.size());
        System.out.println("   S3 operations: " + filerOps.size());
        
        if (mismatchFound) {
            System.out.println("\n   ❌ MISMATCHES FOUND - Streams behave differently!");
        } else {
            System.out.println("\n   ✅ ALL OPERATIONS MATCH - Streams are identical!");
        }

        System.out.println("\n   Detailed operation log:");
        System.out.println("   ----------------------");
        int maxOps = Math.max(localOps.size(), filerOps.size());
        for (int i = 0; i < maxOps; i++) {
            if (i < localOps.size()) {
                System.out.println("   " + localOps.get(i));
            }
            if (i < filerOps.size()) {
                System.out.println("   " + filerOps.get(i));
            }
            if (i < localOps.size() && i < filerOps.size()) {
                WriteOperation localOp = localOps.get(i);
                WriteOperation filerOp = filerOps.get(i);
                if (localOp.positionAfter != filerOp.positionAfter) {
                    System.out.println("      ⚠️  Position mismatch: LOCAL=" + localOp.positionAfter + 
                            " S3=" + filerOp.positionAfter);
                }
            }
        }

        assertFalse("Streams should behave identically", mismatchFound);
    }
}

