# Deployment

**ldap-image-proxy** can run as:

1. a standalone binary;
2. a container;
3. a Docker Compose stack.

## Published images

- `icikowski/ldap-image-proxy`
- `ghcr.io/icikowski/ldap-image-proxy`

Published tags are multi-arch manifests that resolve to `linux/amd64` and `linux/arm64`.

### Tag guide

| Tag      | Published from                | Notes                                       |
| -------- | ----------------------------- | ------------------------------------------- |
| `dev`    | `main`                        | Published on every push to the main branch. |
| `v*`     | Matching Git `v*` tag         | Published for each versioned release tag.   |
| `latest` | Latest built `v*` release tag | Points to the newest built release tag.     |

## Reverse proxy for TLS/SSL

**ldap-image-proxy** does not terminate TLS/SSL itself. Put a reverse proxy (such as `nginx`) in front of it and handle HTTPS there.

## Requirements

Every deployment needs:

- LDAP connection settings;
- a bind user and password;
- the avatar attribute name;
- a reverse proxy if you want HTTPS.

## Standalone binary

### Build

```bash
task build-static
```

### Run

```bash
LIP_LDAP_URL="ldap://lldap:3890" \
LIP_LDAP_BASE_DN="dc=example,dc=com" \
LIP_LDAP_BIND_USER="uid=service_account,ou=people,dc=example,dc=com" \
LIP_LDAP_BIND_PASSWORD="replace-me" \
LIP_LDAP_IMPLEMENTATION="lldap" \
./bin/ldap-image-proxy
```

## Docker

The image listens on `3891` by default and uses Docker `HEALTCHECK`.

```bash
docker run -d \
  --name ldap-image-proxy \
  -p 3891:3891 \
  --tmpfs /tmp:rw,noexec,nosuid,size=64m \
  --env-file .env \
  icikowski/ldap-image-proxy:<tag>
```

## Docker Compose

Minimal Compose example:

```yaml
services:
  ldap-image-proxy:
    image: icikowski/ldap-image-proxy:<tag>
    restart: unless-stopped
    env_file:
      - .env
    tmpfs:
      - /tmp:rw,noexec,nosuid,size=64m
    ports:
      - "3891:3891"
```

### Example `.env`

```dotenv
LIP_LDAP_URL="ldap://lldap:3890"
LIP_LDAP_BASE_DN="dc=example,dc=com"
LIP_LDAP_BIND_USER="uid=service_account,ou=people,dc=example,dc=com"
LIP_LDAP_BIND_PASSWORD="replace-me"
LIP_LDAP_IMPLEMENTATION="lldap"
```

### Compose notes

- The application image is distroless.
- In this repository, `.env` is a documented workflow file loaded by the committed `.envrc`.

## Nginx reverse proxy

Mount the service behind Nginx and expose it (eg. under a subpath):

```nginx
upstream ldapimageproxy {
    server 127.0.0.1:3891;
}

server {
    listen 80;
    listen [::]:80;
    server_name example.com;

    location / {
        include /etc/nginx/snippets/reverse-proxy.conf;
        proxy_pass http://ldapimageproxy/;
    }
}
```

Or, if you want to expose it under a subpath:

```nginx
location /.avatars/ {
    include /etc/nginx/snippets/reverse-proxy.conf;
    proxy_pass http://ldapimageproxy/;
}
```

That makes avatars available at `https://example.com/.avatars/{username}`.

## See also

- [Environment variables](environment-variables.md)
- [LDAP configuration](ldap.md)
