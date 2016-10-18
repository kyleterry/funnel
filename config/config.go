package config

import (
	"flag"
	"fmt"
	"os"
)

const DatabaseFilename = "funnel.db"

type Config struct {
	HTTPAddr     string
	Debug        bool
	DatabaseDir  string
	DatabaseFile string
}

func New() *Config {
	c := &Config{}

	flag.StringVar(&c.HTTPAddr, "bind", "localhost:6060", "Address and port to bind the HTTP server too")
	flag.BoolVar(&c.Debug, "debug", false, "Enable debugging")
	flag.StringVar(&c.DatabaseDir, "database-dir", fmt.Sprintf("%s/.config/funnel", os.Getenv("HOME")), "Path to the sqlite database")

	if _, err := os.Stat(c.DatabaseDir); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(c.DatabaseDir, os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	c.DatabaseFile = fmt.Sprintf("%s/%s", c.DatabaseDir, DatabaseFilename)

	flag.Parse()

	return c
}
