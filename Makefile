SHELL := /bin/bash

UBUNTU_IMAGE ?= ubuntu:24.04
TEST_COMMAND ?= uon --help
DOCKER_ARGS ?=
SYSTEM_IMAGE ?= uon-ubuntu-system:24.04

.PHONY: build test system-image

build:
	go build -o uon .

system-image:
	docker build --build-arg UBUNTU_IMAGE=$(UBUNTU_IMAGE) --tag $(SYSTEM_IMAGE) .

# Open a root shell and remove the container when the shell exits.
test: build system-image
	@set -e; container_id=$$(docker run --rm -d --privileged --cgroupns=private \
		--tmpfs /run --tmpfs /run/lock --tmpfs /tmp $(DOCKER_ARGS) \
		--volume "$(CURDIR)/uon:/usr/local/bin/uon:ro" \
		--workdir /tmp \
		$(SYSTEM_IMAGE)); \
	trap 'docker stop --time 10 "$$container_id" >/dev/null' EXIT; \
	docker exec "$$container_id" timeout 30 bash -c \
		'until systemctl is-active --quiet systemd-networkd 2>/dev/null; do sleep 1; done'; \
	docker exec -it "$$container_id" /bin/bash
