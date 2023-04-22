package comm

import (
	"encoding/json"
	"os"
)

type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"api"`
}

type Config struct {
	Redis RedisConfig `json:"redis"`
}

var Conf = &Config{}

func InitConfig() {
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, Conf)
}
