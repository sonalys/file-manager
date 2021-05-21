package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sonalys/file-manager/manager/model"
)

func (s Service) ReceiveFile(reader io.Reader, filename, destination string) error {
	path := fmt.Sprintf("%s/%s", destination, filename)
	dst, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file in destination")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return errors.Wrap(err, "failed to write buffer in file destination")
	}

	metadata := []*model.ScriptOutput{}

	for _, rule := range s.rules {
		match, err := rule.Match.Validate(filename, destination)
		if err != nil {
			return errors.Wrap(err, "failed to check rule conditions")
		}
		if !match {
			continue
		}

		for _, scriptName := range rule.Pipeline {
			output := s.Run(scriptName)
			if output == nil {
				continue
			}
			metadata = append(metadata, output)
		}

		encodedOutput, err := json.Marshal(metadata)
		if err != nil {
			return errors.Wrap(err, "failed to encode metadata")
		}

		err = os.WriteFile(fmt.Sprintf("%s.metadata", path), encodedOutput, os.ModeDevice)
		if err != nil {
			return errors.Wrap(err, "failed to create metadata in destination")
		}
	}

	return nil
}
