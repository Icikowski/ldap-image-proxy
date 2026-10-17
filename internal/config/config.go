package config

import (
	"errors"
	"fmt"
	"net"
	"slices"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/multierr"

	"git.sr.ht/~icikowski/ldap-image-proxy/internal/model"
)

const implementationCustom = "custom"

var allowedImplementations = append(model.KnownLDAPImplementations(), implementationCustom)

// Config holds the configuration for the application.
type Config struct {
	Server ServerConfig `envPrefix:"SERVER_"`
	LDAP   LDAPConfig   `envPrefix:"LDAP_"`
	Log    LogConfig    `envPrefix:"LOG_"`
}

// Validate checks if the [Config] is valid.
func (c *Config) Validate() error {
	errs := make([]error, 0, 3)

	if err := c.Server.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.LDAP.Validate(); err != nil {
		errs = append(errs, err)
	}

	if err := c.Log.Validate(); err != nil {
		errs = append(errs, err)
	}

	return multierr.Combine(errs...)
}

// ServerConfig holds the configuration for the server.
type ServerConfig struct {
	Address  string        `env:"ADDRESS"       envDefault:""`
	Port     uint16        `env:"PORT,notEmpty" envDefault:"3891"`
	ImageTTL time.Duration `env:"IMAGE_TTL"     envDefault:"12h"`
}

// ServerAddr returns the server address in the format "address:port".
func (c *ServerConfig) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Address, c.Port)
}

// Validate checks if the [ServerConfig] is valid.
func (c *ServerConfig) Validate() error {
	errs := make([]error, 0, 3)

	if c.Address != "" {
		if ip := net.ParseIP(c.Address); ip == nil {
			errs = append(errs, errors.New("invalid server address"))
		}
	}

	if c.Port == 0 {
		errs = append(errs, errors.New("server port must be between 1 and 65535"))
	}

	if c.ImageTTL < 0 {
		errs = append(errs, errors.New("image TTL must be a non-negative duration"))
	}

	return multierr.Combine(errs...)
}

// LDAPConfig holds the configuration for the LDAP connection.
type LDAPConfig struct {
	URL               string `env:"URL,notEmpty"`
	BaseDN            string `env:"BASE_DN,notEmpty"`
	BindUser          string `env:"BIND_USER,notEmpty"`
	BindPassword      string `env:"BIND_PASSWORD,notEmpty,unset"`
	Implementation    string `env:"IMPLEMENTATION,notEmpty"      envDefault:"custom"`
	AvatarAttribute   string `env:"ATTRIBUTE_AVATAR,notEmpty"    envDefault:"jpegphoto"`
	UsernameAttribute string `env:"ATTRIBUTE_USERNAME"`
	MailAttribute     string `env:"ATTRIBUTE_EMAIL"`
	SearchFilter      string `env:"SEARCH_FILTER"`
}

// QueryBase returns the [model.LDAPQueryBase] for the given [LDAPConfig].
func (c *LDAPConfig) QueryBase() *model.LDAPQueryBase {
	if c.Implementation == implementationCustom {
		return model.NewLDAPQueryBase(c.UsernameAttribute, c.MailAttribute, c.SearchFilter)
	}

	return model.LDAPImplementation(c.Implementation).QueryBase()
}

// Validate checks if the [LDAPConfig] is valid.
func (c *LDAPConfig) Validate() error {
	errs := make([]error, 0, 8)

	if c.URL == "" {
		errs = append(errs, errors.New("LDAP URL is required"))
	}

	if c.BaseDN == "" {
		errs = append(errs, errors.New("LDAP Base DN is required"))
	}

	if c.BindUser == "" {
		errs = append(errs, errors.New("LDAP bind user is required"))
	}

	if c.BindPassword == "" {
		errs = append(errs, errors.New("LDAP bind password is required"))
	}

	if c.Implementation == "" {
		errs = append(errs, errors.New("LDAP implementation is required"))
	} else if !slices.Contains(allowedImplementations, c.Implementation) {
		errs = append(errs, fmt.Errorf("LDAP implementation must be one of: %v", allowedImplementations))
	}

	if c.Implementation == implementationCustom {
		if c.UsernameAttribute == "" {
			errs = append(errs, errors.New("LDAP username attribute is required for custom implementation"))
		}

		if c.MailAttribute == "" {
			errs = append(errs, errors.New("LDAP mail attribute is required for custom implementation"))
		}

		if c.SearchFilter == "" {
			errs = append(errs, errors.New("LDAP search filter is required for custom implementation"))
		}
	}

	return multierr.Combine(errs...)
}

// LogConfig holds the configuration for logging.
type LogConfig struct {
	Level  zerolog.Level `env:"LEVEL"  envDefault:"info"`
	Pretty bool          `env:"PRETTY" envDefault:"false"`
}

// Validate checks if the [LogConfig] is valid.
func (c LogConfig) Validate() error {
	if c.Level < zerolog.TraceLevel || c.Level > zerolog.PanicLevel {
		return errors.New("invalid log level")
	}
	return nil
}
