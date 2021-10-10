package event

import (
	"time"
)

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
	FatalLevel   = "fatal"
)

type ExtraProps map[string]string

type Event struct {
	Type      string    `json:"type"`
	TimeStamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`

	ExtraProps `json:"extraProps,omitempty"`
}
