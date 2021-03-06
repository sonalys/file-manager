package controller

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

func (s Service) Run(scriptName string, upload model.UploadData) *model.UploadData {
	logger := s.Logger.WithFields(logrus.Fields{
		"scriptName": scriptName,
	})
	logger.Debug("searching for script name on service")
	script, found := s.scripts[scriptName]
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

	deadline, _ := timeoutCtx.Deadline()

	logger.Info("started script")
	binary, parameters := script.GetCommand(s.ResolveMount(upload.Destination), scriptName, upload)
	output, err := s.executor.Run(timeoutCtx, binary, parameters...)
	if err != nil || time.Now().After(deadline) {
		logger.Error(errors.Wrap(err, "failed to execute script command"))
		return nil
	}

	lastIndex := strings.LastIndex(string(output), "\n{")

	logger.Debug("deserializing script response")
	var scriptOutput model.UploadData
	if err := json.Unmarshal(output[lastIndex:], &scriptOutput); err != nil {
		logger.Error(errors.Wrap(err, "failed to deserialize script output"), string(output))
		return nil
	}

	logger.Debug("done")
	return &scriptOutput
}
