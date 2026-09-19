ARG UBUNTU_IMAGE=ubuntu:24.04
FROM ${UBUNTU_IMAGE}

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        systemd systemd-sysv dbus udev netplan.io iproute2 iputils-ping \
        dnsutils curl ca-certificates less nano \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV container=docker
RUN systemctl enable systemd-networkd.service \
    && systemctl mask systemd-networkd-wait-online.service \
    && truncate -s 0 /etc/machine-id

STOPSIGNAL SIGRTMIN+3
CMD ["/sbin/init"]
