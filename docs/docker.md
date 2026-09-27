# Docker commands

Use `uon docker` to install and manage Docker Engine on Ubuntu. Installation
uses Docker's official APT repository and installs Docker Engine, the Docker
CLI, containerd, Buildx, and the Docker Compose plugin.

The workflow follows Docker's
[Ubuntu installation guide](https://docs.docker.com/engine/install/ubuntu/).

Before installation, review Docker's firewall limitations. Published container
ports can bypass `ufw` and `firewalld` rules, and custom filtering rules should
be added to the `DOCKER-USER` chain.

## Install

```sh
uon docker install
```

This command:

1. Removes conflicting Ubuntu packages such as `docker.io`, `containerd`, and
   `runc` when they are installed.
2. Installs `ca-certificates` and `curl`.
3. Adds Docker's signing key and official Ubuntu APT repository.
4. Installs the latest stable Docker Engine packages.

## Run

Start Docker immediately and enable it at boot:

```sh
uon docker run
```

## Remove

Remove Docker Engine packages, its APT source, and its signing key while
preserving images, containers, volumes, and custom configuration:

```sh
uon docker remove
```

Also delete Docker and containerd's stored images, containers, and volumes:

```sh
uon docker remove --prune
```

The `--prune` operation deletes `/var/lib/docker` and `/var/lib/containerd` and
cannot be undone. Edited configuration files outside those directories must be
removed manually.

See [Development infrastructure](dev.md) for the commands that run MinIO and
Redis with Docker Compose.
