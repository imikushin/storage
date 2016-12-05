FROM ubuntu:16.04
COPY storage /usr/bin/
COPY common/common.sh longhorn/rancher-longhorn common/start.sh /usr/bin/
CMD ["start.sh", "storage", "--driver-name", "rancher-longhorn"]
