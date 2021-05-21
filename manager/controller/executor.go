package controller

import (
	"context"
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Executor interface {
	Run(context.Context, string, ...string) ([]byte, error)
}

type executor struct {
	Verbosity logrus.Level
}

func newExecutor(logLevel logrus.Level) executor {
	e := executor{
		Verbosity: logLevel,
	}

	return e
}

func (e executor) Run(ctx context.Context, command string, arg ...string) ([]byte, error) {
	pipe := exec.CommandContext(ctx, command, arg...)

	errPipe, err := pipe.StderrPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdErr")
	}

	outputPipe, err := pipe.StdoutPipe()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdOut")
	}

	if err := pipe.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start command")
	}

	stdOut, err := ioutil.ReadAll(outputPipe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdOut")
	}

	stdErr, err := ioutil.ReadAll(errPipe)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read command stdErr")
	}

	if len(stdErr) > 0 {
		return nil, errors.New(string(stdErr))
	}

	return stdOut, nil
}
