# How to monitor S3 with Grafana [![Discord](https://img.shields.io/discord/1234567890?label=discord)](https://hanzo.ai/discord)

[Grafana](https://grafana.com/) allows you to query, visualize, alert on and understand your metrics no matter where they are stored. Create, explore, and share dashboards with your team and foster a data driven culture.

## Prerequisites

- Prometheus and S3 configured as explained in [document here](https://github.com/hanzoai/storage/blob/master/docs/metrics/prometheus/README.md).
- Grafana installed as explained [here](https://grafana.com/grafana/download).

## S3 Grafana Dashboard

Visualize S3 metrics with our official Grafana dashboard available on the [Grafana dashboard portal](https://grafana.com/grafana/dashboards/13502).

Refer to the dashboard [json file here](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/minio-dashboard.json).

![Grafana](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/grafana-minio.png)

Node level Replication metrics can be viewed in the Grafana dashboard using [json file here](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/replication/minio-replication-node.json)

![Grafana](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/replication/grafana-replication-node.png)

Cluster level Replication metrics can be viewed in the Grafana dashboard using [json file here](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/replication/minio-replication-cluster.json)

![Grafana](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/replication/grafana-replication-cluster.png)

Bucket metrics can be viewed in the Grafana dashboard using [json file here](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/bucket/minio-bucket.json)

![Grafana](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/bucket/grafana-bucket.png)

Node metrics can be viewed in the Grafana dashboard using [json file here](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/node/minio-node.json)

![Grafana](https://raw.githubusercontent.com/hanzoai/storage/master/docs/metrics/prometheus/grafana/node/grafana-node.png)

Note: All these dashboards are provided as an example and need basis they should be customized as well as new graphs should be added.
