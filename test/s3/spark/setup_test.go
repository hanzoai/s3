package spark

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hanzoai/s3/test/testutil"
	"github.com/testcontainers/testcontainers-go"
)

var (
	miniProcessMu   sync.Mutex
	lastMiniProcess *exec.Cmd
)

func stopPreviousMini() {
	miniProcessMu.Lock()
	defer miniProcessMu.Unlock()

	if lastMiniProcess != nil && lastMiniProcess.Process != nil {
		_ = lastMiniProcess.Process.Kill()
		_ = lastMiniProcess.Wait()
	}
	lastMiniProcess = nil
}

func registerMiniProcess(cmd *exec.Cmd) {
	miniProcessMu.Lock()
	lastMiniProcess = cmd
	miniProcessMu.Unlock()
}

func clearMiniProcess(cmd *exec.Cmd) {
	miniProcessMu.Lock()
	if lastMiniProcess == cmd {
		lastMiniProcess = nil
	}
	miniProcessMu.Unlock()
}

type TestEnvironment struct {
	dockerAvailable  bool
	s3Binary       string
	hanzoDataDir string
	s3LogPath      string
	s3LogFile      *os.File
	masterPort       int
	filerPort        int
	s3Port           int
	accessKey        string
	secretKey        string
	sparkContainer   testcontainers.Container
	masterProcess    *exec.Cmd
}

func NewTestEnvironment() *TestEnvironment {
	env := &TestEnvironment{
		accessKey: "test",
		secretKey: "test",
	}

	cmd := exec.Command("docker", "version")
	env.dockerAvailable = cmd.Run() == nil

	if s3Path, err := exec.LookPath("s3"); err == nil {
		env.s3Binary = s3Path
	}

	return env
}

func (env *TestEnvironment) StartHanzo(t *testing.T) {
	t.Helper()

	if env.s3Binary == "" {
		t.Skip("s3 binary not found in PATH, skipping Spark S3 integration test")
	}

	stopPreviousMini()

	var err error
	env.hanzoDataDir, err = os.MkdirTemp("", "hanzo-s3-spark-test-")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	ports := testutil.MustFreeMiniPorts(t, []string{"Master", "Filer", "S3"})
	env.masterPort = ports[0]
	env.filerPort = ports[1]
	env.s3Port = ports[2]

	bindIP := testutil.FindBindIP()
	iamConfigPath, err := testutil.WriteIAMConfig(env.hanzoDataDir, env.accessKey, env.secretKey)
	if err != nil {
		t.Fatalf("failed to create IAM config: %v", err)
	}

	env.masterProcess = exec.Command(
		env.s3Binary, "mini",
		"-ip", bindIP,
		"-ip.bind", "0.0.0.0",
		"-master.port", fmt.Sprintf("%d", env.masterPort),
		"-filer.port", fmt.Sprintf("%d", env.filerPort),
		"-s3.port", fmt.Sprintf("%d", env.s3Port),
		"-s3.config", iamConfigPath,
		"-dir", env.hanzoDataDir,
	)
	s3LogPath := filepath.Join(env.hanzoDataDir, "s3-mini.log")
	s3LogFile, err := os.Create(s3LogPath)
	if err != nil {
		t.Fatalf("failed to create s3 log file: %v", err)
	}
	env.s3LogPath = s3LogPath
	env.s3LogFile = s3LogFile
	env.masterProcess.Stdout = s3LogFile
	env.masterProcess.Stderr = s3LogFile
	env.masterProcess.Env = append(os.Environ(),
		"AWS_ACCESS_KEY_ID="+env.accessKey,
		"AWS_SECRET_ACCESS_KEY="+env.secretKey,
	)

	if err := env.masterProcess.Start(); err != nil {
		t.Fatalf("failed to start s3 mini: %v", err)
	}
	registerMiniProcess(env.masterProcess)

	if !testutil.WaitForPort(env.masterPort, testutil.HanzoMiniStartupTimeout) {
		t.Fatalf("s3 mini failed to start - master port %d not listening", env.masterPort)
	}
	if !testutil.WaitForPort(env.filerPort, testutil.HanzoMiniStartupTimeout) {
		t.Fatalf("s3 mini failed to start - filer port %d not listening", env.filerPort)
	}
	if !testutil.WaitForService(fmt.Sprintf("http://127.0.0.1:%d/status", env.s3Port), testutil.HanzoMiniStartupTimeout) {
		t.Fatalf("s3 mini failed to start - s3 endpoint http://127.0.0.1:%d/status not responding", env.s3Port)
	}
}

func (env *TestEnvironment) startSparkContainer(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	req := testcontainers.ContainerRequest{
		Image:        "apache/spark:3.5.1",
		ExposedPorts: []string{"4040/tcp"},
		Env: map[string]string{
			"SPARK_LOCAL_IP": "localhost",
		},
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
		Cmd:        []string{"/bin/sh", "-c", "sleep 7200"},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start spark container: %v", err)
	}
	env.sparkContainer = container
}

