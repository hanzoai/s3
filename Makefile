.PHONY: test admin-generate admin-build admin-clean admin-dev admin-run admin-test admin-fmt admin-help s3-commands

BINARY = s3
ADMIN_DIR = s3/admin

SOURCE_DIR = .
debug ?= 0

all: install

install: admin-generate
	cd s3; go install

s3-commands:
	cd s3 && $(MAKE) s3-db s3-sql

full_install: admin-generate
	cd s3; go install -tags "elastic gocdk sqlite ydb tarantool tikv"

server: install
	s3 -v 0 server -s3 -filer -filer.maxMB=64 -volume.max=0 -master.volumeSizeLimitMB=100 -volume.preStopSeconds=1 -s3.port=8000 -s3.allowDeleteBucketNotEmpty=true -s3.config=./docker/compose/s3.json -metricsPort=9324

test: admin-generate
	cd s3; go test -tags "elastic gocdk sqlite ydb tarantool tikv" -v ./...

# Admin component targets
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
