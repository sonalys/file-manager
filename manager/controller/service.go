package controller

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

type Service struct {
	rules    []model.Rule
	scripts  map[string]model.ScriptConfiguration
	Mounts   map[string]string
	executor Executor
	Logger   *logrus.Logger
	ctx      context.Context
}

func NewService(ctx context.Context, c model.Config) *Service {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.SetLevel(level)

	s := &Service{
		ctx:      ctx,
		scripts:  c.Scripts,
		rules:    c.Rules,
		Mounts:   c.Mounts,
		Logger:   logger,
		executor: newExecutor(logger),
	}

	if err := s.pullImages(); err != nil {
		panic(err)
	}

	return s
}
