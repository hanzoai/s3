#!/bin/sh

# Enable FIPS 140-3 mode by default (Go 1.24+)
# To disable: docker run -e GODEBUG=fips140=off ...
export GODEBUG="${GODEBUG:+$GODEBUG,}fips140=on"

# Fix permissions for mounted volumes
# If /data is mounted from host, it might have different ownership
# Fix this by ensuring hanzo user owns the directory
if [ "$(id -u)" = "0" ]; then
  # Running as root, check and fix permissions if needed
  SEAWEED_UID=$(id -u hanzo)
  SEAWEED_GID=$(id -g hanzo)
  
  # Verify hanzo user and group exist
  if [ -z "$SEAWEED_UID" ] || [ -z "$SEAWEED_GID" ]; then
    echo "Error: 'hanzo' user or group not found. Cannot fix permissions." >&2
    exit 1
  fi
  
  DATA_UID=$(stat -c '%u' /data 2>/dev/null)
  DATA_GID=$(stat -c '%g' /data 2>/dev/null)

  # Only run chown -R if ownership doesn't already match (avoids expensive
  # recursive chown on subsequent starts, and is a no-op on OpenShift when
  # fsGroup has already set correct ownership on the PVC).
  if [ "$DATA_UID" != "$SEAWEED_UID" ] || [ "$DATA_GID" != "$SEAWEED_GID" ]; then
    echo "Fixing /data ownership for hanzo user (uid=$SEAWEED_UID, gid=$SEAWEED_GID)"
    if ! chown -R hanzo:hanzo /data; then
      echo "Warning: Failed to change ownership of /data. This may cause permission errors." >&2
      echo "If /data is read-only or has mount issues, the application may fail to start." >&2
    fi
  fi
  
  # Use su-exec to drop privileges and run as hanzo user
  exec su-exec hanzo "$0" "$@"
fi

isArgPassed() {
  # Match both `-flag` and `--flag` (and their `=value` forms): the Go fla9
  # library accepts both, and users may pick either form on the CLI.
  arg="$1"
  argWithEqualSign="$1="
  argDouble="-$1"
  argDoubleWithEqualSign="-$1="
  shift
  while [ $# -gt 0 ]; do
    passedArg="$1"
    shift
    case $passedArg in
    "$arg"|"$argDouble")
      return 0
      ;;
    "$argWithEqualSign"*|"$argDoubleWithEqualSign"*)
      return 0
      ;;
    esac
  done
  return 1
}

case "$1" in

  'master')
  	ARGS="-mdir=/data -volumeSizeLimitMB=1024"
  	shift
  	exec /usr/bin/s3 -logtostderr=true master $ARGS $@
	;;

  'volume')
  	ARGS="-dir=/data -max=0"
  	if isArgPassed "-max" "$@"; then
  	  ARGS="-dir=/data"
  	fi
  	shift
  	exec /usr/bin/s3 -logtostderr=true volume $ARGS $@
	;;

  'volume-rust')
  	ARGS="-dir /data -max 0"
  	if isArgPassed "-max" "$@"; then
  	  ARGS="-dir /data"
  	fi
  	shift
  	if [ ! -s /usr/bin/s3-volume ]; then
  	  echo "Error: Rust volume server is not available on this platform ($(uname -m))." >&2
  	  echo "Use 'volume' for the Go volume server instead." >&2
  	  exit 1
  	fi
  	exec /usr/bin/s3-volume $ARGS $@
	;;

  'server')
  	ARGS="-dir=/data -volume.max=0 -master.volumeSizeLimitMB=1024"
  	if isArgPassed "-volume.max" "$@"; then
  	  ARGS="-dir=/data -master.volumeSizeLimitMB=1024"
  	fi
 	shift
  	exec /usr/bin/s3 -logtostderr=true server $ARGS $@
  	;;

  'mini')
  	ARGS="-dir=/data"
  	if isArgPassed "-dir" "$@"; then
  	  ARGS=""
  	fi
  	shift
  	exec /usr/bin/s3 -logtostderr=true mini $ARGS $@
  	;;

  'filer')
  	ARGS=""
  	shift
  	exec /usr/bin/s3 -logtostderr=true filer $ARGS $@
	;;

  's3')
  	ARGS="-domainName=$S3_DOMAIN_NAME -key.file=$S3_KEY_FILE -cert.file=$S3_CERT_FILE"
  	shift
  	exec /usr/bin/s3 -logtostderr=true s3 $ARGS $@
	;;

  'shell')
  	ARGS="-cluster=$SHELL_CLUSTER -filer=$SHELL_FILER -filerGroup=$SHELL_FILER_GROUP -master=$SHELL_MASTER -options=$SHELL_OPTIONS"
  	shift
  	exec echo "$@" | /usr/bin/s3 -logtostderr=true shell $ARGS
  ;;

  *)
  	exec /usr/bin/s3 $@
	;;
esac
