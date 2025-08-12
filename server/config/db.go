package config

import (
	"errors"
	"fmt"
	"net/url"
)

var (
	ErrInvalidDialect = errors.New("invalid SQL Dialect")
)

const (
	SQLITE   = "sqlite3"
	POSTGRES = "postgresql"
	MYSQL    = "mysql"
)

// DatabaseConfig - Server config
type DatabaseConfig struct {
	Dialect  string `json:"dialect"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`

	Params map[string]string `json:"params"`

	MaxIdleConns int `json:"max_idle_conns"`
	MaxOpenConns int `json:"max_open_conns"`

	LogLevel string `json:"log_level"`
}

func encodeParams(rawParams map[string]string) string {
	params := url.Values{}
	for key, value := range rawParams {
		params.Add(key, value)
	}
	return params.Encode()
}
func (c *DatabaseConfig) DSN() (string, error) {
	switch c.Dialect {
	case MYSQL:
		user := url.QueryEscape(c.Username)
		password := url.QueryEscape(c.Password)
		db := url.QueryEscape(c.Database)
		host := fmt.Sprintf("%s:%d", url.QueryEscape(c.Host), c.Port)
		params := encodeParams(c.Params)
		return fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", user, password, host, db, params), nil
	case POSTGRES:
		user := url.QueryEscape(c.Username)
		password := url.QueryEscape(c.Password)
		db := url.QueryEscape(c.Database)
		host := url.QueryEscape(c.Host)
		port := c.Port
		params := encodeParams(c.Params)
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s %s", host, port, user, password, db, params), nil
	default:
		return "", ErrInvalidDialect
	}
}
