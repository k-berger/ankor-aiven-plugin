// Plugin for ankorstore CLI to manage Aiven services

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/k-berger/ankor-aiven-plugin/aiven"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}
	arguments := os.Args[1:]
	switch arguments[0] {
	case "help", "--help", "-h", "-?":
		printHelp()
	case "projects", "--projects", "-p":
		result, err := aiven.ListProjects()
		if err != nil {
			fmt.Printf("ERROR:" + err.Error())
		}
		for key, value := range result {
			fmt.Println(key + strings.Repeat(" ", 50-len(key)) + value)
		}
	case "services", "--services", "-svc":
		if len(arguments) < 2 {
			fmt.Println("USAGE: goaiven services [projectName]")
			os.Exit(1)
		}
		projectName := arguments[1]
		result, err := aiven.ListServices(projectName)
		if err != nil {
			fmt.Printf("ERROR:" + err.Error())
		}
		renderServiceDescriptionList(result)
	case "start", "--start":
		if len(arguments) < 3 {
			fmt.Println("USAGE: goaiven start [projectName] [serviceName]")
			os.Exit(1)
		}
		projectName := arguments[1]
		serviceName := arguments[2]
		result, err := aiven.StartService(projectName, serviceName)
		if err != nil {
			fmt.Printf("ERROR:" + err.Error())
		}
		renderServiceDescription(result)
	case "stop", "--stop":
		if len(arguments) < 3 {
			fmt.Println("USAGE: goaiven stop [projectName] [serviceName]")
			os.Exit(1)
		}
		projectName := arguments[1]
		serviceName := arguments[2]
		result, err := aiven.StopService(projectName, serviceName)
		if err != nil {
			fmt.Printf("ERROR:" + err.Error())
		}
		renderServiceDescription(result)
	default:
		fmt.Println(arguments[0] + " is an unrecognized command.")
		printHelp()
	}
}

func renderServiceDescription(service *aiven.ServiceDescription) {
	fmt.Println(service.Name + strings.Repeat(" ", 50-len(service.Name)) + service.Type + strings.Repeat(" ", 30-len(service.Type)) + service.State)
}

func renderServiceDescriptionList(services map[string]aiven.ServiceDescription) {
	for _, service := range services {
		renderServiceDescription(&service)
	}
}

func printHelp() {
	fmt.Println("NOTE: In order to use this utility you need to have access to Aiven services. The current version only supports API-token auth. Please create a token for your account and store it in the ENV-VAR AIVEN_API_TOKEN.")
	fmt.Println("")
	fmt.Println("USAGE: aiven projects | services [projectName] | (start|stop) [projectName] [serviceName]")
}
