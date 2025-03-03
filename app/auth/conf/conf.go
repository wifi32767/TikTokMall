package conf

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/wifi32767/TikTokMall/common/config"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	Kitex Kitex `yaml:"kitex"`
	Redis Redis `yaml:"redis"`
	Log   Log   `yaml:"log"`
}

type Kitex struct {
	Service        string `yaml:"service"`
	Address        string `yaml:"address"`
	Log_level      string `yaml:"log_level"`
	Consul_address string `yaml:"consul_address"`
}

type Redis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type Log struct {
	Level           string `yaml:"level"`
	RabbitmqAddress string `yaml:"rabbitmq_address"`
}

func GetConf() *Config {
	once.Do(func() {
		conf = new(Config)
		config.Init(conf)
	})
	return conf
}

func LogLevel() klog.Level {
	level := GetConf().Log.Level
	switch level {
	case "trace":
		return klog.LevelTrace
	case "debug":
		return klog.LevelDebug
	case "info":
		return klog.LevelInfo
	case "notice":
		return klog.LevelNotice
	case "warn":
		return klog.LevelWarn
	case "error":
		return klog.LevelError
	case "fatal":
		return klog.LevelFatal
	default:
		return klog.LevelInfo
	}
}
