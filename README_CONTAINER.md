# LDAP Image Proxy

**A proxy service for serving LDAP user avatars as images.**

**LDAP Image Proxy** serves LDAP user images at a simple URL like `https://example.com/.avatars/{username}`. It was built so LDAP avatar bytes can be exposed as an HTTP image source for services using OIDC's `picture` claim.

## Images

- `icikowski/ldap-image-proxy`
- `ghcr.io/icikowski/ldap-image-proxy`

## Required environment variables

- `LIP_LDAP_URL`
- `LIP_LDAP_BASE_DN`
- `LIP_LDAP_BIND_USER`
- `LIP_LDAP_BIND_PASSWORD`

## Common container settings

- `LIP_LDAP_IMPLEMENTATION`
- `LIP_LDAP_ATTRIBUTE_AVATAR`
- `LIP_LDAP_ATTRIBUTE_USERNAME`
- `LIP_LDAP_ATTRIBUTE_EMAIL`
- `LIP_LDAP_SEARCH_FILTER`
- `LIP_SERVER_ADDRESS`
- `LIP_SERVER_PORT`
- `LIP_SERVER_IMAGE_TTL`
- `LIP_LOG_LEVEL`
- `LIP_LOG_PRETTY`

## Docker

```bash
docker run -d \
  --name ldap-image-proxy \
  -p 3891:3891 \
  --tmpfs /tmp:rw,noexec,nosuid,size=64m \
  --env-file .env \
  icikowski/ldap-image-proxy:latest
```

## Docker Compose

```yaml
services:
  ldap-image-proxy:
    image: icikowski/ldap-image-proxy:latest
    restart: unless-stopped
    env_file:
      - .env
    tmpfs:
      - /tmp:rw,noexec,nosuid,size=64m
    ports:
      - "3891:3891"
```
