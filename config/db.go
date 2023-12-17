package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// DatabaseConfiguration is a config struct for postgres used by fluxy
type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

// RedisConfiguration is a config struct for redis used by fluxy
type RedisConfiguration struct {
	Host		string
	Password	string
}

// GetDSNConfig constructs and returns a dsn using env config for database
func GetDSNConfig() string {
	dbName := viper.GetString("DB_NAME")
	user := viper.GetString("DB_USER")
	password := viper.GetString("DB_PASSWORD")
	host := viper.GetString("DB_HOST")
	port := viper.GetString("DB_PORT")
	sslMode := viper.GetString("SSL_MODE")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbName, port, sslMode,
	)

	return dsn
}

// GetRedisConfig function returns the redis configuration
func GetRedisConfig() *RedisConfiguration {
	redisHost := viper.GetString("REDIS_HOST")
	redisPort := viper.GetString("REDIS_PORT")
	redisPassword := viper.GetString("REDIS_PASSWORD")

	host := fmt.Sprintf("%s:%s", redisHost, redisPort)
	return &RedisConfiguration{Host: host, Password: redisPassword}
}