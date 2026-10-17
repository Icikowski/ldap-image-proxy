package ldap

import (
	"errors"
	"fmt"
	"io"

	goldap "github.com/go-ldap/ldap/v3"

	"git.sr.ht/~icikowski/ldap-image-proxy/internal/model"
)

// Client is an interface that defines the methods for interacting with an LDAP server.
type Client interface {
	// UserExists checks if a user with the given username exists in the LDAP directory.
	UserExists(username string) (bool, error)

	// UserPhoto retrieves the photo of a user with the given username from the LDAP directory.
	UserPhoto(username string) ([]byte, error)

	io.Closer
}

type client struct {
	conn            *goldap.Conn
	baseDN          string
	attributeAvatar string
	queryBase       *model.LDAPQueryBase
}

// NewClient creates a new LDAP [Client] with the given configuration.
func NewClient(
	url, username, password, baseDN string,
	avatarAttribute string,
	queryBase *model.LDAPQueryBase,
) (Client, error) {
	conn, err := goldap.DialURL(url)
	if err != nil {
		return nil, err
	}

	if queryBase == nil {
		return nil, errors.New("valid query base must be provided")
	}

	if err := conn.Bind(username, password); err != nil {
		//nolint:errcheck
		defer conn.Close()
		return nil, err
	}

	return &client{
		conn:            conn,
		baseDN:          baseDN,
		attributeAvatar: avatarAttribute,
		queryBase:       queryBase,
	}, nil
}

// UserExists implements [Client].
func (c *client) UserExists(username string) (bool, error) {
	result, err := c.lookupSingleUser(username, []string{c.queryBase.UsernameAttribute(), c.queryBase.MailAttribute()})
	if err != nil {
		return false, err
	}

	return result != nil, nil
}

// UserPhoto implements [Client].
func (c *client) UserPhoto(username string) ([]byte, error) {
	result, err := c.lookupSingleUser(username, []string{c.attributeAvatar})
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	photo := result.GetRawAttributeValue(c.attributeAvatar)
	if len(photo) == 0 {
		return nil, nil
	}

	return photo, nil
}

// Close implements [Client].
func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) lookupSingleUser(username string, attrs []string) (*goldap.Entry, error) {
	searchRequest := goldap.NewSearchRequest(
		c.baseDN,
		goldap.ScopeWholeSubtree, goldap.NeverDerefAliases, 0, 0, false,
		c.queryBase.BuildQuery(username),
		attrs,
		nil,
	)

	searchResult, err := c.conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(searchResult.Entries) == 0 {
		//nolint:nilnil
		return nil, nil
	}

	if len(searchResult.Entries) > 1 {
		return nil, fmt.Errorf("multiple entries found for user %s", username)
	}

	return searchResult.Entries[0], nil
}
