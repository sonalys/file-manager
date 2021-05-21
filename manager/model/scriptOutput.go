package model

import "time"

// ScriptOutput represents all the possible script output structures.
// Manager will interprete it and aggregate the metadata and execute all the actions.
type ScriptOutput struct {
	Metadata      interface{}   `json:"metadata"`       // The generic output returned from the script execution.
	MovedTo       string        `json:"moved_to"`       // Changes the dest for the file
	Children      []string      `json:"children"`       // All files created from this script. children doesn't trigger scripts again.
	Version       int           `json:"version"`        // Script version used in this execution.
	LastRun       time.Time     `json:"last_run"`       // Time which the script runned.
	Log           string        `json:"log"`            // Register of the script stdout and stderr.
	ExecutionTime time.Duration `json:"execution_time"` // Time spent running the script.
}
