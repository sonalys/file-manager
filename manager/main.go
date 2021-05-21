package main

import (
	"context"

	"github.com/sonalys/file-manager/manager/controller"
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

	_ = controller.NewService(ctx, config)
}
