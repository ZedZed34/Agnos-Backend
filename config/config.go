package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost         string `mapstructure:"DB_HOST"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPassword     string `mapstructure:"DB_PASSWORD"`
	DBName         string `mapstructure:"DB_NAME"`
	DBPort         string `mapstructure:"DB_PORT"`
	HospitalAPIURL string `mapstructure:"HOSPITAL_A_API_URL"`
	JWTSecret      string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// 1. Set Defaults
	viper.SetDefault("DB_HOST", "db")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("HOSPITAL_A_API_URL", "https://hospital-a.api.co.th")
	viper.SetDefault("JWT_SECRET", "change_this_secret_in_prod")

	// 2. EXPLICITLY BIND ENV VARS (Critical Fix)
	// Viper often ignores env vars if it doesn't know the key exists yet.
	_ = viper.BindEnv("DB_USER")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_PORT")

	// If running in Docker, these env vars will overwrite defaults
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found, using environment variables")
	}

	var config Config
	err := viper.Unmarshal(&config)
	return config, err
}
