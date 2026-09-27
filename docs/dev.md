# Development infrastructure commands

Use `uon dev` to manage the repository's local MinIO and Redis services with
Docker Compose. Run these commands from the repository root so UON can find
`static/docker-compose.dev.yml`.

## Start

```sh
uon dev up
```

The services bind their backend ports to the loopback interface. To expose them
to trusted computers on the LAN through the host NGINX service, run:

```sh
uon nginx install
uon nginx run
uon nginx load --stream static/uon.dev.nginx.conf
```

NGINX then listens on all IPv4 interfaces for:

- MinIO API: `http://<server-LAN-IP>:9000`
- MinIO console: `http://<server-LAN-IP>:9001` (`minio` / `password`)
- Redis: `<server-LAN-IP>:6379`

The Docker backend ports remain bound to `127.0.0.1` and cannot be reached
directly from the LAN. If a host firewall is enabled, allow incoming TCP
traffic on ports `9000`, `9001`, and `6379` only for trusted LAN clients. The
development Redis service does not require authentication.

## Stop

```sh
uon dev down
```

This stops the services and deletes their Docker volumes, including all MinIO
and Redis development data.
