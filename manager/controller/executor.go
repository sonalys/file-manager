package controller

import (
	"context"
	"io/ioutil"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Executor interface {
	Run(context.Context, string, ...string) (<-chan []byte, <-chan error)
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

func (e executor) Run(ctx context.Context, command string, arg ...string) (<-chan []byte, <-chan error) {
	var respChan chan []byte
	var errChan chan error
	cmd := exec.CommandContext(ctx, command, arg...)
	go func() {
		if err := cmd.Run(); err != nil {
			errChan <- err
			return
		}
		stdErr, err := cmd.StderrPipe()
		if err != nil {
			errChan <- err
			return
		}

		errOutput, err := ioutil.ReadAll(stdErr)
		switch true {
		case err != nil:
			errChan <- err
			return
		case len(errOutput) > 0:
			errChan <- errors.New(string(errOutput))
			return
		}

		output, err := cmd.Output()
		if err != nil {
			errChan <- err
			return
		}
		respChan <- output
	}()

	return respChan, errChan
}
