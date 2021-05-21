package models

// ScriptConfiguration represents a script name with parameters, to be executed.
type ScriptConfiguration struct {
	Image      string   `json:"image"`      // Reference to the docker image to be used.
	Name       string   `json:"name"`       // Name of the script to be referenced.
	Parameters []string `json:"parameters"` // Extra parameters to be called with the script execution.
	Version    string   `json:"version"`    // Version of the script to be used.
}
