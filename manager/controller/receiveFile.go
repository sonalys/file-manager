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
func (s Service) ResolvePath(d *model.UploadData) (string, error) {
	if strings.Contains(d.Destination, "../") {
		return "", errors.New("invalid path: cannot include ../")
	}

	if strings.Contains(d.Filename, "/") {
		return "", errors.New("invalid filename: cannot include /")
	}

	match, _ := regexp.MatchString(".*:.*", d.Destination)
	if !match {
		return "", errors.New("invalid destination: must have format mount:path")
	}

	splitPath := strings.Split(d.Destination, ":")

	resolve, found := s.Mounts[splitPath[0]]
	if !found {
		return "", errors.New("invalid mount")
	}
	return fmt.Sprintf("%s/%s/%s", resolve, splitPath[1], d.GetFullName()), nil
}

func (s Service) ReceiveFile(reader io.Reader, filename, destination string) (*model.UploadData, error) {
	dotIndex := strings.LastIndex(filename, ".")
	uploadData := &model.UploadData{
		Filename:    filename[:dotIndex],
		Extension:   filename[dotIndex+1:],
		Destination: destination,
		Metadata:    map[string]model.ScriptOutput{},
	}

	path, err := s.ResolvePath(uploadData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse destination during file upload")
	}

	dst, err := os.Create(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file in destination")
	}

	if _, err := io.Copy(dst, reader); err != nil {
		return nil, errors.Wrap(err, "failed to write buffer in file destination")
	}
	dst.Close()

	for _, rule := range s.rules {
		match, err := rule.Match.Validate(filename, destination)
		if err != nil {
			return nil, errors.Wrap(err, "failed to check rule conditions")
		}
		if !match {
			continue
		}

		for _, scriptName := range rule.Pipeline {
			newUploadData := s.Run(scriptName, *uploadData)
			if newUploadData != nil {
				uploadData = newUploadData
			}
		}

		encodedOutput, err := json.Marshal(uploadData)
		if err != nil {
			return nil, errors.Wrap(err, "failed to encode metadata")
		}

		resolvedPath, _ := s.ResolvePath(uploadData)
		err = os.WriteFile(fmt.Sprintf("%s.metadata", resolvedPath), encodedOutput, os.ModeDevice)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create metadata in destination")
		}
	}

	return uploadData, nil
}
