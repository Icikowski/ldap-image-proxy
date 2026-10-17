# LDAP configuration

`ldap-image-proxy` supports several LDAP schemas. The built-in presets choose the username and mail attributes and the search filter for you.

Mentioned LDAP implementations were based on [Authelia](https://www.authelia.com/) search patterns.

## Implementations

Supported values for `LIP_LDAP_IMPLEMENTATION`:

- `activedirectory`
- `freeipa`
- `glauth`
- `lldap`
- `rfc2307bis`
- `custom`

The built-in presets were based on the LDAP search patterns used by Authelia.

## Presets

### `activedirectory`

This implementation represents configuration for Microsoft Active Directory.

| Property                  | Value                                                          |
| ------------------------- | -------------------------------------------------------------- |
| username attribute        | `sAMAccountName`                                               |
| mail attribute            | `mail`                                                         |
| search filter explanation | matches enabled, non-expired user accounts with a password set |

**Search filter:**

```text
(&(|({username_attribute}={input})({mail_attribute}={input}))(sAMAccountType=805306368)(!(userAccountControl:1.2.840.113556.1.4.803:=2))(!(pwdLastSet=0))(|(!(accountExpires=*))(accountExpires=0)(accountExpires>={datetime-ms})))
```

### `freeipa`

This implementation represents configuration for FreeIPA.

| Property                  | Value                                                 |
| ------------------------- | ----------------------------------------------------- |
| username attribute        | `uid`                                                 |
| mail attribute            | `mail`                                                |
| search filter explanation | matches person entries that are not locked or expired |

**Search filter:**

```text
(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=person)(!(nsAccountLock=TRUE))(krbPasswordExpiration>={datetime})(|(!(krbPrincipalExpiration=*))(krbPrincipalExpiration>={datetime})))
```

### `glauth`

This implementation represents configuration for GLAuth.

| Property                  | Value                         |
| ------------------------- | ----------------------------- |
| username attribute        | `cn`                          |
| mail attribute            | `mail`                        |
| search filter explanation | matches active POSIX accounts |

**Search filter:**

```text
(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=posixAccount)(!(accountStatus=inactive)))
```

### `lldap`

This implementation represents configuration for LLDAP.

| Property                  | Value                  |
| ------------------------- | ---------------------- |
| username attribute        | `uid`                  |
| mail attribute            | `mail`                 |
| search filter explanation | matches person entries |

**Search filter:**

```text
(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=person))
```

### `rfc2307bis`

This implementation represents configuration for the RFC2307bis LDAP schema.

| Property                  | Value                                                                               |
| ------------------------- | ----------------------------------------------------------------------------------- |
| username attribute        | `uid`                                                                               |
| mail attribute            | `mail`                                                                              |
| search filter explanation | matches person and organizational person entries that are not password-reset marked |

**Search filter:**

```text
(&(|({username_attribute}={input})({mail_attribute}={input}))(|(objectClass=inetOrgPerson)(objectClass=organizationalPerson))(!(pwdReset=TRUE)))
```

## Custom LDAP schema

Use `custom` when none of the built-in presets fit your LDAP schema.

Set these variables yourself:

- `LIP_LDAP_ATTRIBUTE_USERNAME` - the attribute that stores the login name
- `LIP_LDAP_ATTRIBUTE_EMAIL` - the attribute that stores the email address
- `LIP_LDAP_SEARCH_FILTER` - the LDAP filter used to find one user entry

Example:

```dotenv
LIP_LDAP_IMPLEMENTATION="custom"
LIP_LDAP_ATTRIBUTE_USERNAME="uid"
LIP_LDAP_ATTRIBUTE_EMAIL="mail"
LIP_LDAP_SEARCH_FILTER="(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=person))"
```

That example matches a person entry by either username or email.

## Custom search filters

Use these placeholders in `LIP_LDAP_SEARCH_FILTER`:

- `{username_attribute}` - resolved username field name
- `{mail_attribute}` - resolved mail field name
- `{input}` - requested username from the URL
- `{datetime}` - current UTC generalized time
- `{datetime-ms}` - current Microsoft epoch timestamp
- `{datetime-unix}` - current Unix epoch timestamp

Example:

```text
(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=person))
```

## See also

- [Configuration overview](configuration.md)
- [Environment variables](environment-variables.md)
- [Deployment](deployment.md)
