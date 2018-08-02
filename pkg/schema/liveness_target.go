package schema

import (
	"time"
)

// LivenessOptions configuration options for a LivenessTarget.
type LivenessOptions struct {
	ServiceName      string `yaml:"serviceName" json:"serviceName"`
	LivenessURL      string `yaml:"livenessUrl" json:"livenessUrl"`
	LivenessInterval int    `yaml:"livenessInterval" json:"livenessInterval"`
	Restart          bool   `yaml:"restart" json:"restart"`
	FailAfter        uint8  `yaml:"failAfter" json:"failAfter"`
}

// LivenessTarget service to check for liveness.
type LivenessTarget struct {
	ServiceName      string
	LivenessURL      string
	LivenessInterval time.Duration
	Restart          bool
	FailAfter        uint8
	FailedAttempts   uint8
}

// NewLivenessTarget creates a new LivenessTarget based on the provided options.
func NewLivenessTarget(opts LivenessOptions) LivenessTarget {
	return LivenessTarget{
		ServiceName:      opts.ServiceName,
		LivenessURL:      opts.LivenessURL,
		LivenessInterval: time.Duration(opts.LivenessInterval) * time.Second,
		Restart:          opts.Restart,
		FailAfter:        opts.FailAfter,
		FailedAttempts:   0,
	}
}

// Wait sleeps for the duration of the targets liveness interval.
func (t *LivenessTarget) Wait() {
	time.Sleep(t.LivenessInterval)
}

// AddFailed increment the recorded number of failed health checks between restarts.
func (t *LivenessTarget) AddFailed() {
	t.FailedAttempts++
}

// ClearFailed sets the number of failed health checks to zero.
func (t *LivenessTarget) ClearFailed() {
	t.FailedAttempts = 0
}

// ShouldRestart returns a boolean indicating if a liveness targets service should be restarted.
func (t *LivenessTarget) ShouldRestart() bool {
	return t.Restart && t.FailedAttempts >= t.FailAfter
}
