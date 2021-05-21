package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

type Service struct {
	rules    []model.Rule
	scripts  map[string]model.ScriptConfiguration
	executor Executor
	ctx      context.Context
}

func NewService(ctx context.Context, c model.Config) *Service {
	level, err := logrus.ParseLevel(c.LogLevel)
	if err != nil {
		panic(err)
	}

	return &Service{
		ctx:      ctx,
		scripts:  c.Scripts,
		rules:    c.Rules,
		executor: newExecutor(level),
	}
}

func (s Service) Run(name string) <-chan model.ScriptOutput {
	script, found := s.scripts[name]
	if !found {
		logrus.Error(errors.New(fmt.Sprintf("script not found %s", name)))
		return nil
	}

	var resp chan model.ScriptOutput

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	timeout, err := time.ParseDuration(script.Timeout)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to parse script timeout"))
		return resp
	}

	switch true {
	case len(script.Image) == 0:
		outputChan, errChan := s.executor.Run(ctx, script.GetCommand())
		select {
		case output := <-outputChan:
			var scriptOutput model.ScriptOutput
			err := json.Unmarshal(output, &scriptOutput)
			if err != nil {
				logrus.Error(errors.Wrap(err, "failed to deserialize script output"))
				return resp
			}
			resp <- scriptOutput
		case err := <-errChan:
			logrus.Error(errors.Wrap(err, "failed to execute script"))
		case <-time.After(timeout):
			cancel()
			logrus.Errorf("script %s timed out", name)
		}
	}

	return resp
}
