package initializers

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost           string `mapstructure:"POSTGRES_HOST"`
	DBUserName       string `mapstructure:"POSTGRES_USER"`
	DBUserPassword   string `mapstructure:"POSTGRES_PASSWORD"`
	DBName           string `mapstructure:"POSTGRES_DB"`
	DBPort           string `mapstructure:"POSTGRES_PORT"`
	ServerPort       string `mapstructure:"PORT"`
	ResendMailAPIKey string `mapstructure:"RESEND_MAIL_API_KEY"`
	EmailFrom        string `mapstructure:"EMAIL_FROM"`
	SessionSecret    string `mapstructure:"SESSION_SECRET"`

	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
	BcryptCost             int           `mapstructure:"BCRYPT_COST"`
}

func MustLoadConfig(path string, stage string) (config *Config) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(stage)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Could not load environment variables", err)
		return nil
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Could not load environment variables", err)
		return nil
	}

	return
}
