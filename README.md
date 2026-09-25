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
uon system network list
uon system network static eth0 --address 192.168.1.10/24 --gateway 192.168.1.1/24
uon system network dhcp eth0
uon system sleep disable # for Ubuntu GUI
uon system ssh enable
uon nginx install
uon nginx config
uon nginx run
uon nginx status
```

## NGINX

Manage the Ubuntu NGINX package and service with:

```sh
uon nginx install # install with apt-get
uon nginx run     # start now and enable at boot
uon nginx status  # show installed, active, and enabled states
uon nginx remove  # remove the package but preserve its configuration
```

Run `uon nginx config` to see Ubuntu's standard configuration locations and the
commands for enabling, validating, and reloading a site. Define server blocks in
`/etc/nginx/sites-available/` and enable them with symbolic links in
`/etc/nginx/sites-enabled/`.

## Cobra Quick Start

```sh
go mod init uon
go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest
cobra-cli init
cobra-cli add serve
```
