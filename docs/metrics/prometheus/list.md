# Cluster Metrics

S3 collects the following metrics at the cluster level.
Metrics may include one or more labels, such as the server that calculated that metric.

These metrics can be obtained from any S3 server once per collection by using the following URL:

```shell
https://HOSTNAME:PORT/minio/v2/metrics/cluster
```

Replace ``HOSTNAME:PORT`` with the hostname of your S3 deployment.
For deployments behind a load balancer, use the load balancer hostname instead of a single node hostname.

## Audit Metrics

| Name                              | Description                                               |
|:----------------------------------|:----------------------------------------------------------|
| `s3_audit_failed_messages`     | Total number of messages that failed to send since start. |
| `s3_audit_target_queue_length` | Number of unsent messages in queue for target.            |
| `s3_audit_total_messages`      | Total number of messages sent since start.                |

## Cluster Capacity Metrics

| Name                                         | Description                                                    |
|:---------------------------------------------|:---------------------------------------------------------------|
| `s3_cluster_capacity_raw_free_bytes`      | Total free capacity online in the cluster.                     |
| `s3_cluster_capacity_raw_total_bytes`     | Total capacity online in the cluster.                          |
| `s3_cluster_capacity_usable_free_bytes`   | Total free usable capacity online in the cluster.              |
| `s3_cluster_capacity_usable_total_bytes`  | Total usable capacity online in the cluster.                   |
| `s3_cluster_objects_size_distribution`    | Distribution of object sizes across a cluster                  |
| `s3_cluster_objects_version_distribution` | Distribution of object versions across a cluster               |
| `s3_cluster_usage_object_total`           | Total number of objects in a cluster                           |
| `s3_cluster_usage_total_bytes`            | Total cluster usage in bytes                                   |
| `s3_cluster_usage_version_total`          | Total number of versions (includes delete marker) in a cluster |
| `s3_cluster_usage_deletemarker_total`     | Total number of delete markers in a cluster                    |
| `s3_cluster_bucket_total`                 | Total number of buckets in the cluster                         |

## Cluster Drive Metrics

| Name                                | Description                           |
|:------------------------------------|:--------------------------------------|
| `s3_cluster_drive_offline_total` | Total drives offline in this cluster. |
| `s3_cluster_drive_online_total`  | Total drives online in this cluster.  |
| `s3_cluster_drive_total`         | Total drives in this cluster.         |

## Cluster ILM Metrics

| Name                                      | Description                                      |
|:------------------------------------------|:-------------------------------------------------|
| `s3_cluster_ilm_transitioned_bytes`    | Total bytes transitioned to a tier.              |
| `s3_cluster_ilm_transitioned_objects`  | Total number of objects transitioned to a tier.  |
| `s3_cluster_ilm_transitioned_versions` | Total number of versions transitioned to a tier. |

## Cluster KMS Metrics

| Name                                | Description                                                                              |
|:------------------------------------|:-----------------------------------------------------------------------------------------|
| `s3_cluster_kms_online`          | Reports whether the KMS is online (1) or offline (0).                                    |
| `s3_cluster_kms_request_error`   | Number of KMS requests that failed due to some error. (HTTP 4xx status code).            |
| `s3_cluster_kms_request_failure` | Number of KMS requests that failed due to some internal failure. (HTTP 5xx status code). |
| `s3_cluster_kms_request_success` | Number of KMS requests that succeeded.                                                   |
| `s3_cluster_kms_uptime`          | The time the KMS has been up and running in seconds.                                     |

## Cluster Health Metrics

| Name                                              | Description                                    |
|:--------------------------------------------------|:-----------------------------------------------|
| `s3_cluster_nodes_offline_total`               | Total number of S3 nodes offline.           |
| `s3_cluster_nodes_online_total`                | Total number of S3 nodes online.            |
| `s3_cluster_write_quorum`                      | Maximum write quorum across all pools and sets |
| `s3_cluster_health_status`                     | Get current cluster health status              |
| `s3_cluster_health_erasure_set_healing_drives` | Count of healing drives in the erasure set     |
| `s3_cluster_health_erasure_set_online_drives`  | Count of online drives in the erasure set      |
| `s3_cluster_health_erasure_set_read_quorum`    | Get read quorum of the erasure set             |
| `s3_cluster_health_erasure_set_write_quorum`   | Get write quorum of the erasure set            |
| `s3_cluster_health_erasure_set_status`         | Get current health status of the erasure set   |

