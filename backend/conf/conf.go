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
	Gin Gin `yaml:"gin"`
	Rpc Rpc `yaml:"rpc"`
}

type Gin struct {
	Port string `yaml:"port"`
}

type Rpc struct {
	Consul_address string `yaml:"consul_address"`
}

func GetConf() *Config {
	once.Do(func() {
		conf = new(Config)
		config.Init(conf)
	})
	return conf
}
