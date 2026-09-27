# UON

UON is a CLI for setting up an Ubuntu server for lightweight self-hosting.

Install the latest version:

```sh
go install github.com/tuanta7/uon@latest
```

## Quick Start

List the server's network interfaces, configure a static address, prevent the
machine from sleeping, enable SSH, and start NGINX:

```sh
uon system network list
uon system network static eth0 --address 192.168.1.10/24 --gateway 192.168.1.1/24
uon system sleep disable
uon system ssh enable
uon nginx install
uon nginx run
uon nginx status
```

Run `uon --help` to list the available command groups and options.
