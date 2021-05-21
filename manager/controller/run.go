package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sonalys/file-manager/manager/model"
)

func (s Service) Run(name string) *model.ScriptOutput {
	logger := s.logger.WithField("script", name)
	logger.Debug("searching for script name on service")
	script, found := s.scripts[name]
	if !found {
		logger.Error("script not found")
		return nil
	}
	logger.Debug("script found")

	timeout, err := time.ParseDuration(script.Timeout)
	if err != nil {
		logger.Error(errors.Wrap(err, "failed to parse script timeout"))
		return nil
	}
	logger.Debug("timeout defined as %s", timeout)

	timeoutCtx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()

	logger.Info("started script")
	output, err := s.executor.Run(timeoutCtx, script.GetCommand())
	if err != nil {
		logger.Error(errors.Wrap(err, fmt.Sprintf("failed to execute script command %s", name)))
		return nil
	}

	logger.Debug("deserializing script response")
	var scriptOutput model.ScriptOutput
	if err := json.Unmarshal(output, &scriptOutput); err != nil {
		logger.Error(errors.Wrap(err, "failed to deserialize script output"))
		return nil
	}

	logger.Debug("done")
	return &scriptOutput
}
