package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sonalys/file-manager/manager/model"
)

// ResolvePath is a security validator to avoid reaching system files.
// It also resolves any mounts configured on the service.
func (s Service) ResolvePath(destination, filename string) (string, error) {
	if strings.Contains(destination, "../") {
		return "", errors.New("invalid path: cannot include ../")
	}

	if strings.Contains(filename, "/") {
		return "", errors.New("invalid filename: cannot include /")
	}

	match, _ := regexp.MatchString(".*:.*", destination)
	if !match {
		return "", errors.New("invalid destination: must have format mount:path")
	}

	splitPath := strings.Split(destination, ":")

	resolve, found := s.Mounts[splitPath[0]]
	if !found {
		return "", errors.New("invalid mount")
	}
	return fmt.Sprintf("%s/%s/%s", resolve, splitPath[1], filename), nil
}

func (s Service) ReceiveFile(reader io.Reader, filename, destination string) ([]string, error) {
	path, err := s.ResolvePath(destination, filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse destination during file upload")
	}

	dst, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file in destination")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return nil, errors.Wrap(err, "failed to write buffer in file destination")
	}

	metadata := map[string]*model.ScriptOutput{}

	createdFiles := []string{}

	for _, rule := range s.rules {
		match, err := rule.Match.Validate(filename, destination)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check rule conditions")
		}
		if !match {
			continue
		}

		for _, scriptName := range rule.Pipeline {
			output := s.Run(scriptName, filename)
			if output == nil {
				continue
			}

			if len(output.MovedTo) > 0 {
				newName := output.MovedTo[strings.LastIndex(output.MovedTo, "/")+1:]
				path, _ = s.ResolvePath(destination, newName)
				createdFiles = append(createdFiles, fmt.Sprintf("%s/%s", destination, newName))
			} else {
				createdFiles = append(createdFiles, fmt.Sprintf("%s/%s", destination, filename))
			}

			for i := range output.Children {
				createdFiles = append(createdFiles, fmt.Sprintf("%s/%s", destination, output.Children[i]))
			}

			metadata[scriptName] = output
		}

		encodedOutput, err := json.Marshal(metadata)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode metadata")
		}

		err = os.WriteFile(fmt.Sprintf("%s.metadata", path), encodedOutput, os.ModeDevice)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create metadata in destination")
		}
	}

	return createdFiles, nil
}
