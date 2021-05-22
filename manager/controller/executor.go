package controller

import (
	"context"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Executor interface {
	Run(context.Context, string, ...string) ([]byte, error)
}

type executor struct {
	log *logrus.Logger
}

func newExecutor(logger *logrus.Logger) executor {
	e := executor{
		log: logger,
	}

	return e
}

func (e executor) Run(ctx context.Context, command string, arg ...string) ([]byte, error) {
	t1 := time.Now()
	pipe := exec.CommandContext(ctx, command, arg...)

	logger := e.log.WithFields(logrus.Fields{
		"command": command,
	})

	logger.Debug("created pipe")
	errPipe, err := pipe.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdErr")
	}
	logger.Debug("listening stdErr pipe")

	outputPipe, err := pipe.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdOut")
	}
	logger.Debug("listening stdOut pipe")

	logger.Debug("starting command")
	if err := pipe.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	logger.Debug("reading stdOut buffer")
	stdOut, err := ioutil.ReadAll(outputPipe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdOut")
	}

	logger.Debug("reading stdErr buffer")
	stdErr, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdErr")
	}
	logger.Infof("command executed in %s", time.Since(t1))

	if len(stdErr) > 0 {
		return nil, errors.New(string(stdErr))
	}
	return stdOut, nil
}
