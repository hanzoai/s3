#!/bin/sh
#

# If command starts with an option, prepend minio.
if [ "${1}" != "minio" ]; then
	if [ -n "${1}" ]; then
		set -- minio "$@"
	fi
fi

docker_switch_user() {
	if [ -n "${S3_USERNAME}" ] && [ -n "${S3_GROUPNAME}" ]; then
		if [ -n "${S3_UID}" ] && [ -n "${S3_GID}" ]; then
			chroot --userspec=${S3_UID}:${S3_GID} / "$@"
		else
			echo "${S3_USERNAME}:x:1000:1000:${S3_USERNAME}:/:/sbin/nologin" >>/etc/passwd
			echo "${S3_GROUPNAME}:x:1000" >>/etc/group
			chroot --userspec=${S3_USERNAME}:${S3_GROUPNAME} / "$@"
		fi
	else
		exec "$@"
	fi
}

## DEPRECATED and unsupported - switch to user if applicable.
docker_switch_user "$@"
