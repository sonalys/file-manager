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

func (s Service) Run(name string) <-chan model.ScriptOutput {
	script, found := s.scripts[name]
	if !found {
		logrus.Error(errors.New(fmt.Sprintf("script not found %s", name)))
		return nil
	}

	var resp chan model.ScriptOutput

	timeout, err := time.ParseDuration(script.Timeout)
	if err != nil {
		logrus.Error(errors.Wrap(err, "failed to parse script timeout"))
		return resp
	}

	timeoutCtx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	output, err := s.executor.Run(timeoutCtx, script.GetCommand())
	if err != nil {
		logrus.Error(errors.Wrap(err, fmt.Sprintf("failed to execute script command %s", name)))
		return resp
	}

	var scriptOutput model.ScriptOutput
	if err := json.Unmarshal(output, &scriptOutput); err != nil {
		logrus.Error(errors.Wrap(err, "failed to deserialize script output"))
		return resp
	}

	return resp
}
