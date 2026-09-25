# UON

CLI tool for quickly setting up an Ubuntu server on a laptop or PC, with everything needed for a lightweight self-hosted environment. Includes:

- No sleep & automatic restart
- Static IP configuration
- SSH
- Docker & Kubernetes (K3s)
- Ollama, NGINX, Cloudflared
- Multipass for experimenting with multi-cluster Kubernetes

## Quick Start

```sh
go install github.com/tuanta7/uon@latest
uon system ip --interface eth0 --static 192.168.1.10/24 --gateway 192.168.1.1/24
uon system sleep off # for Ubuntu GUI
uon system ssh on
```

## Cobra Quick Start

```sh
go mod init uon
go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest
cobra-cli init
cobra-cli add serve
```
