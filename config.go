package couchdb

// Config mongodb configuration parameters
type Config struct {
	URL      string
	DB       string
	Username string
	Password string
}

// NewConfig create mongodb configuration
func NewConfig(url, db string, username string, password string) *Config {
	return &Config{
		URL:      url,
		DB:       db,
		Username: username,
		Password: password,
	}
}
