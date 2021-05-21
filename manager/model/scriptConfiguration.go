package model

import (
	"fmt"
)

// ScriptConfiguration represents a script name with parameters, to be executed.
type ScriptConfiguration struct {
	Image   string `json:"image"`   // Optional: Reference to the docker image to be used.
	Version string `json:"version"` // Optional: Version of the script to be used.
	Command string `json:"command"` // Binary to be called on cli.
	Timeout string `json:"timeout"` // Max execution time for the script to run.
}

func (s ScriptConfiguration) GetCommand() string {
	if len(s.Image) > 0 {
		return fmt.Sprintf("docker run -it --rm -v $(pwd):/buffer %s:%s %s", s.Image, s.Version, s.Command)
	}
	return s.Command
}
