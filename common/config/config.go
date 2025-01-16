package config

import (
	"fmt"
	"io/ioutil"
	"log/slog"
	"path/filepath"

	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v2"
)

func Init(conf any) {
	prefix := "conf"
	confFileRelPath := filepath.Join(prefix, "conf.yaml")
	content, err := ioutil.ReadFile(confFileRelPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, conf)
	if err != nil {
		slog.Error("parse yaml error - ", "error", err)
		panic(err)
	}
	if err := validator.Validate(conf); err != nil {
		slog.Error("validate config error - ", "error", err)
		panic(err)
	}
	fmt.Printf("%+v\n", conf)
}
