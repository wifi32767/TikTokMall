package conf

import (
	"sync"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/wifi32767/TikTokMall/common/config"
)

type Config struct {
	ConsulAddress string `yaml:"consul_address"`
	MqAddress     string `yaml:"mq_address"`
	Log           Log    `yaml:"log"`
}

type Log struct {
	Level     string `yaml:"level"`
	MqAddress string `yaml:"mq_address"`
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
