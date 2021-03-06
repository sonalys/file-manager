package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/controller"
	"github.com/sonalys/file-manager/manager/handler"
	"github.com/sonalys/file-manager/manager/model"
	"github.com/sonalys/file-manager/manager/util"
)

func main() {
	ctx := context.Background()
	var config model.Config
	err := util.ReadJSON("config.json", &config)
	if err != nil {
		panic(err)
	}

	s := controller.NewService(ctx, config)

	handler := handler.NewHandler(s)
	go handler.Start()

	gracefulShutdown := make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGINT, syscall.SIGTERM)

	<-gracefulShutdown
	logrus.Info("Service stopped")
}
