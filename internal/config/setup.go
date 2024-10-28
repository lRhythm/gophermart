package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
)

func New() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return cfg.withFlags().withDefault(), nil
}

func (c *Config) withFlags() *Config {
	sa, dd, aa := new(string), new(string), new(string)
	var parse bool
	if flag.Lookup(a) == nil {
		sa = flag.String(a, "", "Server net address (host:port / host / :port)")
		parse = true
	}
	if flag.Lookup(d) == nil {
		dd = flag.String(d, "", "Database DSN")
		parse = true
	}
	if flag.Lookup(r) == nil {
		aa = flag.String(r, "", "Accrual net address (host:port / host / :port)")
		parse = true
	}
	if parse {
		flag.Parse()
	}
	if *sa != "" && c.Server.Address == "" {
		c.Server.Address = *sa
	}
	if *dd != "" && c.Database.DSN == "" {
		c.Database.DSN = *dd
	}
	if *aa != "" && c.Accrual.Address == "" {
		c.Accrual.Address = *aa
	}
	return c
}

func (c *Config) withDefault() *Config {
	if c.Server.Address == "" {
		c.Server.Address = ":1234"
	}
	if c.Database.DSN == "" {
		c.Database.DSN = "host=127.0.0.1 port=5432 user=user password=secret dbname=postgres sslmode=disable"
	}
	if c.Accrual.Address == "" {
		c.Accrual.Address = "http://localhost:8080"
	}
	return c
}
