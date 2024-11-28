package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-required:"true"`
	PostgresURL string     `yaml:"postgres_url" env-required:"true"`
	RedisURL    string     `yaml:"redis_url" env-required:"true"`
	SSO         SSOService `yaml:"sso" env-required:"true"`
	Jwt         JwtConfig  `yaml:"jwt" env-required:"true"`
}

type JwtConfig struct {
	Secret     string        `yaml:"secret" env-required:"true"`
	RefreshTTL time.Duration `yaml:"refresh_ttl" env-required:"true"`
	AccessTTL  time.Duration `yaml:"access_ttl" env-required:"true"`
}

type SSOService struct {
	Port      int           `yaml:"port" env-required:"true"`
	Timeout   time.Duration `yaml:"timeout" env-required:"true"`
	SecretKey string        `yaml:"secret_key" env-required:"true"`
}

func MustLoad() *Config {
	path := fectConfigPath()
	if path == "" {
		panic("config path empty")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config doesnot exists: " + path)
	}
	cfg := new(Config)

	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}

func fectConfigPath() (res string) {
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return
}
