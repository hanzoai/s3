FROM alpine:3.20
RUN apk add --no-cache curl && \
    addgroup -S hanzo && \
    adduser -S hanzo -G hanzo
COPY s3 /usr/bin/s3
RUN chmod +x /usr/bin/s3 && \
    chown hanzo:hanzo /usr/bin/s3 && \
    mkdir -p /etc/hanzo /data/filerldb2 && \
    chown -R hanzo:hanzo /etc/hanzo /data && \
    chmod 755 /data /etc/hanzo /data/filerldb2
WORKDIR /data
USER hanzo
