.PHONY: test test-short build-s3 clean

# Build s3 binary first if needed
build-s3:
	cd ../../s3 && go build -o s3 .

# Run all integration tests
test: build-s3
	go test -v -timeout 5m .

# Run tests in short mode (skip integration tests)
test-short:
	go test -v -short .

# Run specific test
test-basic: build-s3
	go test -v -timeout 3m -run TestMetadataSubscribeBasic .

test-stall: build-s3
	go test -v -timeout 5m -run TestMetadataSubscribeSingleFilerNoStall .

test-resume: build-s3
	go test -v -timeout 3m -run TestMetadataSubscribeResumeFromDisk .

clean:
	rm -f ../../s3/s3


