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
uon nginx load ./example.conf # copy, enable, validate, and reload
uon nginx remove          # remove the package but preserve its configuration
uon nginx remove --prune  # remove the package and its configuration
```

Run `uon nginx config` to see Ubuntu's standard configuration locations and the
commands for enabling, validating, and reloading a site. Define server blocks in
`/etc/nginx/sites-available/` and enable them with symbolic links in
`/etc/nginx/sites-enabled/`.

`uon nginx load <file>` automates that workflow. It copies the file into
`/etc/nginx/sites-available/`, enables it under the same file name, checks the
configuration with `nginx -t`, and reloads NGINX. Existing files and enabled
entries with the same name are replaced.

## Development Infrastructure

From the repository root, start the local MinIO and Redis services with:

```sh
sudo uon dev up
```

MinIO and Redis bind private backend ports on loopback. Expose them through the
host NGINX service:

```sh
sudo uon nginx install
sudo uon nginx run
sudo uon nginx load --stream static/uon.dev.nginx.conf
```

NGINX listens on all IPv4 interfaces, so other computers on the LAN can use:

- MinIO API: `http://<server-LAN-IP>:9000`
- MinIO console: `http://<server-LAN-IP>:9001` (`minio` / `password`)
- Redis: `<server-LAN-IP>:6379`

The Docker backend ports remain bound to `127.0.0.1` and cannot be reached
directly from the LAN. If a host firewall is enabled, allow incoming TCP
traffic on ports `9000`, `9001`, and `6379` only for trusted LAN clients. The
development Redis service does not require authentication.

Stop the services with `sudo uon dev down`. This also deletes the MinIO and
Redis development volumes and all data stored in them.

## Cobra Quick Start

```sh
go mod init uon
go get -u github.com/spf13/cobra@latest

go install github.com/spf13/cobra-cli@latest
cobra-cli init
cobra-cli add serve
```
