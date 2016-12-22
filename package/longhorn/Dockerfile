FROM imikushin/longhorn-orc:890410e
RUN apt-get update && \
    apt-get install -y jq
COPY storage common/* longhorn/rancher-longhorn /usr/bin/
CMD ["start.sh", "storage", "--driver-name", "rancher-longhorn"]
