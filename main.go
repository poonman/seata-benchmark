package main

import (
	"github.com/dotamixer/doom/pkg/lion"
	"github.com/dotamixer/doom/pkg/lion/source/file"
	"github.com/poonman/seata-benchmark/config"
	"github.com/poonman/seata-benchmark/handler"
	"github.com/sirupsen/logrus"
)

func main() {

	err := lion.Load(file.NewSource(file.WithPath("config.yaml")))
	if err != nil {
		logrus.Fatal(err)
	}

	conf := config.NewConfig()

	h := handler.NewHandler(conf)

	h.Run()
}
