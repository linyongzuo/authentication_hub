package main

import (
	"github.com/authentication_hub/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("服务启动")
	service.Start()
}
