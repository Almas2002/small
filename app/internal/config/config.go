package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"small/pkg/constans"
	"small/pkg/store/postgres"
	"small/pkg/tools/mail"
	"small/pkg/tracing"
	"strconv"
)

var configPath string

type Http struct {
	Port               string `mapstructure:"port"`
	Development        bool   `mapstructure:"development"`
	BaseUserPath       string `mapstructure:"baseUserPath"`
	BaseProductPath    string `mapstructure:"baseProductPath"`
	DebugErrorResponse bool   `mapstructure:"debugErrorResponse"`
}
type Config struct {
	ServiceName string           `mapstructure:"serviceName"`
	Http        Http             `mapstructure:"http"`
	JwtSecret   string           `mapstructure:"jwtSecret"`
	Postgres    *postgres.Config `mapstructure:"postgres"`
	Jaeger      *tracing.Config  `mapstructure:"jaeger"`
	Email       *mail.Config     `mapstructure:"mail"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constans.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getWd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			parentDir := filepath.Join(getWd, "..")
			err = os.Chdir(parentDir)
			if err != nil {
				fmt.Println("Error:", err)
				return nil, nil
			}
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error:", err)
				return nil, nil
			}
			configPath = fmt.Sprintf("%s/configs/config.yml", cwd)
		}

	}

	cfg := &Config{}
	viper.SetConfigType(constans.YAML)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	// postgres
	pgHost := os.Getenv(constans.PGHOST)
	if pgHost != "" {
		cfg.Postgres.Host = pgHost
	}
	pgPassword := os.Getenv(constans.PGPASSWORD)
	if pgPassword != "" {
		cfg.Postgres.Password = pgPassword
	}
	pgPort := os.Getenv(constans.PGPORT)
	if pgPort != "" {
		cfg.Postgres.Port = pgPort
	}
	pgUser := os.Getenv(constans.PGUSER)
	if pgUser != "" {
		cfg.Postgres.User = pgUser
	}

	jaegerHostPort := os.Getenv(constans.JAEGERHOSTPORT)
	if jaegerHostPort != "" {
		cfg.Jaeger.HostPort = jaegerHostPort
	}

	smsHost := os.Getenv(constans.SMSHOST)
	if smsHost != "" {
		parseInt, err := strconv.ParseInt(smsHost, 10, 64)
		if err != nil {

		} else {
			cfg.Email.Host = int(parseInt)
		}
	}

	smsEmail := os.Getenv(constans.SMSEMAIL)
	if smsEmail != "" {
		cfg.Email.Email = smsEmail
	}

	smsPass := os.Getenv(constans.SMSPASS)
	if smsPass != "" {
		cfg.Email.Pass = smsPass
	}

	smsSMTP := os.Getenv(constans.SMSSMTP)
	if smsSMTP != "" {
		cfg.Email.Smtp = smsSMTP
	}

	return cfg, nil

}