func (env *TestEnvironment) Cleanup(t *testing.T) {
	if env.masterProcess != nil && env.masterProcess.Process != nil {
		_ = env.masterProcess.Process.Kill()
		_ = env.masterProcess.Wait()
	}
	clearMiniProcess(env.masterProcess)
	if env.s3LogFile != nil {
		_ = env.s3LogFile.Close()
	}

	if t.Failed() && os.Getenv("CI") != "" && env.s3LogPath != "" {
		logData, err := os.ReadFile(env.s3LogPath)
		if err != nil {
			t.Logf("failed to read s3 mini log file %s: %v", env.s3LogPath, err)
		} else {
			// Print the tail to keep CI output manageable while preserving failure context.
			const maxTailBytes = 64 * 1024
			start := 0
			if len(logData) > maxTailBytes {
				start = len(logData) - maxTailBytes
			}
			t.Logf("s3 mini logs (tail, %d bytes):\n%s", len(logData)-start, string(logData[start:]))
		}
	}

	if env.sparkContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = env.sparkContainer.Terminate(ctx)
	}

	if env.hanzoDataDir != "" {
		_ = os.RemoveAll(env.hanzoDataDir)
	}
}

func runSparkPyScript(t *testing.T, container testcontainers.Container, script string, s3Port int) (int, string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	pythonScript := fmt.Sprintf(`
import glob
import os
import sys

spark_home = os.environ.get("SPARK_HOME", "/opt/spark")
python_path = os.path.join(spark_home, "python")
py4j_glob = glob.glob(os.path.join(python_path, "lib", "py4j-*.zip"))
ivy_dir = "/tmp/ivy"
os.makedirs(ivy_dir, exist_ok=True)
os.environ["AWS_REGION"] = "us-east-1"
os.environ["AWS_DEFAULT_REGION"] = "us-east-1"
os.environ["AWS_ACCESS_KEY_ID"] = "test"
os.environ["AWS_SECRET_ACCESS_KEY"] = "test"

if python_path not in sys.path:
    sys.path.insert(0, python_path)
if py4j_glob and py4j_glob[0] not in sys.path:
    sys.path.insert(0, py4j_glob[0])

from pyspark.sql import SparkSession

spark = (SparkSession.builder
    .master("local[2]")
    .appName("Hanzo S3 Spark Issue 8234 Repro")
    .config("spark.sql.catalogImplementation", "hive")
    .config("spark.executor.extraJavaOptions", "-Djdk.tls.client.protocols=TLSv1")
    .config("spark.sql.parquet.int96RebaseModeInRead", "CORRECTED")
    .config("spark.sql.parquet.int96RebaseModeInWrite", "CORRECTED")
    .config("spark.sql.parquet.datetimeRebaseModeInRead", "CORRECTED")
    .config("spark.sql.parquet.datetimeRebaseModeInWrite", "CORRECTED")
    .config("spark.sql.parquet.enableVectorizedReader", "false")
    .config("spark.sql.parquet.mergeSchema", "false")
    .config("spark.sql.parquet.writeLegacyFormat", "true")
    .config("spark.sql.broadcastTimeout", "-1")
    .config("spark.network.timeout", "600")
    .config("spark.hadoop.hive.metastore.schema.verification", "false")
    .config("spark.hadoop.hive.metastore.schema.verification.record.version", "false")
    .config("spark.jars.ivy", ivy_dir)
    .config("spark.jars.packages", "org.apache.hadoop:hadoop-aws:3.3.2,com.amazonaws:aws-java-sdk-bundle:1.12.262")
    .config("spark.hadoop.fs.s3a.access.key", "test")
    .config("spark.hadoop.fs.s3a.secret.key", "test")
    .config("spark.hadoop.fs.s3a.endpoint", "host.docker.internal:%d")
    .config("spark.hadoop.mapreduce.fileoutputcommitter.algorithm.version", "1")
    .config("spark.hadoop.fs.s3a.path.style.access", "true")
    .config("spark.hadoop.fs.s3a.fast.upload", "true")
    .config("spark.hadoop.fs.s3a.connection.ssl.enabled", "false")
    .config("spark.hadoop.fs.s3a.multiobjectdelete.enable", "true")
    .config("spark.hadoop.fs.s3a.directory.marker.retention", "keep")
    .config("spark.hadoop.fs.s3a.change.detection.version.required", "false")
    .config("spark.hadoop.fs.s3a.change.detection.mode", "warn")
    .config("spark.local.dir", "/tmp/spark-temp")
    .getOrCreate())

%s
`, s3Port, script)

	code, out, err := container.Exec(ctx, []string{"python3", "-c", pythonScript})
	var output string
	if out != nil {
		outputBytes, readErr := io.ReadAll(out)
		if readErr != nil {
			output = fmt.Sprintf("failed to read container output: %v", readErr)
		} else {
			output = string(outputBytes)
		}
	}
	if err != nil {
		output = output + fmt.Sprintf("\ncontainer exec error: %v\n", err)
	}
	return code, output
}

func createObjectBucket(t *testing.T, env *TestEnvironment, bucketName string) {
	t.Helper()

	cfg := aws.Config{
		Region:       "us-east-1",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(env.accessKey, env.secretKey, "")),
		BaseEndpoint: aws.String(fmt.Sprintf("http://localhost:%d", env.s3Port)),
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	_, err := client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		t.Fatalf("failed to create object bucket %s: %v", bucketName, err)
	}
}
