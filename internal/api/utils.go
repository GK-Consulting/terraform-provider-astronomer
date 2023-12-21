package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	RequestId  string `json:"requestId"`
	StatusCode int    `json:"statusCode"`
}

const urlBase string = "https://api.astronomer.io/platform/v1beta1/organizations/"

func makeAuthorizedRequest(req *http.Request, apiKey string) (*http.Response, error) {
	client := &http.Client{}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	httpResp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}

	return httpResp, nil
}

func readErrorFirst(source io.ReadCloser, decoded any) error {
	defer source.Close()

	//Check to see if there was an error response. If so, read the message and error
	errorResponse := new(ErrorResponse)
	b, _ := io.ReadAll(source)
	if err := json.Unmarshal(b, &errorResponse); err != nil {
		return fmt.Errorf("%s", errorResponse.Message)
	}

	if errorResponse.Message != "" {
		return fmt.Errorf("%s", errorResponse.Message)
	}

	if err := json.Unmarshal(b, &decoded); err != nil {
		return fmt.Errorf("%s", err)
	}
	return nil
}

func getObjectFromApi(apiKey string, req *http.Request, decoded any) error { //TODO get objects
	httpResp, httpErr := makeAuthorizedRequest(req, apiKey)

	if httpErr != nil {
		return fmt.Errorf("Request Error: %s", httpErr)
	}

	apiErr := readErrorFirst(httpResp.Body, &decoded)
	if apiErr != nil {
		return fmt.Errorf("API Error: %s", apiErr)
	}
	return nil
}
