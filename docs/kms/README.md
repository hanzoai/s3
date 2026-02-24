# KMS Guide [![Discord](https://img.shields.io/discord/1234567890?label=discord)](https://hanzo.ai/discord)

S3 uses a key-management-system (KMS) to support SSE-S3. If a client requests SSE-S3, or auto-encryption is enabled, the S3 server encrypts each object with a unique object key which is protected by a master key managed by the KMS.

## Quick Start

S3 supports multiple KMS implementations via the [KES](https://github.com/minio/kes#kes) project. To run S3 with a KMS, fetch the root identity, set the following environment variables and then start your S3 server. If you haven't installed S3 yet, then follow the S3 [install instructions](https://docs.hanzo.ai/storage/operations/deployments) first.

### 1. Fetch the root identity

As the initial step, fetch the private key and certificate of the root identity:

```sh
curl -sSL --tlsv1.2 \
     -O 'https://raw.githubusercontent.com/minio/kes/master/root.key' \
     -O 'https://raw.githubusercontent.com/minio/kes/master/root.cert'
```

### 2. Set the S3-KES configuration

```sh
export S3_KMS_KES_ENDPOINT=https://kes.hanzo.ai:7373
export S3_KMS_KES_KEY_FILE=root.key
export S3_KMS_KES_CERT_FILE=root.cert
export S3_KMS_KES_KEY_NAME=my-minio-key
```

### 3. Start the S3 Server

```sh
export S3_ROOT_USER=minio
export S3_ROOT_PASSWORD=minio123
minio server ~/export
```

> For production deployments, you should run your own KES instance.

## Configuration Guides

A typical S3 deployment that uses a KMS for SSE-S3 looks like this:

```
    ┌────────────┐
    │ ┌──────────┴─┬─────╮          ┌────────────┐
    └─┤ ┌──────────┴─┬───┴──────────┤ ┌──────────┴─┬─────────────────╮
      └─┤ ┌──────────┴─┬─────┬──────┴─┤ KES Server ├─────────────────┤
        └─┤  S3  ├─────╯        └────────────┘            ┌────┴────┐
          └────────────┘                                        │   KMS   │
                                                                └─────────┘
```

In a given setup, there are `n` S3 instances talking to `m` KES servers but only `1` central KMS. The most simple setup consists of `1` S3 server or cluster talking to `1` KMS via `1` KES server.

The main difference between various S3-KMS deployments is the KMS implementation. The following table helps you select the right option for your use case:

| KMS                                                                                          | Purpose                                                           |
|:---------------------------------------------------------------------------------------------|:------------------------------------------------------------------|
| [Hashicorp Vault](https://github.com/minio/kes/wiki/Hashicorp-Vault-Keystore)                | Local KMS. S3 and KMS on-prem (**Recommended**)             |
| [AWS-KMS + SecretsManager](https://github.com/minio/kes/wiki/AWS-SecretsManager)             | Cloud KMS. S3 in combination with a managed KMS installation |
| [Gemalto KeySecure /Thales CipherTrust](https://github.com/minio/kes/wiki/Gemalto-KeySecure) | Local KMS. S3 and KMS On-Premises.                          |
| [Google Cloud Platform SecretManager](https://github.com/minio/kes/wiki/GCP-SecretManager)   | Cloud KMS. S3 in combination with a managed KMS installation |
| [FS](https://github.com/minio/kes/wiki/Filesystem-Keystore)                                  | Local testing or development (**Not recommended for production**) |

The S3-KES configuration is always the same - regardless of the underlying KMS implementation. Checkout the S3-KES [configuration example](https://github.com/minio/kes/wiki/Hanzo-S3-Object-Storage).

### Further references

- [Run S3 with TLS / HTTPS](https://docs.hanzo.ai/storage/operations/network-encryption)
- [Tweak the KES server configuration](https://github.com/minio/kes/wiki/Configuration)
- [Run a load balancer in front of KES](https://github.com/minio/kes/wiki/TLS-Proxy)
- [Understand the KES server concepts](https://github.com/minio/kes/wiki/Concepts)

## Auto Encryption

Auto-Encryption is useful when the S3 administrator wants to ensure that all data stored on S3 is encrypted at rest.

### Using `mc encrypt` (recommended)

S3 automatically encrypts all objects on buckets if KMS is successfully configured and bucket encryption configuration is enabled for each bucket as shown below:

```
mc encrypt set sse-s3 myminio/bucket/
```

Verify if S3 has `sse-s3` enabled

```
mc encrypt info myminio/bucket/
Auto encryption 'sse-s3' is enabled
```

### Using environment (not-recommended)

S3 automatically encrypts all objects on buckets if KMS is successfully configured and following ENV is enabled:

```
export S3_KMS_AUTO_ENCRYPTION=on
```

### Verify auto-encryption

> Note that auto-encryption only affects requests without S3 encryption headers. So, if a S3 client sends
> e.g. SSE-C headers, S3 will encrypt the object with the key sent by the client and won't reach out to
> the configured KMS.

To verify auto-encryption, use the following `mc` command:

```
mc cp test.file myminio/bucket/
test.file:              5 B / 5 B  ▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓  100.00% 337 B/s 0s
```

```
mc stat myminio/bucket/test.file
Name      : test.file
...
Encrypted :
  X-Amz-Server-Side-Encryption: AES256
```

## Encrypted Private Key

S3 supports encrypted KES client private keys. Therefore, you can use
an password-protected private keys for `S3_KMS_KES_KEY_FILE`.

When using password-protected private keys for accessing KES you need to
provide the password via:

```
export S3_KMS_KES_KEY_PASSWORD=<your-password>
```

Note that S3 only supports encrypted private keys - not encrypted certificates.
Certificates are no secrets and sent in plaintext as part of the TLS handshake.

## Explore Further

- [Use `mc` with S3](https://docs.hanzo.ai/storage/reference/mc)
- [Use `aws-cli` with S3](https://docs.hanzo.ai/storage/integrations/aws-cli)
- [Use `hanzo-s3` Go SDK](https://docs.hanzo.ai/storage/developers/go)
- [The S3 documentation website](https://docs.hanzo.ai/storage)
