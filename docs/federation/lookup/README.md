# Federation Quickstart Guide *Federation feature is deprecated and should be avoided for future deployments*

This document explains how to configure S3 with `Bucket lookup from DNS` style federation.

## Get started

### 1. Prerequisites

Install S3 - [S3 Quickstart Guide](https://docs.hanzo.ai/community/minio-object-store/operations/deployments/baremetal-deploy-minio-on-redhat-linux.html).

### 2. Run S3 in federated mode

Bucket lookup from DNS federation requires two dependencies

- etcd (for bucket DNS service records)
- CoreDNS (for DNS management based on populated bucket DNS service records, optional)

## Architecture

![bucket-lookup](https://github.com/hanzoai/s3/blob/master/docs/federation/lookup/bucket-lookup.png?raw=true)

### Environment variables

#### S3_ETCD_ENDPOINTS

This is comma separated list of etcd servers that you want to use as the S3 federation back-end. This should
be same across the federated deployment, i.e. all the S3 instances within a federated deployment should use same
etcd back-end.

#### S3_DOMAIN

This is the top level domain name used for the federated setup. This domain name should ideally resolve to a load-balancer
running in front of all the federated S3 instances. The domain name is used to create sub domain entries to etcd. For
example, if the domain is set to `domain.com`, the buckets `bucket1`, `bucket2` will be accessible as `bucket1.domain.com`
and `bucket2.domain.com`.

#### S3_PUBLIC_IPS

This is comma separated list of IP addresses to which buckets created on this S3 instance will resolve to. For example,
a bucket `bucket1` created on current S3 instance will be accessible as `bucket1.domain.com`, and the DNS entry for
`bucket1.domain.com` will point to IP address set in `S3_PUBLIC_IPS`.

- This field is mandatory for standalone and erasure code S3 server deployments, to enable federated mode.
- This field is optional for distributed deployments. If you don't set this field in a federated setup, we use the IP addresses of
hosts passed to the S3 server startup and use them for DNS entries.

### Run Multiple Clusters

> cluster1

```sh
export S3_ETCD_ENDPOINTS="http://remote-etcd1:2379,http://remote-etcd2:4001"
export S3_DOMAIN=domain.com
export S3_PUBLIC_IPS=44.35.2.1,44.35.2.2,44.35.2.3,44.35.2.4
minio server http://rack{1...4}.host{1...4}.domain.com/mnt/export{1...32}
```

> cluster2

```sh
export S3_ETCD_ENDPOINTS="http://remote-etcd1:2379,http://remote-etcd2:4001"
export S3_DOMAIN=domain.com
export S3_PUBLIC_IPS=44.35.1.1,44.35.1.2,44.35.1.3,44.35.1.4
minio server http://rack{5...8}.host{5...8}.domain.com/mnt/export{1...32}
```

In this configuration you can see `S3_ETCD_ENDPOINTS` points to the etcd backend which manages S3's
`config.json` and bucket DNS SRV records. `S3_DOMAIN` indicates the domain suffix for the bucket which
will be used to resolve bucket through DNS. For example if you have a bucket such as `mybucket`, the
client can use now `mybucket.domain.com` to directly resolve itself to the right cluster. `S3_PUBLIC_IPS`
points to the public IP address where each cluster might be accessible, this is unique for each cluster.

NOTE: `mybucket` only exists on one cluster either `cluster1` or `cluster2` this is random and
is decided by how `domain.com` gets resolved, if there is a round-robin DNS on `domain.com` then
it is randomized which cluster might provision the bucket.

### 3. Test your setup

To test this setup, access the S3 server via browser or [`mc`](https://docs.hanzo.ai/community/minio-object-store/reference/minio-mc.html#quickstart). You’ll see the uploaded files are accessible from the all the S3 endpoints.

## Explore Further

- [Use `mc` with S3 Server](https://docs.hanzo.ai/community/minio-object-store/reference/minio-mc.html)
- [Use `aws-cli` with S3 Server](https://docs.hanzo.ai/community/minio-object-store/integrations/aws-cli-with-minio.html)
- [Use `minio-go` SDK with S3 Server](https://docs.hanzo.ai/community/minio-object-store/developers/go/minio-go.html)
- [The S3 documentation website](https://docs.hanzo.ai/community/minio-object-store/index.html)
