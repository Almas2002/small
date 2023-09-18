package mail

import "gopkg.in/gomail.v2"

type Config struct {
	Host  int    `yaml:"host"`
	Email string `yaml:"email"`
	Pass  string `yaml:"pass"`
	Smtp  string `yaml:"smtp"`
}

func New(config *Config) *gomail.Dialer {
	return gomail.NewDialer(config.Smtp, config.Host, config.Email, config.Pass)
}
