# System commands

Use `uon system` to manage network interfaces, system sleep, and SSH access.

## Network

List network interfaces and their assigned IP addresses:

```sh
uon system network list
```

Configure an interface with a static IP address:

```sh
uon system network static eth0 \
  --address 192.168.1.10/24 \
  --gateway 192.168.1.1/24
```

`--address` is required and must include a CIDR prefix. `--gateway` is
optional; when provided, it must also include a CIDR prefix.

Return an interface to DHCP:

```sh
uon system network dhcp eth0
```

## Sleep

Prevent suspend and hibernation, which is useful when a laptop or desktop is
being used as a server:

```sh
uon system sleep disable
```

Allow suspend and hibernation again:

```sh
uon system sleep enable
```

## SSH

Start SSH now and at boot:

```sh
uon system ssh enable
```

If `openssh-server` is not installed, UON asks for confirmation before
installing it.

Stop SSH and disable it at boot:

```sh
uon system ssh disable
```
