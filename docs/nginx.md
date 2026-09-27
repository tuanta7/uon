# NGINX commands

Use `uon nginx` to install and manage the Ubuntu NGINX package and service.

## Install and run

```sh
uon nginx install
uon nginx run
uon nginx status
```

- `install` installs NGINX with APT.
- `run` starts NGINX immediately and enables it at boot.
- `status` reports whether NGINX is installed, active, and enabled.

## Configuration

Run the following command to display Ubuntu's standard NGINX configuration
locations and the steps for enabling, validating, and reloading a site:

```sh
uon nginx config
```

HTTP server blocks belong in `/etc/nginx/sites-available/` and are enabled by
symbolic links in `/etc/nginx/sites-enabled/`.

Load an HTTP configuration from a local file:

```sh
uon nginx load ./example.conf
```

This copies the file into `/etc/nginx/sites-available/`, enables it under the
same file name, validates the complete configuration with `nginx -t`, and
reloads NGINX. Existing files and enabled entries with the same name are
replaced.

Load a main-context TCP/UDP stream configuration instead:

```sh
uon nginx load --stream static/uon.dev.nginx.conf
```

## Remove

Remove the package while preserving its configuration:

```sh
uon nginx remove
```

Remove both the package and its configuration:

```sh
uon nginx remove --prune
```
