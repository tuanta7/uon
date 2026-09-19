SHELL := /bin/bash

UBUNTU_IMAGE ?= ubuntu:24.04
TEST_COMMAND ?= uon --help
DOCKER_ARGS ?=

.PHONY: build test

build:
	go build -o uon .

# Run the requested command as root in a fresh Ubuntu container. Only the
# freshly-built binary is mounted, so the container starts from a clean image
# on every invocation and is removed when the command exits.
# Example: make test TEST_COMMAND="uon system"
test: build
	docker run --rm --pull=missing $(DOCKER_ARGS) \
		--volume "$(CURDIR)/uon:/usr/local/bin/uon:ro" \
		--workdir /tmp \
		$(UBUNTU_IMAGE) \
		/bin/bash -lc '$(TEST_COMMAND)'
