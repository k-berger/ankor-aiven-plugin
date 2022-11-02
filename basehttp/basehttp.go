package basehttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"
)

const AIVEN_BASE_URL = "https://api.aiven.io/v1/"
const AIVEN_API_TOKEN_ENVVAR = "AIVEN_API_TOKEN"

var client = http.Client{Timeout: 30 * time.Second}

func ExecuteGetRequest(path string) (*http.Response, error) {
	request, err := CreateGetRequest(path)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil

}

func ExecutePutRequest(path string, data map[string]interface{}) (*http.Response, error) {
	request, err := CreatePutRequest(path, data)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil

}
func CreateGetRequest(path string) (*http.Request, error) {
	return createRequest("GET", path, nil)
}

func CreatePostRequest(path string, data map[string]interface{}) (*http.Request, error) {
	req, err := createRequest("POST", path, data)
	if err != nil {
		return nil, err
	}

	// annotate request with body-data to be posted
	return req, nil
}

func CreatePutRequest(path string, data map[string]interface{}) (*http.Request, error) {
	req, err := createRequest("PUT", path, data)
	if err != nil {
		return nil, err
	}

	// annotate request with body-data to be posted

	return req, nil
}
func createRequest(method string, path string, data map[string]interface{}) (*http.Request, error) {

	// create the basic request.
	// early out in case of error.
	var err error
	var req *http.Request

	if data != nil {
		reqBody, marshalError := json.Marshal(data)
		if marshalError != nil {
			return nil, errors.New("Error marshaling request-body!")
		}
		req, err = http.NewRequest(method, AIVEN_BASE_URL+path, bytes.NewBuffer(reqBody))
	} else {
		req, err = http.NewRequest(method, AIVEN_BASE_URL+path, nil)
	}
	if err != nil {
		return nil, err
	}

	// annotate any additional common information on the request
	err = annotateRequestWithAuth(req)
	return req, err
}

func annotateRequestWithAuth(req *http.Request) error {
	token := os.Getenv(AIVEN_API_TOKEN_ENVVAR)
	if token == "" {
		return errors.New("Please provide an AIVEN access-token in " + AIVEN_API_TOKEN_ENVVAR + " environment-variable")
	}
	req.Header.Add("Authorization", "aivenv1 "+token)
	return nil
}
