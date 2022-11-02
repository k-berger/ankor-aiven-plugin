package aiven

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/k-berger/ankor-aiven-plugin/basehttp"
)

type ServiceDescription struct {
	Type  string
	Name  string
	State string
}

func ListProjects() (map[string]string, error) {
	response, err := basehttp.ExecuteGetRequest("project")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// something went wrong. Aiven does not only sent a corresponding http-status in that case
		// but also write messages into the response-body, which we can extract and use for creating
		// an error that can actually help users in understanding the issue.
		return nil, errors.New("Aiven returned " + response.Status + " so deal with it. It won't work.")

	} else {

		result := make(map[string]string)
		var data = make(map[string]interface{})
		json.NewDecoder(response.Body).Decode(&data)
		for key, value := range data["project_membership"].(map[string]interface{}) {
			result[key] = value.(string)
		}
		return result, nil

	}
}

func ListServices(projectName string) (map[string]ServiceDescription, error) {
	response, err := basehttp.ExecuteGetRequest("project/" + projectName + "/service")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		// something went wrong. Aiven does not only sent a corresponding http-status in that case
		// but also write messages into the response-body, which we can extract and use for creating
		// an error that can actually help users in understanding the issue.
		return nil, errors.New("Aiven returned " + response.Status + " so deal with it. It won't work.")

	} else {

		result := make(map[string]ServiceDescription)
		var data = make(map[string]interface{})
		json.NewDecoder(response.Body).Decode(&data)
		for _, value := range data["services"].([]interface{}) {
			accessor := value.(map[string]interface{})
			result[accessor["service_name"].(string)] = transformToServiceDescription(accessor)
		}
		return result, nil

	}
}

func transformToServiceDescription(data map[string]interface{}) ServiceDescription {
	return ServiceDescription{
		Name:  data["service_name"].(string),
		Type:  data["service_type"].(string),
		State: data["state"].(string)}
}

// Stops an existing service without deleting it (power off)
// The function will return a ServiceDescription with the state of the service right after
// sending the stop signal or an error if something goes wrong.
// Only running services can be powered off.
func StopService(projectName string, serviceName string) (*ServiceDescription, error) {
	inputData := make(map[string]interface{})
	inputData["powered"] = false

	response, err := updateService(projectName, serviceName, inputData)
	return response, err
}

// Starts an existing service (power on)
// The function will return a ServiceDescription with the state of the service right after
// sending the stop signal or an error if something goes wrong. Only powered off services
// can be started.
func StartService(projectName string, serviceName string) (*ServiceDescription, error) {
	inputData := make(map[string]interface{})
	inputData["powered"] = true

	response, err := updateService(projectName, serviceName, inputData)
	return response, err
}

func updateService(projectName string, serviceName string, inputData map[string]interface{}) (*ServiceDescription, error) {

	response, err := basehttp.ExecutePutRequest("project/"+projectName+"/service/"+serviceName, inputData)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Could not stop service, API returned " + response.Status)
	}

	// looks like this worked. Now we will query the service-statuses for this project
	// and return the status for the service that should be stopped.

	allServices, err := ListServices(projectName)
	if err != nil {
		return nil, err
	}

	if val, ok := allServices[serviceName]; ok {
		return &val, nil
	}
	return nil, errors.New("I don't know what to do man, seriously.")
}
