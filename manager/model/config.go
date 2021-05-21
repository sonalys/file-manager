package model

// Config holds data for the service configuration.
type Config struct {
	LogLevel string                         `json:"log_level"` // debug, info, error, panic ...
	Mounts   map[string]string              `json:"mounts"`    // Used to allow multidisk storage without compromising system files.
	Rules    []Rule                         `json:"rules"`
	Scripts  map[string]ScriptConfiguration `json:"scripts"`
}
