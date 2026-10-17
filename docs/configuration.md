# Configuration overview

`ldap-image-proxy` reads configuration from `LIP_*` environment variables and LDAP implementation presets.

The application serves HTTP only. Deploy it behind a reverse proxy that terminates TLS/SSL and exposes a stable public path.

## Route

| Route                                       | Purpose                                                      |
| ------------------------------------------- | ------------------------------------------------------------ |
| `/`                                         | Returns `204 No Content` for health checks and root requests |
| `/{username}`                               | Returns avatar with default size of 256 pixels               |
| `/{username}?s=<pixels>` / `?size=<pixels>` | Returns avatar resized to the requested size                 |
| `/{username}?s=orig` / `?size=original`     | Returns the avatar without resizing                          |

Size query parameters are optional. When omitted, the default size of 256 pixels is used. Parameters are clamped to a minimum of 16 and a maximum of 2048 pixels.

## Startup

At startup the application:

1. loads and validates environment variables;
2. resolves LDAP implementation defaults;
3. creates the LDAP client;
4. starts the HTTP server.

Successful startup requires:

- valid LDAP connection settings;
- a reachable LDAP server;
- a valid LDAP implementation choice.

## Request handling

The response for `/{username}` is cached on the client side for the duration configured by `LIP_SERVER_IMAGE_TTL` (via `Cache-Control` header). Resized responses are also cached locally in a temporary directory created at startup for the same duration. Requests that send `Cache-Control: no-cache`, `Cache-Control: no-store`, or `Pragma: no-cache` bypass the local cache and get a `no-store` response instead. Requests may also lower the local cache freshness window with `Cache-Control: max-age=<seconds>`.

A reverse proxy can mount the endpoint under a public subpath such as `/.avatars/`.

## See also

- [Environment variables](environment-variables.md)
- [Deployment](deployment.md)
- [LDAP configuration](ldap.md)
