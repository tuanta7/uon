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

## Cobra Quick Start

```sh
go mod init uon
go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest
cobra-cli init
cobra-cli add serve
```
