# UON

CLI tool for quickly setting up an Ubuntu server on a laptop or PC, with everything needed for a lightweight self-hosted environment. Includes:

- No sleep & automatic restart
- Static IP configuration
- SSH
- Docker & Kubernetes (K3s)
- Ollama, NGINX, Cloudflared
- Multipass for experimenting with multi-cluster Kubernetes

## Getting Started

```sh
go build -o uon
./uon brew install
```

List all network interfaces and their IP addresses:

```sh
sudo ./uon system
```

Configure a network interface with a static IP or DHCP:

```sh
sudo ./uon system ip --interface eth0 --static 192.168.1.10/24 --gateway 192.168.1.1/24
sudo ./uon system ip --interface eth0 --dhcp
```

The static address and optional gateway must include a CIDR prefix. These
commands apply the network configuration through Netplan.

## Testing in a clean Ubuntu image

Docker is required. The `test` target builds the CLI, mounts only the binary
into a temporary Ubuntu 24.04 container, and runs `uon --help` as root:

```sh
make test
```

Run another command by overriding `TEST_COMMAND`:

```sh
make test TEST_COMMAND='uon brew install'
```

The image can be changed with `UBUNTU_IMAGE`, and extra Docker options can be
passed with `DOCKER_ARGS` (for example, `DOCKER_ARGS=--privileged`).

## Cobra Quick Start

```sh
go mod init uon
go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest
cobra-cli init
cobra-cli add serve
```
