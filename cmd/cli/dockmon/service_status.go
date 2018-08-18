package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/CzarSimon/dockmon/pkg/schema"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

// GetServicesCommand returns command for listing services tracked by dockmon.
func GetServicesCommand() cli.Command {
	return cli.Command{
		Name:   "get-services",
		Usage:  fmt.Sprintf("Lists the services tracked by dockmon and their status"),
		Action: GetServices,
	}
}

// GetServiceCommand returns command for describing a specified service tracked by dockmon.
func GetServiceCommand() cli.Command {
	return cli.Command{
		Name:   "get-service",
		Usage:  fmt.Sprintf("Describes a sepecified service tracked by dockmon"),
		Action: GetService,
	}
}

// GetServices display the list of services tracked by dockmon.
func GetServices(c *cli.Context) error {
	api := GetApiClientAndTestCredentials()
	services := api.GetStatuses()
	printServicesList(services)

	return nil
}

// GetService describes a sepecified service tracked by dockmon".
func GetService(c *cli.Context) error {
	serviceName := getServiceName(c)
	api := GetApiClientAndTestCredentials()

	service := api.GetStatus(serviceName)
	svcJSON, _ := json.MarshalIndent(service, "", "    ")
	fmt.Println(string(svcJSON))

	return nil
}

func getServiceName(c *cli.Context) string {
	serviceName := c.Args().First()
	if serviceName == "" {
		fmt.Println("No service name provided")
		os.Exit(1)
	}
	return serviceName
}

func printServicesList(services []schema.ServiceStatus) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Status", "Should Restart", "Restarts", "Age"})
	for _, svc := range services {
		table.Append(makeServiceRow(svc))
	}
	table.Render()
}

func makeServiceRow(svc schema.ServiceStatus) []string {
	return []string{
		svc.ServiceName,
		selectString(svc.IsHealty, "healthy", "unhealthy"),
		selectString(svc.ShouldRestart, "Yes", "No"),
		fmt.Sprintf("%d", svc.Restarts),
		makeAgeString(svc.CreatedAt),
	}
}

func selectString(selector bool, trueOption, falseOption string) string {
	if selector {
		return trueOption
	}
	return falseOption
}

func makeAgeString(fromTime time.Time) string {
	age := time.Now().UTC().Sub(fromTime)

	HOURS_IN_DAY := 24
	HOURS_IN_YEAR := HOURS_IN_DAY * 365
	DAY := 24 * time.Hour
	YEAR := 365 * DAY
	if age < time.Minute {
		return fmt.Sprintf("%d seconds", int(age.Seconds()))
	} else if age < time.Hour {
		return fmt.Sprintf("%d minutes", int(age.Minutes()))
	} else if age < DAY {
		return fmt.Sprintf("%d hours", int(age.Hours()))
	} else if age > YEAR {
		return fmt.Sprintf("%d years", int(age.Hours())/HOURS_IN_YEAR)
	}
	return fmt.Sprintf("%d days", int(age.Hours())/HOURS_IN_DAY)
}
