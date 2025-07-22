package conf

import (
	"sync"

	"github.com/wifi32767/TikTokMall/common/config"
)

type Config struct {
	ConsulAddress string
}

var (
	conf *Config
	once sync.Once
)

func GetConf() *Config {
	once.Do(func() {
		conf = new(Config)
		config.Init(conf)
	})
	return conf
}
