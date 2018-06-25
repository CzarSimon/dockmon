package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	docker "docker.io/go-docker"
)

// ErrServiceUnhealthy error inidicating that a service is unhealthy.
var ErrServiceUnhealthy = errors.New("Service unhealthy")

// runHealthChecks creates livenessTargets and starts perpetual
// health check loops for each in different goroutines.
func (env *Env) runHealthChecks(waitGroup *sync.WaitGroup) {
	for _, livenessTarget := range getLivenessTargets(env.serviceOptions) {
		go env.runServiceHealthChecks(livenessTarget, waitGroup)
	}
	log.Println("Waiting")
	time.Sleep(1 * time.Second)
	waitGroup.Wait()
}

// runServiceHealthChecks runs a perpetual health check loop for a given LivenessTarget.
func (env *Env) runServiceHealthChecks(livenessTarget LivenessTarget, waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	for {
		livenessTarget.Wait()
		err := callLivenessTarget(&livenessTarget, env.httpClient)
		if err != nil {
			log.Println(err)
			env.handleLivenessFailure(&livenessTarget)
			continue
		}
		livenessTarget.ClearFailed()
	}
}

// callLivenessTarget performes a health check on a livenessTarget.
func callLivenessTarget(livenessTarget *LivenessTarget, client *http.Client) error {
	resp, err := client.Get(livenessTarget.livenessURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrServiceUnhealthy
	}
	return nil
}

// handleLivenessFailure updates the livenessTarget state and
// restarts the underlying service if needed.
func (env *Env) handleLivenessFailure(livenessTarget *LivenessTarget) {
	livenessTarget.AddFailed()
	if !livenessTarget.ShouldRestart() {
		return
	}
	err := restartService(livenessTarget.serviceName, env.dockerClient, &env.dockerTimeout)
	if err != nil {
		return
	}
	livenessTarget.ClearFailed()
}

// restartService restarts a given service.
func restartService(serviceName string, dockerClient *docker.Client, restartTimeout *time.Duration) error {
	log.Printf("Restarting %s\n", serviceName)
	err := dockerClient.ContainerRestart(context.Background(), serviceName, restartTimeout)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// getLivenessTargets maps a list of LivenessOptions to a list of LivenessTargets
func getLivenessTargets(targetOptions []LivenessOptions) []LivenessTarget {
	targets := make([]LivenessTarget, 0, len(targetOptions))
	for _, opts := range targetOptions {
		targets = append(targets, NewLivenessTarget(opts))
	}
	return targets
}
