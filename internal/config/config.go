package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type RedisConfig struct {
	RedisHost string `yaml:"host"`
	RedisPort string `yaml:"port"`
	Password  string `yaml:"password"`
	DBName    int    `yaml:"DBName"`
}

type FloodControlCfg struct {
	CheckTime  int `yaml:"check_time"`
	CheckCount int `yaml:"check_count"`
}

func GetRedisConfig(path string) *RedisConfig {
	cfg := &RedisConfig{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		log.Fatalf("Can't read Redis cfg: %s\n", err.Error())
		return nil
	}
	return cfg
}

func GetFloodControlConfig(path string) *FloodControlCfg {
	cfg := &FloodControlCfg{}
	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		log.Fatalf("Can't read FloodControl cfg: %s\n", err.Error())
		return nil
	}
	return cfg
}
