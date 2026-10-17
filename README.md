# LDAP Image Proxy

**A proxy service for serving LDAP user avatars as images.**

`ldap-image-proxy` serves LDAP user images at a simple URL like `https://example.com/.avatars/{username}`. It was built so LDAP avatar bytes can be exposed as an HTTP image source for services using OIDC's `picture` claim.

It does one thing:

- looks up a user by username or email in LDAP;
- reads the configured avatar attribute;
- serves the image with a cache TTL (`Cache-Control`) and a local `/tmp` cache;
- optionally resizes the image with a Gravatar-style `?s=<pixels>` query parameter, clamped to 16-2048 pixels with a default of 256, or returns the original image with `?s=orig` / `?s=original`.

It serves HTTP only - **terminate TLS at a reverse proxy in front of it**.

> **Warning**
>
> This repository is mirrored to GitHub. The SourceHut repository is the primary one.

## Overview

`ldap-image-proxy` is meant for self-hosted setups that have:

- LDAP users with avatar bytes stored in an attribute;
- a need for a stable image URL.

## Repository layout

| Path                   | Description                                   |
| ---------------------- | --------------------------------------------- |
| `cmd/ldap-image-proxy` | Main application entrypoint                   |
| `cmd/healthcheck`      | Docker-specific healthcheck entrypoint        |
| `internal/config`      | Environment variable parsing and validation   |
| `internal/ldap`        | LDAP client and search filter handling        |
| `internal/logs`        | Zerolog setup                                 |
| `internal/web`         | HTTP routes, middleware and response handling |

## Documentation

| Topic                                                  | Purpose                                                              |
| ------------------------------------------------------ | -------------------------------------------------------------------- |
| [Configuration overview](docs/configuration.md)        | Runtime model and request flow                                       |
| [Environment variables](docs/environment-variables.md) | Reference for every `LIP_*` variable                                 |
| [Deployment](docs/deployment.md)                       | Binary, Docker, Docker Compose and reverse proxy deployment guidance |
| [LDAP](docs/ldap.md)                                   | LDAP presets, custom filters and avatar setup                        |

## Development

### Prerequisites

- Go 1.27+
- [`task`](https://taskfile.dev)
- [`golangci-lint`](https://golangci-lint.run)

### Common commands

```bash
task generate
task build
task run
task test
task lint
```

### Setup

1. Start from `.env.example` or `docs/environment-variables.md`.
2. For local development, place the required `LIP_*` variables in `.env`; the committed `.envrc` loads it with direnv.

## Contributing

The Go module path is `git.sr.ht/~icikowski/ldap-image-proxy`. Documentation links should treat SourceHut as the primary forge.

When contributing:

1. Run `task generate` after template, asset, or UI source changes.
2. Run `task fmt`, `task lint`, and `task test` before sending changes.
3. Update the docs when configuration, deployment behavior, or schemas change.
4. Keep changes focused and operator-friendly.

## Container images

Published images:

- `icikowski/ldap-image-proxy`
- `ghcr.io/icikowski/ldap-image-proxy`

See [`docs/deployment.md`](docs/deployment.md) for full examples.

## License

Licensed under the GNU AGPL v3. See [`LICENSE`](LICENSE).

## Disclaimer

Code was written in spare time, without AI assistance as I find vibed crap unmaintainable. Contributions are welcome, but every contributor is expected to review and test the code they submit and to be responsive to feedback on it. I will not merge contributions that I have not personally reviewed and tested.

Use at your own risk, and expect a personal project with limited support and no SLA. Contributions are welcome, but I can't promise any specific turnaround time or feature roadmap.
