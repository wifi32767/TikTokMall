package common

import (
	"testing"

	"github.com/wifi32767/TikTokMall/common/config"
)

type Config struct {
	A A `yaml:"A"`
	B B `yaml:"B"`
}

type A struct {
	A1 string `yaml:"A1"`
	A2 string `yaml:"A2"`
}

type B struct {
	B1 string `yaml:"B1"`
	B2 string `yaml:"B2"`
}

func TestConfig(t *testing.T) {
	conf := &Config{}
	config.Init(conf)
	t.Log(conf.A.A1 == "aaa")
}
