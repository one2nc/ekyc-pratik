package config

import (
	"fmt"
	"os"
)

type RedisConfig struct {
	Address  string
	Port     string
	Endpoint string
}

func NewRedisConfig() RedisConfig {
	address := os.Getenv("REDIS_ADDRESS")
	port := os.Getenv("REDIS_PORT")
	return RedisConfig{
		Address:  address,
		Port:     port,
		Endpoint: fmt.Sprintf("%s:%s", address, port),
	}
}
