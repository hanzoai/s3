.PHONY: all build install full_install clean test server benchmark benchmark_with_pprof warp_install s3-commands admin-generate admin-build admin-clean admin-dev admin-run admin-test admin-fmt admin-help

BINARY = s3
ADMIN_DIR = weed/admin

SOURCE_DIR = .
debug ?= 0

all: build

# Build the rebranded `s3` binary (Hanzo S3, fork of SeaweedFS) at the repo root.
build:
	go build -ldflags="-s -w" -o $(BINARY) ./weed

# Install the `s3` binary into $GOBIN/$GOPATH/bin.
install:
	cd weed && go build -ldflags="-s -w" -o s3 . && mv s3 "$$(go env GOPATH)/bin/s3"

s3-commands:
	cd weed && $(MAKE) weed-db weed-sql

warp_install:
	go install github.com/minio/warp@v0.7.6

full_install:
	cd weed && go build -ldflags="-s -w" -tags "elastic gocdk sqlite ydb tarantool tikv rclone" -o s3 . && mv s3 "$$(go env GOPATH)/bin/s3"

server: build
	./$(BINARY) -v 0 server -s3 -filer -filer.maxMB=64 -volume.max=0 -master.volumeSizeLimitMB=100 -volume.preStopSeconds=1 -s3.port=8000 -s3.allowDeleteBucketNotEmpty=true -s3.config=./docker/compose/s3.json -metricsPort=9324

benchmark: build warp_install
	pkill $(BINARY) || true
	pkill warp || true
	./$(BINARY) server -debug=$(debug) -s3 -filer -volume.max=0 -master.volumeSizeLimitMB=100 -volume.preStopSeconds=1 -s3.port=8000 -s3.allowDeleteBucketNotEmpty=false -s3.config=./docker/compose/s3.json &
	warp client &
	while ! nc -z localhost 8000 ; do sleep 1 ; done
	warp mixed --host=127.0.0.1:8000 --access-key=some_access_key1 --secret-key=some_secret_key1 --autoterm
	pkill warp
	pkill $(BINARY)

# curl -o profile "http://127.0.0.1:6060/debug/pprof/profile?debug=1"
benchmark_with_pprof: debug = 1
benchmark_with_pprof: benchmark

clean:
	go clean $(SOURCE_DIR)
	rm -f $(BINARY)

test:
	cd weed && go test -tags "elastic gocdk sqlite ydb tarantool tikv rclone" -v ./...

# Admin component targets (templ/sqlc codegen lives under weed/admin).
admin-generate:
	@cd $(ADMIN_DIR) && $(MAKE) generate

admin-build: admin-generate
	@echo "Building admin component..."
	@cd $(ADMIN_DIR) && $(MAKE) build

admin-clean:
	@echo "Cleaning admin component..."
	@cd $(ADMIN_DIR) && $(MAKE) clean

admin-dev:
	@echo "Starting admin development server..."
	@cd $(ADMIN_DIR) && $(MAKE) dev

admin-run:
	@echo "Running admin server..."
	@cd $(ADMIN_DIR) && $(MAKE) run

admin-test:
	@echo "Testing admin component..."
	@cd $(ADMIN_DIR) && $(MAKE) test

admin-fmt:
	@echo "Formatting admin component..."
	@cd $(ADMIN_DIR) && $(MAKE) fmt

admin-help:
	@echo "Admin component help..."
	@cd $(ADMIN_DIR) && $(MAKE) help
