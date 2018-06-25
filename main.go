package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("Running dockmon")
	env := SetupEnv(getConfig())
	defer env.Close()

	go env.startAPI()
	waitGroup := &sync.WaitGroup{}
	env.runHealthChecks(waitGroup)
}
