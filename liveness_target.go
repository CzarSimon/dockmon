package main

import (
	"time"
)

// LivenessTarget service to check for liveness.
type LivenessTarget struct {
	serviceName      string
	livenessURL      string
	livenessInterval time.Duration
	restart          bool
	failAfter        uint8
	failedAttempts   uint8
}

// NewLivenessTarget creates a new LivenessTarget based on the provided options.
func NewLivenessTarget(opts LivenessOptions) LivenessTarget {
	return LivenessTarget{
		serviceName:      opts.ServiceName,
		livenessURL:      opts.LivenessURL,
		livenessInterval: time.Duration(opts.LivenessInterval) * time.Second,
		restart:          opts.Restart,
		failAfter:        opts.FailAfter,
		failedAttempts:   0,
	}
}

// Wait sleeps for the duration of the targets liveness interval.
func (t *LivenessTarget) Wait() {
	time.Sleep(t.livenessInterval)
}

// AddFailed increment the recorded number of failed health checks between restarts.
func (t *LivenessTarget) AddFailed() {
	t.failedAttempts++
}

// ClearFailed sets the number of failed health checks to zero.
func (t *LivenessTarget) ClearFailed() {
	t.failedAttempts = 0
}

// ShouldRestart returns a boolean indicating if a liveness targets service should be restarted.
func (t *LivenessTarget) ShouldRestart() bool {
	return t.restart && t.failedAttempts >= t.failAfter
}
