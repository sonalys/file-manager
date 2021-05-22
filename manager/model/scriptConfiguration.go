package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// ScriptConfiguration represents a script name with parameters, to be executed.
type ScriptConfiguration struct {
	Image      string   `json:"image"`      // Optional: Reference to the docker image to be used.
	Name       string   `json:"name"`       // Script name.
	Version    string   `json:"version"`    // Optional: Version of the script to be used.
	Binary     string   `json:"binary"`     // Binary to be called on cli.
	Parameters []string `json:"parameters"` // Binary to be called on cli.
	Timeout    string   `json:"timeout"`    // Max execution time for the script to run.
}

func (s ScriptConfiguration) GetCommand(location, name string, u UploadData) (string, []string) {
	s.Name = name
	for i := range s.Parameters {
		s.Parameters[i] = strings.ReplaceAll(s.Parameters[i], "%filename", u.Filename)
		s.Parameters[i] = strings.ReplaceAll(s.Parameters[i], "%extension", u.Extension)
	}

	if len(s.Image) > 0 {
		cwd, _ := os.Getwd()
		parameters := []string{
			"run",
			"--rm",
			"-v", fmt.Sprintf("%s:/buffer", location),
			"-v", fmt.Sprintf("%s/scripts:/scripts", cwd),
			s.GetImagePath(),
			s.Binary,
		}

		serializedData, _ := json.Marshal(u)
		serializedConfig, _ := json.Marshal(s)
		return "docker", append(
			parameters,
			string(serializedData),
			string(serializedConfig),
		)
	}
	return s.Binary, s.Parameters
}

func (s ScriptConfiguration) GetImagePath() string {
	tag := s.Version
	if len(s.Version) == 0 {
		tag = "latest"
	}

	return fmt.Sprintf("%s:%s", s.Image, tag)
}
