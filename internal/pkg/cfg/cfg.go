package cfg

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type DBCfg struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	User     string `env:"USER"`
	Password string `env:"PASSWORD"`
	DBName   string `env:"DATABASE_NAME"`
}

type Cfg struct {
	DB DBCfg
}

func init() {
	err := godotenv.Load(os.Getenv("ENV_FILE_PATH"))
	if err != nil {
		log.Panic("can't load .env file")
	}
}

func Read() (*Cfg, error) {
	cfg := Cfg{}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
