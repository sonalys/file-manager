package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sonalys/file-manager/manager/model"
)

func (s Service) pullImages() error {
	logrus.Info("Pulling script images")
	for _, script := range s.scripts {
		if len(script.Image) == 0 {
			continue
		}

		if err := s.pullImage(script); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) pullImage(script model.ScriptConfiguration) error {
	imgPath := fmt.Sprintf("%s:%s", script.Image, script.Version)
	logrus.Infof("Checking %s", imgPath)

	if s.checkImage(imgPath) {
		logrus.Infof("%s already exists", imgPath)
		return nil
	}

	_, err := s.executor.Run(s.ctx, "docker", "pull", imgPath)
	if err != nil {
		errors.Wrap(err, fmt.Sprintf("error pulling image %s:%s", script.Image, script.Version))
	}

	logrus.Infof("%s downloaded", imgPath)
	return nil
}

func (s Service) checkImage(image string) bool {
	_, err := s.executor.Run(s.ctx, "docker", "image", "inspect", image)
	return err == nil
}
