package controller

import (
	"github.com/pkg/errors"
	"github.com/sonalys/file-manager/manager/model"
)

func (s Service) pullImages() error {
	logger := s.Logger
	logger.Info("pulling script images")
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
	imgPath := script.GetImagePath()
	logger := s.Logger.WithField("image", imgPath)
	logger.Info("checking image in cache")

	if s.checkImage(imgPath) {
		logger.Info("version already exists")
		return nil
	}

	logger.Info("downloading image")
	_, err := s.executor.Run(s.ctx, "docker", "pull", imgPath)
	if err != nil {
		return errors.Wrap(err, "error pulling image")
	}

	logger.Info("image downloaded")
	return nil
}

func (s Service) checkImage(image string) bool {
	_, err := s.executor.Run(s.ctx, "docker", "image", "inspect", image)
	return err == nil
}
