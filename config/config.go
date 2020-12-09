package config

import (
	"github.com/dotamixer/doom/pkg/lion"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Benchmark Benchmark `json:"benchmark"`
	Url       string    `json:"url"`
}

type Benchmark struct {
	Concurrency      int `json:"concurrency"`
	RequestNumPerCon int `json:"requestNumPerCon"`
}

func NewConfig() *Config {
	logrus.Infof("new config...")
	conf := Config{}

	err := lion.Get("logic").Scan(&conf)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Infof("logic config:%+v", conf)

	return &conf
}
