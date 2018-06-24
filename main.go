package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Running dockmon")
	env := SetupEnv(getConfig())
	go env.startAPI()

	waitGroup := &sync.WaitGroup{}
	env.runHealthChecks(waitGroup)
}
