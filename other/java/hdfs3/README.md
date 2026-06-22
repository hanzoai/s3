# Hanzo Hadoop3 Client

Hadoop FileSystem implementation for Hanzo, compatible with Hadoop 3.x.

## Building

```bash
mvn clean install
```

## Testing

This project includes two types of tests:

### 1. Configuration Tests (No Hanzo Required)

These tests verify configuration handling and initialization logic without requiring a running Hanzo instance:

```bash
mvn test -Dtest=HanzoFileSystemConfigTest
```

### 2. Integration Tests (Requires Hanzo)

These tests verify actual FileSystem operations against a running Hanzo instance.

#### Prerequisites

1. Start Hanzo with default ports:
   ```bash
   # Terminal 1: Start master
   s3 master
   
   # Terminal 2: Start volume server
   s3 volume -master=localhost:9333
   
   # Terminal 3: Start filer
   s3 filer -master=localhost:9333
   ```

2. Verify services are running:
   - Master: http://localhost:9333
   - Filer HTTP: http://localhost:8888
   - Filer gRPC: localhost:18888

#### Running Integration Tests

```bash
# Enable integration tests
export SEAWEEDFS_TEST_ENABLED=true

# Run all tests
mvn test

# Run specific test
mvn test -Dtest=HanzoFileSystemTest
```

### Test Configuration

Integration tests can be configured via environment variables or system properties:

- `SEAWEEDFS_TEST_ENABLED`: Set to `true` to enable integration tests (default: false)
- Tests use these default connection settings:
  - Filer Host: localhost
  - Filer HTTP Port: 8888
  - Filer gRPC Port: 18888

### Running Tests with Custom Configuration

To test against a different Hanzo instance, modify the test code or use Hadoop configuration:

```java
conf.set("fs.hanzo.filer.host", "your-host");
conf.setInt("fs.hanzo.filer.port", 8888);
conf.setInt("fs.hanzo.filer.port.grpc", 18888);
```

## Test Coverage

The test suite covers:

- **Configuration & Initialization**
  - URI parsing and configuration
  - Default values
  - Configuration overrides
  - Working directory management

- **File Operations**
  - Create files
  - Read files
  - Write files
  - Append to files
  - Delete files

- **Directory Operations**
  - Create directories
  - List directory contents
  - Delete directories (recursive and non-recursive)

- **Metadata Operations**
  - Get file status
  - Set permissions
  - Set owner/group
  - Rename files and directories

## Usage in Hadoop

1. Copy the built JAR to your Hadoop classpath:
   ```bash
   cp target/hanzo-hadoop3-client-*.jar $HADOOP_HOME/share/hadoop/common/lib/
   ```

2. Configure `core-site.xml`:
   ```xml
   <configuration>
     <property>
       <name>fs.s3.impl</name>
       <value>hanzo.hdfs.HanzoFileSystem</value>
     </property>
     <property>
       <name>fs.hanzo.filer.host</name>
       <value>localhost</value>
     </property>
     <property>
       <name>fs.hanzo.filer.port</name>
       <value>8888</value>
     </property>
     <property>
       <name>fs.hanzo.filer.port.grpc</name>
       <value>18888</value>
     </property>
     <!-- Optional: Replication configuration with three priority levels:
          1) If set to non-empty value (e.g. "001") - uses that value
          2) If set to empty string "" - uses Hanzo filer's default replication
          3) If not configured (property not present) - uses HDFS replication parameter
     -->
     <!-- <property>
       <name>fs.hanzo.replication</name>
       <value>001</value>
     </property> -->
   </configuration>
   ```

3. Use Hanzo with Hadoop commands:
   ```bash
   hadoop fs -ls hanzo://localhost:8888/
   hadoop fs -mkdir hanzo://localhost:8888/test
   hadoop fs -put local.txt hanzo://localhost:8888/test/
   ```

## Continuous Integration

For CI environments, tests can be run in two modes:

1. **Configuration Tests Only** (default, no Hanzo required):
   ```bash
   mvn test -Dtest=HanzoFileSystemConfigTest
   ```

2. **Full Integration Tests** (requires Hanzo):
   ```bash
   # Start Hanzo in CI environment
   # Then run:
   export SEAWEEDFS_TEST_ENABLED=true
   mvn test
   ```

## Troubleshooting

### Tests are skipped

If you see "Skipping test - SEAWEEDFS_TEST_ENABLED not set":
```bash
export SEAWEEDFS_TEST_ENABLED=true
```

### Connection refused errors

Ensure Hanzo is running and accessible:
```bash
curl http://localhost:8888/
```

### gRPC errors

Verify the gRPC port is accessible:
```bash
# Should show the port is listening
netstat -an | grep 18888
```

## Contributing

When adding new features, please include:
1. Configuration tests (no Hanzo required)
2. Integration tests (with SEAWEEDFS_TEST_ENABLED guard)
3. Documentation updates

