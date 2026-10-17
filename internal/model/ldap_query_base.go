package model

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

const (
	placeholderUsernameAttribute      = "{username_attribute}"
	placeholderMailAttribute          = "{mail_attribute}"
	placeholderInput                  = "{input}"
	placeholderDateTime               = "{datetime}"
	placeholderDateTimeEpochMicrosoft = "{datetime-ms}"
	placeholderDateTimeEpochUnix      = "{datetime-unix}"

	formatDateTime         = "20060102150405.0Z"
	offsetMicrosoftNTEpoch = 116444736000000000
)

// LDAPQueryBase represents the base structure for an LDAP query, including the username and mail attributes and
// the search filter.
type LDAPQueryBase struct {
	attributeUsername string
	attributeMail     string
	searchFilter      string
}

// NewLDAPQueryBase creates a new instance of [LDAPQueryBase] with the provided attributes and search filter.
func NewLDAPQueryBase(attributeUsername, attributeMail, searchFilter string) *LDAPQueryBase {
	lqb := &LDAPQueryBase{
		attributeUsername: attributeUsername,
		attributeMail:     attributeMail,
		searchFilter:      searchFilter,
	}

	lqb.prepareSearchFilter()

	return lqb
}

// BuildQuery constructs the LDAP search filter by replacing placeholders with actual values, including the provided
// username and current date & time.
func (lqb *LDAPQueryBase) BuildQuery(username string) string {
	var (
		out = lqb.searchFilter
		now = time.Now()

		dateTime               = now.UTC().Format(formatDateTime)
		dateTimeEpochMicrosoft = strconv.FormatUint(offsetMicrosoftNTEpoch+uint64(now.UnixNano()/100), 10)
		dateTimeEpochUnix      = strconv.FormatInt(now.Unix(), 10)
	)

	out = strings.ReplaceAll(out, placeholderInput, ldap.EscapeFilter(username))
	out = strings.ReplaceAll(out, placeholderDateTime, dateTime)
	out = strings.ReplaceAll(out, placeholderDateTimeEpochMicrosoft, dateTimeEpochMicrosoft)
	out = strings.ReplaceAll(out, placeholderDateTimeEpochUnix, dateTimeEpochUnix)

	return out
}

// SearchFilter returns the base search filter.
func (lqb *LDAPQueryBase) SearchFilter() string {
	return lqb.searchFilter
}

// UsernameAttribute returns the username attribute used in the LDAP query.
func (lqb *LDAPQueryBase) UsernameAttribute() string {
	return lqb.attributeUsername
}

// MailAttribute returns the mail attribute used in the LDAP query.
func (lqb *LDAPQueryBase) MailAttribute() string {
	return lqb.attributeMail
}

func (lqb *LDAPQueryBase) prepareSearchFilter() {
	lqb.searchFilter = strings.ReplaceAll(lqb.searchFilter, placeholderUsernameAttribute, lqb.attributeUsername)
	lqb.searchFilter = strings.ReplaceAll(lqb.searchFilter, placeholderMailAttribute, lqb.attributeMail)
}
