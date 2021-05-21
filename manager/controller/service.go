package controller

import (
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/models"
)

type Executor interface {
	Run(string, ...string) (<-chan string, <-chan error)
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

func (e executor) Run(command string, arg ...string) (<-chan string, <-chan error) {
	var respChan chan string
	var errChan chan error
	cmd := exec.Command(command, arg...)
	go func() {
		if err := cmd.Run(); err != nil {
			errChan <- err
			return
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			errChan <- err
			return
		}
		respChan <- string(output)
	}()

	return respChan, errChan
}

type Service struct {
	rules    []models.Rule
	scripts  []models.ScriptConfiguration
	executor Executor
}

func NewService(c models.Config) *Service {
	return &Service{
		executor: newExecutor(c.LogLevel),
	}
}
