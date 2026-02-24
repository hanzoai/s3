# Deploy S3 on Kubernetes [![Slack](https://hanzo.ai/discord/slack?type=svg)](https://hanzo.ai/discord)  [![Docker Pulls](https://img.shields.io/docker/pulls/minio/minio.svg?maxAge=604800)](https://hub.docker.com/r/minio/minio/)

S3 is a high performance distributed object storage server, designed for large-scale private cloud infrastructure. S3 is designed in a cloud-native manner to scale sustainably in multi-tenant environments. Orchestration platforms like Kubernetes provide perfect cloud-native environment to deploy and scale S3.

## S3 Deployment on Kubernetes

There are multiple options to deploy S3 on Kubernetes:

- S3-Operator: Operator offers seamless way to create and update highly available distributed S3 clusters. Refer [S3 Operator documentation](https://github.com/hanzoai/s3-operator/blob/master/README.md) for more details.

- Helm Chart: S3 Helm Chart offers customizable and easy S3 deployment with a single command. Refer [S3 Helm Chart documentation](https://github.com/hanzoai/s3/tree/master/helm/minio) for more details.

## Monitoring S3 in Kubernetes

S3 server exposes un-authenticated liveness endpoints so Kubernetes can natively identify unhealthy S3 containers. S3 also exposes Prometheus compatible data on a different endpoint to enable Prometheus users to natively monitor their S3 deployments.

## Explore Further

- [S3 Erasure Code QuickStart Guide](https://docs.hanzo.ai/community/minio-object-store/operations/concepts/erasure-coding.html)
- [Kubernetes Documentation](https://kubernetes.io/docs/home/)
- [Helm package manager for kubernetes](https://helm.sh/)
