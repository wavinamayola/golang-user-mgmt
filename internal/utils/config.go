package utils

import (
	"fmt"

	"github.com/spf13/viper"
	"github.com/wavinamayola/user-management/internal/config"
)

func LoadConfig(path string) (cfg config.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)
	return
}

func NewDBStringFromDBConfig(cfg config.Config) (string, error) {
	switch "" {
	case cfg.Database.User:
		return "", fmt.Errorf("database.user is empty")
	case cfg.Database.Password:
		return "", fmt.Errorf("database.password is empty")
	case cfg.Database.DBName:
		return "", fmt.Errorf("database.dbname is empty")
	case cfg.Database.Host:
		return "", fmt.Errorf("database.host is empty")
	case cfg.Database.Port:
		return "", fmt.Errorf("database.port is empty")
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	), nil
}
