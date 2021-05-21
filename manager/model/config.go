package model

// Config holds data for the service configuration.
type Config struct {
	LogLevel string                         `json:"log_level"`
	Rules    []Rule                         `json:"rules"`
	Scripts  map[string]ScriptConfiguration `json:"scripts"`
}
