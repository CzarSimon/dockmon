package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	docker "docker.io/go-docker"
	"github.com/CzarSimon/dockmon/pkg/schema"
)

// ErrServiceUnhealthy error inidicating that a service is unhealthy.
var ErrServiceUnhealthy = errors.New("Service unhealthy")

// runHealthChecks creates livenessTargets and starts perpetual
// health check loops for each in different goroutines.
func (env *Env) runHealthChecks(waitGroup *sync.WaitGroup) {
	for _, livenessTarget := range getLivenessTargets(env.serviceOptions) {
		go env.runServiceHealthChecks(livenessTarget, waitGroup)
	}
	time.Sleep(1 * time.Second)
	waitGroup.Wait()
}

// runServiceHealthChecks runs a perpetual health check loop for a given LivenessTarget.
func (env *Env) runServiceHealthChecks(livenessTarget schema.LivenessTarget, waitGroup *sync.WaitGroup) {
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
		err = env.serviceRepo.SaveHealthSuccess(livenessTarget.ServiceName, now())
		if err != nil {
			log.Println(err)
		}
		livenessTarget.ClearFailed()
	}
}

// callLivenessTarget performes a health check on a livenessTarget.
func callLivenessTarget(livenessTarget *schema.LivenessTarget, client *http.Client) error {
	resp, err := client.Get(livenessTarget.LivenessURL)
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
func (env *Env) handleLivenessFailure(livenessTarget *schema.LivenessTarget) {
	livenessTarget.AddFailed()
	err := env.serviceRepo.SaveHealthFailure(livenessTarget.ServiceName, now())
	if err != nil {
		log.Println(err)
	}
	if !livenessTarget.ShouldRestart() {
		return
	}
	err = restartService(livenessTarget.ServiceName, env.dockerClient, &env.dockerTimeout)
	if err != nil {
		return
	}
	err = env.serviceRepo.SaveRestart(livenessTarget.ServiceName, now())
	if err != nil {
		log.Println(err)
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
func getLivenessTargets(targetOptions []schema.LivenessOptions) []schema.LivenessTarget {
	targets := make([]schema.LivenessTarget, 0, len(targetOptions))
	for _, opts := range targetOptions {
		targets = append(targets, schema.NewLivenessTarget(opts))
	}
	return targets
}

// now returns the current UTC timestamp.
func now() time.Time {
	return time.Now().UTC()
}
