package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Config struct {
	HTTPAddr         string
	Debug            bool
	DatabaseEngine   string
	ConnectionString string
	Multiuser        bool
}

func New() (*Config, error) {
	c := &Config{}

	var databaseURL string

	flag.StringVar(&c.HTTPAddr, "bind", "localhost:6060", "Address and port to bind the HTTP server too")
	flag.BoolVar(&c.Debug, "debug", false, "Enable debugging")
	flag.StringVar(&databaseURL, "database-url", fmt.Sprintf("sqlite3://%s/.config/funnel/funnel.db", os.Getenv("HOME")), "URL style connection string for the database (sqlite3 and postgres are currently supported)")
	flag.BoolVar(&c.Multiuser, "multiuser", false, "Enable multiuser")

	flag.Parse()

	dbURL, err := url.Parse(databaseURL)

	if err != nil {
		return nil, errors.Wrap(err, "error while parsing database URL")
	}

	switch dbURL.Scheme {
	case "sqlite3":
		if err := sqlite3Config(c, dbURL); err != nil {
			return nil, errors.Wrap(err, "could not parse sqlite3 URL")
		}
	case "postgres":
		if err := postgresConfig(c, dbURL); err != nil {
			return nil, errors.Wrap(err, "could not parse postgres URL")
		}
	}

	return c, nil
}

func sqlite3Config(c *Config, u *url.URL) error {
	if _, err := os.Stat(u.Path); err != nil {
		if os.IsNotExist(err) {
			parts := strings.Split(u.Path, "/")
			err := os.MkdirAll(strings.Join(parts[:len(parts)-1], "/"), os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "directory creation failed")
			}
		} else {
			return errors.Wrap(err, "error while trying to stat database file")
		}
	}

	c.ConnectionString = u.Path
	c.DatabaseEngine = u.Scheme

	return nil
}

func postgresConfig(c *Config, u *url.URL) error {
	c.ConnectionString = u.String()
	c.DatabaseEngine = u.Scheme
	return nil
}
