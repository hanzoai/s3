all: gen

.PHONY : gen

gen: dev

cgo ?= 0
ldflags_extra ?= -extldflags -static

binary:
	export SWCOMMIT=$(shell git rev-parse --short HEAD)
	export SWLDFLAGS="-X github.com/hanzoai/s3/s3/util/version.COMMIT=$(SWCOMMIT)"
	cd ../s3 && CGO_ENABLED=$(cgo) GOOS=linux go build $(options) -tags "$(tags)" -ldflags "-s -w $(ldflags_extra) $(SWLDFLAGS)" -o s3_binary && mv s3_binary ../docker/s3
	cd ../other/mq_client_example/agent_pub_record && CGO_ENABLED=$(cgo) GOOS=linux go build && mv agent_pub_record ../../../docker/
	cd ../other/mq_client_example/agent_sub_record && CGO_ENABLED=$(cgo) GOOS=linux go build && mv agent_sub_record ../../../docker/

# Race detector requires CGO and dynamic linking - don't use -static
binary_race: options = -race
binary_race: cgo = 1
binary_race: ldflags_extra =
binary_race: binary

build: binary
	docker build --no-cache -t chrislusf/hanzo:local -f Dockerfile.local .

build_e2e: binary_race
	docker buildx build \
		--cache-from=type=local,src=/tmp/.buildx-cache \
		--cache-to=type=local,dest=/tmp/.buildx-cache-new,mode=max \
		--load \
		-t chrislusf/hanzo:e2e \
		-f Dockerfile.e2e .
	# Move cache to avoid growing cache size
	rm -rf /tmp/.buildx-cache || true
	mv /tmp/.buildx-cache-new /tmp/.buildx-cache || true

go_build: # make go_build tags=elastic,ydb,gocdk,hdfs,5BytesOffset,tarantool
	docker build --build-arg TAGS=$(tags) --no-cache -t chrislusf/hanzo:go_build -f Dockerfile.go_build .

go_build_large_disk:
	docker build --build-arg TAGS=large_disk --no-cache -t chrislusf/hanzo:large_disk -f Dockerfile.go_build .

build_rocksdb_dev_env:
	docker build --no-cache -t chrislusf/rocksdb_dev_env -f Dockerfile.rocksdb_dev_env .

build_rocksdb_local: build_rocksdb_dev_env
	cd .. ; docker build --no-cache -t chrislusf/hanzo:rocksdb_local -f docker/Dockerfile.rocksdb_large_local .

build_rocksdb:
	docker build --no-cache -t chrislusf/hanzo:rocksdb -f Dockerfile.rocksdb_large .

build_tarantool_dev_env:
	docker build --no-cache -t chrislusf/tarantool_dev_env -f Dockerfile.tarantool.dev_env .

s3tests_build:
	docker build --no-cache -t chrislusf/ceph-s3-tests:local -f Dockerfile.s3tests .

dev: build
	docker compose -f compose/local-dev-compose.yml -p hanzo up

dev_race: binary_race
	docker compose -f compose/local-dev-compose.yml -p hanzo up

dev_tls: build certstrap
	ENV_FILE="tls.env" docker compose -f compose/local-dev-compose.yml -p hanzo up

dev_mount: build
	docker compose -f compose/local-mount-compose.yml -p hanzo up

run_image: build
	docker run --rm -ti --device /dev/fuse --cap-add SYS_ADMIN --entrypoint /bin/sh chrislusf/hanzo:local

profile_mount: build
	docker compose -f compose/local-mount-profile-compose.yml -p hanzo up

k8s: build
	docker compose -f compose/local-k8s-compose.yml -p hanzo up

dev_registry: build
	docker compose -f compose/local-registry-compose.yml -p hanzo up

dev_replicate:
	docker build --build-arg TAGS=gocdk --no-cache -t chrislusf/hanzo:local -f Dockerfile.go_build .
	docker compose -f compose/local-replicate-compose.yml -p hanzo up

dev_auditlog: build
	docker compose -f compose/local-auditlog-compose.yml -p hanzo up

dev_nextcloud: build
	docker compose -f compose/local-nextcloud-compose.yml -p hanzo up

cluster: build
	docker compose -f compose/local-cluster-compose.yml -p hanzo up

2clusters: build
	docker compose -f compose/local-clusters-compose.yml -p hanzo up

2mount: build
	docker compose -f compose/local-sync-mount-compose.yml -p hanzo up

filer_backup: build
	docker compose -f compose/local-filer-backup-compose.yml -p hanzo up

hashicorp_raft: build
	docker compose -f compose/local-hashicorp-raft-compose.yml -p hanzo up

s3tests: build s3tests_build
	docker compose -f compose/local-s3tests-compose.yml -p hanzo up

brokers: build
	docker compose -f compose/local-brokers-compose.yml -p hanzo up

agent: build
	docker compose -f compose/local-mq-test.yml -p hanzo up

filer_etcd: build
	docker stack deploy -c compose/swarm-etcd.yml fs

test_etcd: build
	docker compose -f compose/test-etcd-filer.yml -p hanzo up

test_ydb: tags = ydb
test_ydb: build
	docker compose -f compose/test-ydb-filer.yml -p hanzo up

test_tarantool: tags = tarantool
test_tarantool: build_tarantool_dev_env build
	docker compose -f compose/test-tarantool-filer.yml -p hanzo up

test_keycloak_s3: build
	docker compose -f compose/test-keycloak-s3.yml -p hanzo up

clean:
	rm ./s3

certstrap:
	go install -v github.com/square/certstrap@latest
	certstrap --depot-path compose/tls init --curve P-256 --passphrase "" --common-name "Hanzo CA" || true
	certstrap --depot-path compose/tls request-cert --ou "Hanzo" --curve P-256 --passphrase "" --domain localhost --common-name volume01.dev || true
	certstrap --depot-path compose/tls request-cert --ou "Hanzo" --curve P-256 --passphrase "" --common-name master01.dev || true
	certstrap --depot-path compose/tls request-cert --ou "Hanzo" --curve P-256 --passphrase "" --common-name filer01.dev || true
	certstrap --depot-path compose/tls request-cert --ou "Hanzo" --curve P-256 --passphrase "" --common-name client01.dev || true
	certstrap --depot-path compose/tls sign --CA "Hanzo CA" volume01.dev || true
	certstrap --depot-path compose/tls sign --CA "Hanzo CA" master01.dev || true
	certstrap --depot-path compose/tls sign --CA "Hanzo CA" filer01.dev || true
	certstrap --depot-path compose/tls sign --CA "Hanzo CA" client01.dev || true
