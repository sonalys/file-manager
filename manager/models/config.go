package models

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel logrus.Level `json:"log_level"`
}
