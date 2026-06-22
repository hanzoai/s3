#!/usr/bin/env bash
#
# Runs the SMB test batteries inside the samba container against the local smbd,
# which serves /mnt/hanzo/share over a Hanzo FUSE mount. A second FUSE
# mount (/mnt/hanzo2) backs the distributed-locking tests.
# Invoked via: docker compose exec samba /run_inside_container.sh
set -euo pipefail

export SMB_HOST=127.0.0.1
export SMB_SHARE=hanzo
export SMB_PORT="${SMB_PORT:-445}"
export SMB_USER="${SMB_USER:-smbtest}"
export SMB_PASS="${SMB_PASS:-smbtest}"
export SHARE_FS_PATH="${SHARE_FS_PATH:-/mnt/hanzo/share}"
export MOUNT_SHARE="${MOUNT_SHARE:-/mnt/hanzo/share}"
export MOUNT2_SHARE="${MOUNT2_SHARE:-/mnt/hanzo2/share}"

rc=0
echo "############ SMB functional tests ############"
/smb_tests.sh || rc=1
echo
echo "############ SMB locking / concurrency tests ############"
/lock_tests.sh || rc=1
exit "${rc}"
