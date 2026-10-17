# Environment variables

All runtime configuration is provided through `LIP_*` environment variables.

Durations use Go-style syntax such as `500ms`, `10m`, and `24h`.

## Server

| Variable               | Default | Required | Purpose                                  |
| ---------------------- | ------- | -------- | ---------------------------------------- |
| `LIP_SERVER_ADDRESS`   | empty   | No       | Bind address; empty means all interfaces |
| `LIP_SERVER_PORT`      | `3891`  | No       | HTTP listen port                         |
| `LIP_SERVER_IMAGE_TTL` | `12h`   | No       | Cache TTL for avatar responses           |

## LDAP

| Variable                      | Default     | Required          | Purpose                                                                                                     |
| ----------------------------- | ----------- | ----------------- | ----------------------------------------------------------------------------------------------------------- |
| `LIP_LDAP_URL`                | none        | Yes               | LDAP server URL                                                                                             |
| `LIP_LDAP_BASE_DN`            | none        | Yes               | Base DN                                                                                                     |
| `LIP_LDAP_BIND_USER`          | none        | Yes               | Bind user                                                                                                   |
| `LIP_LDAP_BIND_PASSWORD`      | none        | Yes               | Bind password                                                                                               |
| `LIP_LDAP_IMPLEMENTATION`     | `custom`    | No                | `activedirectory`, `freeipa`, `glauth`, `lldap`, `rfc2307bis`, `custom` (see [LDAP configuration](ldap.md)) |
| `LIP_LDAP_ATTRIBUTE_AVATAR`   | `jpegphoto` | No                | Avatar attribute                                                                                            |
| `LIP_LDAP_ATTRIBUTE_USERNAME` | empty       | Only for `custom` | Username attribute for custom LDAP schemas                                                                  |
| `LIP_LDAP_ATTRIBUTE_EMAIL`    | empty       | Only for `custom` | Mail attribute for custom LDAP schemas                                                                      |
| `LIP_LDAP_SEARCH_FILTER`      | empty       | Only for `custom` | Search filter for custom LDAP schemas                                                                       |

These values are required for the application to start.

## Logging

| Variable         | Default | Required | Purpose             |
| ---------------- | ------- | -------- | ------------------- |
| `LIP_LOG_LEVEL`  | `info`  | No       | Zerolog log level   |
| `LIP_LOG_PRETTY` | `false` | No       | Human-friendly logs |

Pretty logging is useful during development but adds overhead.

## Minimal example

```dotenv
LIP_LDAP_URL="ldap://lldap:3890"
LIP_LDAP_BASE_DN="dc=example,dc=com"
LIP_LDAP_BIND_USER="uid=service_account,ou=people,dc=example,dc=com"
LIP_LDAP_BIND_PASSWORD="replace-me"
LIP_LDAP_IMPLEMENTATION="lldap"
```

## See also

- [Configuration overview](configuration.md)
- [LDAP configuration](ldap.md)
- [Deployment](deployment.md)