## Cluster Replication Metrics

Metrics marked as ``Site Replication Only`` only populate on deployments with [Site Replication](https://docs.hanzo.ai/storage/operations/multi-site-replication) configurations.
For deployments with [bucket](https://docs.hanzo.ai/storage/administration/bucket-replication) or [batch](https://docs.hanzo.ai/storage/administration/batch-framework#replicate) configurations, these metrics populate instead under the [Bucket Metrics](#bucket-metrics) endpoint.

| Name                                                       | Description                                                                                             
|:-----------------------------------------------------------|:---------------------------------------------------------------------------------------------------------|
| `s3_cluster_replication_last_hour_failed_bytes`         | (_Site Replication Only_) Total number of bytes failed at least once to replicate in the last full hour. |
| `s3_cluster_replication_last_hour_failed_count`         | (_Site Replication Only_) Total number of objects which failed replication in the last full hour.        |
| `s3_cluster_replication_last_minute_failed_bytes`       | Total number of bytes failed at least once to replicate in the last full minute.                         |
| `s3_cluster_replication_last_minute_failed_count`       | Total number of objects which failed replication in the last full minute.                                |
| `s3_cluster_replication_total_failed_bytes`             | (_Site Replication Only_) Total number of bytes failed at least once to replicate since server start.    |
| `s3_cluster_replication_total_failed_count`             | (_Site Replication Only_) Total number of objects which failed replication since server start.           |
| `s3_cluster_replication_received_bytes`                 | (_Site Replication Only_) Total number of bytes replicated to this cluster from another source cluster.  |
| `s3_cluster_replication_received_count`                 | (_Site Replication Only_) Total number of objects received by this cluster from another source cluster.  |
| `s3_cluster_replication_sent_bytes`                     | (_Site Replication Only_) Total number of bytes replicated to the target cluster.                        |
| `s3_cluster_replication_sent_count`                     | (_Site Replication Only_) Total number of objects replicated to the target cluster.                      |
| `s3_cluster_replication_credential_errors`              | (_Site Replication Only_) Total number of replication credential errors since server start               |
| `s3_cluster_replication_proxied_get_requests_total` | (_Site Replication Only_)Number of GET requests proxied to replication target                          |
| `s3_cluster_replication_proxied_head_requests_total` | (_Site Replication Only_)Number of HEAD requests proxied to replication target                          |
| `s3_cluster_replication_proxied_delete_tagging_requests_total` | (_Site Replication Only_)Number of DELETE tagging requests proxied to replication target                          |
| `s3_cluster_replication_proxied_get_tagging_requests_total` | (_Site Replication Only_)Number of GET tagging requests proxied to replication target                          |
| `s3_cluster_replication_proxied_put_tagging_requests_total` | (_Site Replication Only_)Number of PUT tagging requests proxied to replication target                          |
| `s3_cluster_replication_proxied_get_requests_failures` | (_Site Replication Only_)Number of failures in GET requests proxied to replication target                          |
| `s3_cluster_replication_proxied_head_requests_failures` | (_Site Replication Only_)Number of failures in HEAD requests proxied to replication target                          |
| `s3_cluster_replication_proxied_delete_tagging_requests_failures` | (_Site Replication Only_)Number of failures proxying DELETE tagging requests to replication target                          |
| `s3_cluster_replication_proxied_get_tagging_requests_failures` | (_Site Replication Only_)Number of failures proxying GET tagging requests to replication target                          |
| `s3_cluster_replication_proxied_put_tagging_requests_failures` | (_Site Replication Only_)Number of failures proxying PUT tagging requests to replication target                          |


## Node Replication Metrics

Metrics marked as ``Site Replication Only`` only populate on deployments with [Site Replication](https://docs.hanzo.ai/storage/operations/multi-site-replication) configurations.
For deployments with [bucket](https://docs.hanzo.ai/storage/administration/bucket-replication) or [batch](https://docs.hanzo.ai/storage/administration/batch-framework#replicate) configurations, these metrics populate instead under the [Bucket Metrics](#bucket-metrics) endpoint.

| Name                                                       | Description
|:-----------------------------------------------------------|:---------------------------------------------------------------------------------------------------------|
| `s3_node_replication_current_active_workers`         | Total number of active replication workers                                                               |
| `s3_node_replication_average_active_workers`         | Average number of active replication workers                                                             |
| `s3_node_replication_max_active_workers`             | Maximum number of active replication workers seen since server start                                     |
| `s3_node_replication_link_online`                    | Reports whether the replication link is online (1) or offline (0).                                       |
| `s3_node_replication_link_offline_duration_seconds`  | Total duration of replication link being offline in seconds since last offline event                     |
| `s3_node_replication_link_downtime_duration_seconds` | Total downtime of replication link in seconds since server start                                         |
| `s3_node_replication_average_link_latency_ms`        | Average replication link latency in milliseconds                                                         |
| `s3_node_replication_max_link_latency_ms`            | Maximum replication link latency in milliseconds seen since server start                                 |
| `s3_node_replication_current_link_latency_ms`        | Current replication link latency in milliseconds                                                         |
| `s3_node_replication_current_transfer_rate`          | Current replication transfer rate in bytes/sec                                                           |
| `s3_node_replication_average_transfer_rate`          | Average replication transfer rate in bytes/sec                                                           |
| `s3_node_replication_max_transfer_rate`              | Maximum replication transfer rate in bytes/sec seen since server start                                   |
| `s3_node_replication_last_minute_queued_count`       | Total number of objects queued for replication in the last full minute                                   |
| `s3_node_replication_last_minute_queued_bytes`       | Total number of bytes queued for replication in the last full minute                                     |
| `s3_node_replication_average_queued_count`           | Average number of objects queued for replication since server start                                      |
| `s3_node_replication_average_queued_bytes`           | Average number of bytes queued for replication since server start                                        |
| `s3_node_replication_max_queued_bytes`               | Maximum number of bytes queued for replication seen since server start                                   |
| `s3_node_replication_max_queued_count`               | Maximum number of objects queued for replication seen since server start                                 |
| `s3_node_replication_recent_backlog_count`           | Total number of objects seen in replication backlog in the last 5 minutes                                |

## Healing Metrics

| Name                                         | Description                                                      |
|:---------------------------------------------|:-----------------------------------------------------------------|
| `s3_heal_objects_errors_total`            | Objects for which healing failed in current self healing run.    |
| `s3_heal_objects_heal_total`              | Objects healed in current self healing run.                      |
| `s3_heal_objects_total`                   | Objects scanned in current self healing run.                     |
| `s3_heal_time_last_activity_nano_seconds` | Time elapsed (in nano seconds) since last self healing activity. |

## Inter Node Metrics

| Name                                      | Description                                             |
|:------------------------------------------|:--------------------------------------------------------|
| `s3_inter_node_traffic_dial_avg_time`  | Average time of internodes TCP dial calls.              |
| `s3_inter_node_traffic_dial_errors`    | Total number of internode TCP dial timeouts and errors. |
| `s3_inter_node_traffic_errors_total`   | Total number of failed internode calls.                 |
| `s3_inter_node_traffic_received_bytes` | Total number of bytes received from other peer nodes.   |
| `s3_inter_node_traffic_sent_bytes`     | Total number of bytes sent to the other peer nodes.     |

## Bucket Notification Metrics

| Name                                           | Description                                                                                                                                 |
|:-----------------------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------|
| `s3_notify_current_send_in_progress`        | Number of concurrent async Send calls active to all targets (deprecated, please use `s3_notify_target_current_send_in_progress` instead) |
| `s3_notify_events_errors_total`             | Events that were failed to be sent to the targets (deprecated, please use `s3_notify_target_failed_events` instead)                      |
| `s3_notify_events_sent_total`               | Total number of events sent to the targets (deprecated, please use `s3_notify_target_total_events` instead)                              |
| `s3_notify_events_skipped_total`            | Events that were skipped to be sent to the targets due to the in-memory queue being full                                                    |
| `s3_notify_target_current_send_in_progress` | Number of concurrent async Send calls active to the target                                                                                  |
| `s3_notify_target_queue_length`             | Number of events currently staged in the queue_dir configured for the target.                                                               |
| `s3_notify_target_total_events`             | Total number of events sent (or) queued to the target                                                                                       |

## S3 API Request Metrics

| Name                                          | Description                                              |
|:----------------------------------------------|:---------------------------------------------------------|
| `s3_s3_requests_4xx_errors_total`          | Total number S3 requests with (4xx) errors.              |
| `s3_s3_requests_5xx_errors_total`          | Total number S3 requests with (5xx) errors.              |
| `s3_s3_requests_canceled_total`            | Total number S3 requests canceled by the client.         |
| `s3_s3_requests_errors_total`              | Total number S3 requests with (4xx and 5xx) errors.      |
| `s3_s3_requests_incoming_total`            | Volatile number of total incoming S3 requests.           |
| `s3_s3_requests_inflight_total`            | Total number of S3 requests currently in flight.         |
| `s3_s3_requests_rejected_auth_total`       | Total number S3 requests rejected for auth failure.      |
| `s3_s3_requests_rejected_header_total`     | Total number S3 requests rejected for invalid header.    |
| `s3_s3_requests_rejected_invalid_total`    | Total number S3 invalid requests.                        |
| `s3_s3_requests_rejected_timestamp_total`  | Total number S3 requests rejected for invalid timestamp. |
| `s3_s3_requests_total`                     | Total number S3 requests.                                |
| `s3_s3_requests_waiting_total`             | Number of S3 requests in the waiting queue.              |
| `s3_s3_requests_ttfb_seconds_distribution` | Distribution of the time to first byte across API calls. |
| `s3_s3_traffic_received_bytes`             | Total number of s3 bytes received.                       |
| `s3_s3_traffic_sent_bytes`                 | Total number of s3 bytes sent.                           |

## Software Metrics

| Name                          | Description                            |
|:------------------------------|:---------------------------------------|
| `s3_software_commit_info`  | Git commit hash for the S3 release. |
| `s3_software_version_info` | S3 release tag for the server.      |

## Drive Metrics

| Name                                   | Description                                                         |
|:---------------------------------------|:--------------------------------------------------------------------|
| `s3_node_drive_free_bytes`          | Total storage available on a drive.                                 |
| `s3_node_drive_free_inodes`         | Total free inodes.                                                  |
| `s3_node_drive_latency_us`          | Average last minute latency in µs for drive API storage operations. |
| `s3_node_drive_offline_total`       | Total drives offline in this node.                                  |
| `s3_node_drive_online_total`        | Total drives online in this node.                                   |
| `s3_node_drive_total`               | Total drives in this node.                                          |
| `s3_node_drive_total_bytes`         | Total storage on a drive.                                           |
| `s3_node_drive_used_bytes`          | Total storage used on a drive.                                      |
| `s3_node_drive_errors_timeout`      | Total number of drive timeout errors since server start             |
| `s3_node_drive_errors_ioerror`      | Total number of drive I/O errors since server start                 |
| `s3_node_drive_errors_availability` | Total number of drive I/O errors, timeouts since server start       |
| `s3_node_drive_io_waiting`          | Total number I/O operations waiting on drive                        |

## Identity and Access Management (IAM) Metrics

| Name                                       | Description                                                 |
|:-------------------------------------------|:------------------------------------------------------------|
| `s3_node_iam_last_sync_duration_millis` | Last successful IAM data sync duration in milliseconds.     |
| `s3_node_iam_since_last_sync_millis`    | Time (in milliseconds) since last successful IAM data sync. |
| `s3_node_iam_sync_failures`             | Number of failed IAM data syncs since server start.         |
| `s3_node_iam_sync_successes`            | Number of successful IAM data syncs since server start.     |

## Information Lifecycle Management (ILM) Metrics

| Name                                                         | Description                                                                                                |
|:-------------------------------------------------------------|:-----------------------------------------------------------------------------------------------------------|
| `s3_node_ilm_expiry_pending_tasks`                        | Number of pending ILM expiry tasks in the queue.                                                           |
| `s3_node_ilm_transition_active_tasks`                     | Number of active ILM transition tasks.                                                                     |
| `s3_node_ilm_transition_pending_tasks`                    | Number of pending ILM transition tasks in the queue.                                                       |
| `s3_node_ilm_transition_missed_immediate_tasks`           | Number of missed immediate ILM transition tasks.                                                           |
| `s3_node_ilm_versions_scanned`                            | Total number of object versions checked for ilm actions since server start.                                |
| `s3_node_ilm_action_count_delete_action`                  | Total action outcome of lifecycle checks since server start for deleting object                            |
| `s3_node_ilm_action_count_delete_version_action`          | Total action outcome of lifecycle checks since server start for deleting a version                         |
| `s3_node_ilm_action_count_transition_action`              | Total action outcome of lifecycle checks since server start for transition of an object                    |
| `s3_node_ilm_action_count_transition_version_action`      | Total action outcome of lifecycle checks since server start for transition of a particular object version  |
| `s3_node_ilm_action_count_delete_restored_action`         | Total action outcome of lifecycle checks since server start for deletion of temporarily restored object    |
| `s3_node_ilm_action_count_delete_restored_version_action` | Total action outcome of lifecycle checks since server start for deletion of a temporarily restored version |
| `s3_node_ilm_action_count_delete_all_versions_action`     | Total action outcome of lifecycle checks since server start for deletion of all versions                   |

## Tier Metrics

| Name                                               | Description                                                                 |
|:---------------------------------------------------|:----------------------------------------------------------------------------|
| `s3_node_tier_tier_ttlb_seconds_distribution`   | Distribution of time to last byte for objects downloaded from warm tier     |
| `s3_node_tier_requests_success`                 | Number of requests to download object from warm tier that were successful   | 
| `s3_node_tier_requests_failure`                 | Number of requests to download object from warm tier that were failure      | 

## System Metrics

| Name                                       | Description                                                                                                     |
|:-------------------------------------------|:----------------------------------------------------------------------------------------------------------------|
| `s3_node_file_descriptor_limit_total`   | Limit on total number of open file descriptors for the S3 server process.                                    |
| `s3_node_file_descriptor_open_total`    | Total number of open file descriptors by the S3 server process.                                              |
| `s3_node_go_routine_total`              | Total number of go routines running.                                                                            |
| `s3_node_io_rchar_bytes`                | Total bytes read by the process from the underlying storage system including cache, /proc/[pid]/io rchar.       |
| `s3_node_io_read_bytes`                 | Total bytes read by the process from the underlying storage system, /proc/[pid]/io read_bytes.                  |
| `s3_node_io_wchar_bytes`                | Total bytes written by the process to the underlying storage system including page cache, /proc/[pid]/io wchar. |
| `s3_node_io_write_bytes`                | Total bytes written by the process to the underlying storage system, /proc/[pid]/io write_bytes.                |
| `s3_node_process_cpu_total_seconds`     | Total user and system CPU time spent in seconds by the process.                                                |
| `s3_node_process_resident_memory_bytes` | Resident memory size in bytes.                                                                                  |
| `s3_node_process_virtual_memory_bytes`  | Virtual memory size in bytes.                                                                                   |
| `s3_node_process_starttime_seconds`     | Start time for S3 process per node, time in seconds since Unix epoc.                                         |
| `s3_node_process_uptime_seconds`        | Uptime for S3 process per node in seconds.                                                                   |

## Scanner Metrics

| Name                                       | Description                                                 |
|:-------------------------------------------|:------------------------------------------------------------|
| `s3_node_scanner_bucket_scans_finished` | Total number of bucket scans finished since server start.   |
| `s3_node_scanner_bucket_scans_started`  | Total number of bucket scans started since server start.    |
| `s3_node_scanner_directories_scanned`   | Total number of directories scanned since server start.     |
| `s3_node_scanner_objects_scanned`       | Total number of unique objects scanned since server start.  |
| `s3_node_scanner_versions_scanned`      | Total number of object versions scanned since server start. |
| `s3_node_syscall_read_total`            | Total read SysCalls to the kernel. /proc/[pid]/io syscr.    |
| `s3_node_syscall_write_total`           | Total write SysCalls to the kernel. /proc/[pid]/io syscw.   |
| `s3_usage_last_activity_nano_seconds`   | Time elapsed (in nano seconds) since last scan activity.    |

# Bucket Metrics

S3 collects the following metrics at the bucket level.
Each metric includes the ``bucket`` label to identify the corresponding bucket.
Metrics may include one or more additional labels, such as the server that calculated that metric.

These metrics can be obtained from any S3 server once per collection by using the following URL:

```shell
https://HOSTNAME:PORT/minio/v2/metrics/bucket
```

Replace ``HOSTNAME:PORT`` with the hostname of your S3 deployment.
For deployments behind a load balancer, use the load balancer hostname instead of a single node hostname.

## Distribution Metrics

| Name                                        | Description                                                                     |
|:--------------------------------------------|:--------------------------------------------------------------------------------|
| `s3_bucket_objects_size_distribution`    | Distribution of object sizes in the bucket, includes label for the bucket name. |
| `s3_bucket_objects_version_distribution` | Distribution of object sizes in a bucket, by number of versions                 |

## Replication Metrics

These metrics only populate on deployments with [Bucket Replication](https://docs.hanzo.ai/storage/administration/bucket-replication) or [Batch Replication](https://docs.hanzo.ai/storage/administration/batch-framework) configurations.
For deployments with [Site Replication](https://docs.hanzo.ai/storage/operations/multi-site-replication) configured, select metrics populate under the [Cluster Metrics](#cluster-metrics) endpoint.

| Name                                                | Description                                                                      |
|:----------------------------------------------------|:---------------------------------------------------------------------------------|
| `s3_bucket_replication_last_minute_failed_bytes` | Total number of bytes failed at least once to replicate in the last full minute. |
| `s3_bucket_replication_last_minute_failed_count` | Total number of objects which failed replication in the last full minute.        |
| `s3_bucket_replication_last_hour_failed_bytes`   | Total number of bytes failed at least once to replicate in the last full hour.   |
| `s3_bucket_replication_last_hour_failed_count`   | Total number of objects which failed replication in the last full hour.          |
| `s3_bucket_replication_total_failed_bytes`       | Total number of bytes failed at least once to replicate since server start.      |
| `s3_bucket_replication_total_failed_count`       | Total number of objects which failed replication since server start.             |
| `s3_bucket_replication_latency_ms`               | Replication latency in milliseconds.                                             |
| `s3_bucket_replication_received_bytes`           | Total number of bytes replicated to this bucket from another source bucket.      |
| `s3_bucket_replication_received_count`           | Total number of objects received by this bucket from another source bucket.      |
| `s3_bucket_replication_sent_bytes`               | Total number of bytes replicated to the target bucket.                           |
| `s3_bucket_replication_sent_count`               | Total number of objects replicated to the target bucket.                         |
| `s3_bucket_replication_credential_errors`        | Total number of replication credential errors since server start                 |
| `s3_bucket_replication_proxied_get_requests_total` | Number of GET requests proxied to replication target                          |
| `s3_bucket_replication_proxied_head_requests_total` | Number of HEAD requests proxied to replication target                          |
| `s3_bucket_replication_proxied_delete_tagging_requests_total` | Number of DELETE tagging requests proxied to replication target                          |
| `s3_bucket_replication_proxied_get_tagging_requests_total` | Number of GET tagging requests proxied to replication target                          |
| `s3_bucket_replication_proxied_put_tagging_requests_total` | Number of PUT tagging requests proxied to replication target                          |
| `s3_bucket_replication_proxied_get_requests_failures` | Number of failures in GET requests proxied to replication target                          |
| `s3_bucket_replication_proxied_head_requests_failures` | Number of failures in HEAD requests proxied to replication target                          |
| `s3_bucket_replication_proxied_delete_tagging_requests_failures` | Number of failures in DELETE tagging proxy requests to replication target                          |
| `s3_bucket_replication_proxied_get_tagging_requests_failures` |Number of failures in GET tagging proxy requests to replication target                          |
| `s3_bucket_replication_proxied_put_tagging_requests_failures` | Number of failures in PUT tagging proxy requests to replication target                          |

## Traffic Metrics

| Name                                  | Description                                        |
|:--------------------------------------|:---------------------------------------------------|
| `s3_bucket_traffic_received_bytes` | Total number of S3 bytes received for this bucket. |
| `s3_bucket_traffic_sent_bytes`     | Total number of S3 bytes sent for this bucket.     |
	
## Usage Metrics

| Name                                    | Description                                       |
|:----------------------------------------|:--------------------------------------------------|
| `s3_bucket_usage_object_total`       | Total number of objects.                          |
| `s3_bucket_usage_version_total`      | Total number of versions (includes delete marker) |
| `s3_bucket_usage_deletemarker_total` | Total number of delete markers.                   |
| `s3_bucket_usage_total_bytes`        | Total bucket size in bytes.                       |
| `s3_bucket_quota_total_bytes`        | Total bucket quota size in bytes.                 |

## Requests Metrics

| Name                                              | Description                                                     |
|:--------------------------------------------------|:----------------------------------------------------------------|
| `s3_bucket_requests_4xx_errors_total`          | Total number of S3 requests with (4xx) errors on a bucket.      |
| `s3_bucket_requests_5xx_errors_total`          | Total number of S3 requests with (5xx) errors on a bucket.      |
| `s3_bucket_requests_inflight_total`            | Total number of S3 requests currently in flight on a bucket.    |
| `s3_bucket_requests_total`                     | Total number of S3 requests on a bucket.                        |
| `s3_bucket_requests_canceled_total`            | Total number S3 requests canceled by the client.                |
| `s3_bucket_requests_ttfb_seconds_distribution` | Distribution of time to first byte across API calls per bucket. |

# Resource Metrics

S3 collects the following resource metrics at the node level.
Each metric includes the `server` label to identify the corresponding node.
Metrics may include one or more additional labels, such as the drive path, interface name, etc.

These metrics can be obtained from any S3 server once per collection by using the following URL:

```shell
https://HOSTNAME:PORT/minio/v2/metrics/resource
```

Replace `HOSTNAME:PORT` with the hostname of your S3 deployment.
For deployments behind a load balancer, use the load balancer hostname instead of a single node hostname.

## Drive Resource Metrics

| Name                                 | Description                                              |
|:-------------------------------------|:---------------------------------------------------------|
| `s3_node_drive_total_bytes`       | Total bytes on a drive.                                  |
| `s3_node_drive_used_bytes`        | Used bytes on a drive.                                   |
| `s3_node_drive_total_inodes`      | Total inodes on a drive.                                 |
| `s3_node_drive_used_inodes`       | Total inodes used on a drive.                            |
| `s3_node_drive_reads_per_sec`     | Reads per second on a drive.                             |
| `s3_node_drive_reads_kb_per_sec`  | Kilobytes read per second on a drive.                    |
| `s3_node_drive_reads_await`       | Average time for read requests to be served on a drive.  |
| `s3_node_drive_writes_per_sec`    | Writes per second on a drive.                            |
| `s3_node_drive_writes_kb_per_sec` | Kilobytes written per second on a drive.                 |
| `s3_node_drive_writes_await`      | Average time for write requests to be served on a drive. |
| `s3_node_drive_perc_util`         | Percentage of time the disk was busy since uptime.       |

## Network Interface Metrics

| Name                          | Description                                                |
|:------------------------------|:-----------------------------------------------------------|
| `s3_node_if_rx_bytes`      | Bytes received on the interface in 60s.                    |
| `s3_node_if_rx_bytes_avg`  | Bytes received on the interface in 60s (avg) since uptime. |
| `s3_node_if_rx_bytes_max`  | Bytes received on the interface in 60s (max) since uptime. |
| `s3_node_if_rx_errors`     | Receive errors in 60s.                                     |
| `s3_node_if_rx_errors_avg` | Receive errors in 60s (avg).                               |
| `s3_node_if_rx_errors_max` | Receive errors in 60s (max).                               |
| `s3_node_if_tx_bytes`      | Bytes transmitted in 60s.                                  |
| `s3_node_if_tx_bytes_avg`  | Bytes transmitted in 60s (avg).                            |
| `s3_node_if_tx_bytes_max`  | Bytes transmitted in 60s (max).                            |
| `s3_node_if_tx_errors`     | Transmit errors in 60s.                                    |
| `s3_node_if_tx_errors_avg` | Transmit errors in 60s (avg).                              |
| `s3_node_if_tx_errors_max` | Transmit errors in 60s (max).                              |

## CPU Metrics

| Name                                 | Description                                |
|:-------------------------------------|:-------------------------------------------|
| `s3_node_cpu_avg_user`            | CPU user time.                             |
| `s3_node_cpu_avg_user_avg`        | CPU user time (avg).                       |
| `s3_node_cpu_avg_user_max`        | CPU user time (max).                       |
| `s3_node_cpu_avg_system`          | CPU system time.                           |
| `s3_node_cpu_avg_system_avg`      | CPU system time (avg).                     |
| `s3_node_cpu_avg_system_max`      | CPU system time (max).                     |
| `s3_node_cpu_avg_idle`            | CPU idle time.                             |
| `s3_node_cpu_avg_idle_avg`        | CPU idle time (avg).                       |
| `s3_node_cpu_avg_idle_max`        | CPU idle time (max).                       |
| `s3_node_cpu_avg_iowait`          | CPU ioWait time.                           |
| `s3_node_cpu_avg_iowait_avg`      | CPU ioWait time (avg).                     |
| `s3_node_cpu_avg_iowait_max`      | CPU ioWait time (max).                     |
| `s3_node_cpu_avg_nice`            | CPU nice time.                             |
| `s3_node_cpu_avg_nice_avg`        | CPU nice time (avg).                       |
| `s3_node_cpu_avg_nice_max`        | CPU nice time (max).                       |
| `s3_node_cpu_avg_steal`           | CPU steam time.                            |
| `s3_node_cpu_avg_steal_avg`       | CPU steam time (avg).                      |
| `s3_node_cpu_avg_steal_max`       | CPU steam time (max).                      |
| `s3_node_cpu_avg_load1`           | CPU load average 1min.                     |
| `s3_node_cpu_avg_load1_avg`       | CPU load average 1min (avg).               |
| `s3_node_cpu_avg_load1_max`       | CPU load average 1min (max).               |
| `s3_node_cpu_avg_load1_perc`      | CPU load average 1min (percentage).        |
| `s3_node_cpu_avg_load1_perc_avg`  | CPU load average 1min (percentage) (avg).  |
| `s3_node_cpu_avg_load1_perc_max`  | CPU load average 1min (percentage) (max).  |
| `s3_node_cpu_avg_load5`           | CPU load average 5min.                     |
| `s3_node_cpu_avg_load5_avg`       | CPU load average 5min (avg).               |
| `s3_node_cpu_avg_load5_max`       | CPU load average 5min (max).               |
| `s3_node_cpu_avg_load5_perc`      | CPU load average 5min (percentage).        |
| `s3_node_cpu_avg_load5_perc_avg`  | CPU load average 5min (percentage) (avg).  |
| `s3_node_cpu_avg_load5_perc_max`  | CPU load average 5min (percentage) (max).  |
| `s3_node_cpu_avg_load15`          | CPU load average 15min.                    |
| `s3_node_cpu_avg_load15_avg`      | CPU load average 15min (avg).              |
| `s3_node_cpu_avg_load15_max`      | CPU load average 15min (max).              |
| `s3_node_cpu_avg_load15_perc`     | CPU load average 15min (percentage).       |
| `s3_node_cpu_avg_load15_perc_avg` | CPU load average 15min (percentage) (avg). |
| `s3_node_cpu_avg_load15_perc_max` | CPU load average 15min (percentage) (max). |

## Memory Metrics

| Name                           | Description                               |
|:-------------------------------|:------------------------------------------|
| `s3_node_mem_available`     | Available memory on the node.             |
| `s3_node_mem_available_avg` | Available memory on the node (avg).       |
| `s3_node_mem_available_max` | Available memory on the node (max).       |
| `s3_node_mem_buffers`       | Buffers memory on the node.               |
| `s3_node_mem_buffers_avg`   | Buffers memory on the node (avg).         |
| `s3_node_mem_buffers_max`   | Buffers memory on the node (max).         |
| `s3_node_mem_cache`         | Cache memory on the node.                 |
| `s3_node_mem_cache_avg`     | Cache memory on the node (avg).           |
| `s3_node_mem_cache_max`     | Cache memory on the node (max).           |
| `s3_node_mem_free`          | Free memory on the node.                  |
| `s3_node_mem_free_avg`      | Free memory on the node (avg).            |
| `s3_node_mem_free_max`      | Free memory on the node (max).            |
| `s3_node_mem_shared`        | Shared memory on the node.                |
| `s3_node_mem_shared_avg`    | Shared memory on the node (avg).          |
| `s3_node_mem_shared_max`    | Shared memory on the node (max).          |
| `s3_node_mem_total`         | Total memory on the node.                 |
| `s3_node_mem_total_avg`     | Total memory on the node (avg).           |
| `s3_node_mem_total_max`     | Total memory on the node (max).           |
| `s3_node_mem_used`          | Used memory on the node.                  |
| `s3_node_mem_used_avg`      | Used memory on the node (avg).            |
| `s3_node_mem_used_max`      | Used memory on the node (max).            |
| `s3_node_mem_used_perc`     | Used memory percentage on the node.       |
| `s3_node_mem_used_perc_avg` | Used memory percentage on the node (avg). |
| `s3_node_mem_used_perc_max` | Used memory percentage on the node (max). |
