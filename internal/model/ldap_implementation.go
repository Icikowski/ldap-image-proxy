package model

const (
	attributeSAMAccountName = "sAMAccountName"
	attributeUID            = "uid"
	attributeCN             = "cn"
	attributeMail           = "mail"

	//nolint:lll
	searchFilterActiveDirectory = "(&(|({username_attribute}={input})({mail_attribute}={input}))(sAMAccountType=805306368)(!(userAccountControl:1.2.840.113556.1.4.803:=2))(!(pwdLastSet=0))(|(!(accountExpires=*))(accountExpires=0)(accountExpires>={datetime-ms})))"
	//nolint:lll
	searchFilterFreeIPA = "(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=person)(!(nsAccountLock=TRUE))(krbPasswordExpiration>={datetime})(|(!(krbPrincipalExpiration=*))(krbPrincipalExpiration>={datetime})))"
	//nolint:lll
	searchFilterGLAuth = "(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=posixAccount)(!(accountStatus=inactive)))"

	searchFilterLLDAP = "(&(|({username_attribute}={input})({mail_attribute}={input}))(objectClass=person))"
	//nolint:lll
	searchFilterRFC2307bis = "(&(|({username_attribute}={input})({mail_attribute}={input}))(|(objectClass=inetOrgPerson)(objectClass=organizationalPerson))(!(pwdReset=TRUE)))"
)

// LDAPImplementation represents the type of LDAP implementation.
type LDAPImplementation string

const (
	// LDAPImplementationActiveDirectory represents the Active Directory LDAP implementation.
	LDAPImplementationActiveDirectory LDAPImplementation = "activedirectory"

	// LDAPImplementationFreeIPA represents the FreeIPA LDAP implementation.
	LDAPImplementationFreeIPA LDAPImplementation = "freeipa"

	// LDAPImplementationGLAuth represents the GLAuth LDAP implementation.
	LDAPImplementationGLAuth LDAPImplementation = "glauth"

	// LDAPImplementationLLDAP represents the LLDAP LDAP implementation.
	LDAPImplementationLLDAP LDAPImplementation = "lldap"

	// LDAPImplementationRFC2307bis represents the RFC2307bis LDAP implementation.
	LDAPImplementationRFC2307bis LDAPImplementation = "rfc2307bis"
)

// QueryBase returns the [LDAPQueryBase] for the given LDAP implementation.
func (i LDAPImplementation) QueryBase() *LDAPQueryBase {
	switch i {
	case LDAPImplementationActiveDirectory:
		return NewLDAPQueryBase(attributeSAMAccountName, attributeMail, searchFilterActiveDirectory)
	case LDAPImplementationFreeIPA:
		return NewLDAPQueryBase(attributeUID, attributeMail, searchFilterFreeIPA)
	case LDAPImplementationGLAuth:
		return NewLDAPQueryBase(attributeCN, attributeMail, searchFilterGLAuth)
	case LDAPImplementationLLDAP:
		return NewLDAPQueryBase(attributeUID, attributeMail, searchFilterLLDAP)
	case LDAPImplementationRFC2307bis:
		return NewLDAPQueryBase(attributeUID, attributeMail, searchFilterRFC2307bis)
	default:
		return nil
	}
}

// KnownLDAPImplementations returns a list of known LDAP implementations.
func KnownLDAPImplementations() []string {
	return []string{
		string(LDAPImplementationActiveDirectory),
		string(LDAPImplementationFreeIPA),
		string(LDAPImplementationGLAuth),
		string(LDAPImplementationLLDAP),
		string(LDAPImplementationRFC2307bis),
	}
}
