package model

import (
	"fmt"
	"os"
	"strings"
)

// ScriptConfiguration represents a script name with parameters, to be executed.
type ScriptConfiguration struct {
	Image      string   `json:"image"`      // Optional: Reference to the docker image to be used.
	Version    string   `json:"version"`    // Optional: Version of the script to be used.
	Binary     string   `json:"binary"`     // Binary to be called on cli.
	Parameters []string `json:"parameters"` // Binary to be called on cli.
	Timeout    string   `json:"timeout"`    // Max execution time for the script to run.
}

func (s ScriptConfiguration) GetCommand(filename string) (string, []string) {
	index := strings.LastIndex(filename, ".")

	file := filename[:index]
	ext := filename[index+1:]

	for i := range s.Parameters {
		s.Parameters[i] = strings.ReplaceAll(s.Parameters[i], "%filename", file)
		s.Parameters[i] = strings.ReplaceAll(s.Parameters[i], "%extension", ext)
	}

	if len(s.Image) > 0 {
		cwd, _ := os.Getwd()
		parameters := []string{"run", "--rm", "-v", fmt.Sprintf("%s/storage:/buffer", cwd), s.GetImagePath()}
		return s.Binary, append(parameters, s.Parameters...)
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
