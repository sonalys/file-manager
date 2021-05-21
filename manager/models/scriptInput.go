package models

import "time"

// ScriptInput serves as generic script input data structure.
type ScriptInput struct {
	Filename  string        `json:"filename"`  // Path for the source file.
	Arguments []string      `json:"arguments"` // Arguments for the script process.
	Timeout   time.Duration `json:"timeout"`   // Max execution time for the script to run.
}
