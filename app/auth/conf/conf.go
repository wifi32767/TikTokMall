package conf

import (
	"sync"

	"github.com/wifi32767/TikTokMall/common/config"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Kitex          Kitex  `yaml:"kitex"`
	Redis          Redis  `yaml:"redis"`
	Consul_address string `yaml:"consul_address"`
}

type Kitex struct {
	Service   string `yaml:"service"`
	Address   string `yaml:"address"`
	Log_level string `yaml:"log_level"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

func GetConf() *Config {
	once.Do(func() {
		conf = new(Config)
		config.Init(conf)
	})
	return conf
}
