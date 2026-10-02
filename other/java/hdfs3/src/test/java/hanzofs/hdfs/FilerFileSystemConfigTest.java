package hanzofs.hdfs;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.junit.Before;
import org.junit.Test;

import static org.junit.Assert.*;

/**
 * Unit tests for FilerFileSystem configuration that don't require a running Hanzo S3 instance.
 * 
 * These tests verify basic properties and constants.
 */
public class FilerFileSystemConfigTest {

    private FilerFileSystem fs;
    private Configuration conf;

    @Before
    public void setUp() {
        fs = new FilerFileSystem();
        conf = new Configuration();
    }

    @Test
    public void testScheme() {
        assertEquals("s3", fs.getScheme());
    }

    @Test
    public void testConstants() {
        // Test that constants are defined correctly
        assertEquals("fs.hanzofs.filer.host", FilerFileSystem.FS_HANZOFS_FILER_HOST);
        assertEquals("fs.hanzofs.filer.port", FilerFileSystem.FS_HANZOFS_FILER_PORT);
        assertEquals("fs.hanzofs.filer.port.grpc", FilerFileSystem.FS_HANZOFS_FILER_PORT_GRPC);
        assertEquals(8888, FilerFileSystem.FS_HANZOFS_DEFAULT_PORT);
        assertEquals("fs.hanzofs.buffer.size", FilerFileSystem.FS_HANZOFS_BUFFER_SIZE);
        assertEquals(4 * 1024 * 1024, FilerFileSystem.FS_HANZOFS_DEFAULT_BUFFER_SIZE);
        assertEquals("fs.hanzofs.replication", FilerFileSystem.FS_HANZOFS_REPLICATION);
        assertEquals("fs.hanzofs.volume.server.access", FilerFileSystem.FS_HANZOFS_VOLUME_SERVER_ACCESS);
        assertEquals("fs.hanzofs.filer.cn", FilerFileSystem.FS_HANZOFS_FILER_CN);
    }

    @Test
    public void testWorkingDirectoryPathOperations() {
        // Test path operations that don't require initialization
        Path testPath = new Path("/test/path");
        assertTrue("Path should be absolute", testPath.isAbsolute());
        assertEquals("/test/path", testPath.toUri().getPath());
        
        Path childPath = new Path(testPath, "child");
        assertEquals("/test/path/child", childPath.toUri().getPath());
    }

    @Test
    public void testConfigurationProperties() {
        // Test that configuration can be set and read
        conf.set(FilerFileSystem.FS_HANZOFS_FILER_HOST, "testhost");
        assertEquals("testhost", conf.get(FilerFileSystem.FS_HANZOFS_FILER_HOST));
        
        conf.setInt(FilerFileSystem.FS_HANZOFS_FILER_PORT, 9999);
        assertEquals(9999, conf.getInt(FilerFileSystem.FS_HANZOFS_FILER_PORT, 0));
        
        conf.setInt(FilerFileSystem.FS_HANZOFS_BUFFER_SIZE, 8 * 1024 * 1024);
        assertEquals(8 * 1024 * 1024, conf.getInt(FilerFileSystem.FS_HANZOFS_BUFFER_SIZE, 0));
        
        conf.set(FilerFileSystem.FS_HANZOFS_REPLICATION, "001");
        assertEquals("001", conf.get(FilerFileSystem.FS_HANZOFS_REPLICATION));
        
        conf.set(FilerFileSystem.FS_HANZOFS_VOLUME_SERVER_ACCESS, "publicUrl");
        assertEquals("publicUrl", conf.get(FilerFileSystem.FS_HANZOFS_VOLUME_SERVER_ACCESS));
        
        conf.set(FilerFileSystem.FS_HANZOFS_FILER_CN, "test-cn");
        assertEquals("test-cn", conf.get(FilerFileSystem.FS_HANZOFS_FILER_CN));
    }

    @Test
    public void testDefaultBufferSize() {
        // Test default buffer size constant
        int expected = 4 * 1024 * 1024; // 4MB
        assertEquals(expected, FilerFileSystem.FS_HANZOFS_DEFAULT_BUFFER_SIZE);
    }

    @Test
    public void testDefaultPort() {
        // Test default port constant
        assertEquals(8888, FilerFileSystem.FS_HANZOFS_DEFAULT_PORT);
    }
}
