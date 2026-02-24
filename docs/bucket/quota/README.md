# Bucket Quota Configuration Quickstart Guide

![quota](https://raw.githubusercontent.com/minio/minio/master/docs/bucket/quota/bucketquota.png)

Buckets can be configured to have `Hard` quota - it disallows writes to the bucket after configured quota limit is reached.

## Prerequisites

- Install S3 - [S3 Quickstart Guide](https://docs.hanzo.ai/community/minio-object-store/operations/deployments/baremetal-deploy-minio-on-redhat-linux.html#procedure).
- [Use `mc` with S3 Server](https://docs.hanzo.ai/community/minio-object-store/reference/minio-mc.html#quickstart)

## Set bucket quota configuration

### Set a hard quota of 1GB for a bucket `mybucket` on S3 object storage

```sh
mc admin bucket quota myminio/mybucket --hard 1gb
```

### Verify the quota configured on `mybucket` on S3

```sh
mc admin bucket quota myminio/mybucket
```

### Clear bucket quota configuration for `mybucket` on S3

```sh
mc admin bucket quota myminio/mybucket --clear
```
