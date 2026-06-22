# Docker

## Compose V2 
Hanzo now uses the `v2` syntax `docker compose`

If you rely on using Docker Compose as docker-compose (with a hyphen), you can set up Compose V2 to act as a drop-in replacement of the previous docker-compose. Refer to the [Installing Compose](https://docs.docker.com/compose/install/) section for detailed instructions on upgrading.

Confirm your system has docker compose v2 with a version check
```bash
$ docker compose version
Docker Compose version v2.10.2
```

## Try it out

```bash

wget https://raw.githubusercontent.com/hanzo/hanzo/master/docker/hanzo-compose.yml

docker compose -f hanzo-compose.yml -p hanzo up

```

## Try latest tip

```bash

wget https://raw.githubusercontent.com/hanzo/hanzo/master/docker/hanzo-dev-compose.yml

docker compose -f hanzo-dev-compose.yml -p hanzo up

```

## Local Development

```bash
cd $GOPATH/src/github.com/hanzoai/s3/docker
make
```

### S3 cmd

list
```
s3cmd --no-ssl --host=127.0.0.1:8333 ls s3://
```

## Build and push a multiarch build

Make sure that `docker buildx` is supported (might be an experimental docker feature)
```bash
BUILDER=$(docker buildx create --driver docker-container --use)
docker buildx build --pull --push --platform linux/386,linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 . -t chrislusf/hanzo
docker buildx stop $BUILDER
```

## Minio debugging
```
mc config host add local http://127.0.0.1:9000 some_access_key1 some_secret_key1
mc admin trace --all --verbose local
```
