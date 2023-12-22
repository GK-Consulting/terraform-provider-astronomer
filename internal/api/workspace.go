package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WorkspaceCreateRequest struct {
	CicdEnforcedDefault bool   `json:"cicdEnforcedDefault"`
	Description         string `json:"description"`
	Name                string `json:"name"`
}

type WorkspaceUpdateRequest struct {
	CicdEnforcedDefault bool   `json:"cicdEnforcedDefault"`
	Description         string `json:"description"`
	Name                string `json:"name"`
}

type Workspace struct {
	CicdEnforcedDefault bool   `json:"cicdEnforcedDefault"`
	CreatedAt           string `json:"createdAt"`
	CreatedBy           User   `json:"createdBy"`
	Description         string `json:"description"`
	Id                  string `json:"id"`
	Name                string `json:"name"`
	OrganizationId      string `json:"organizationId"`
	OrganizationName    string `json:"organizationName"`
	UpdatedAt           string `json:"updatedAt"`
	UpdatedBy           User   `json:"updatedBy"`
	StatusCode          string `json:"statusCode"`
	Message             string `json:"message"`
	RequestId           string `json:"requestId"`
}

type WorkspaceListResponse struct {
	Workspaces []Workspace `json:"workspaces"`
	TotalCount int         `json:"totalCount"`
	Offset     int         `json:"offset"`
	Limit      int         `json:"limit"`
}

func CreateWorkspace(apiKey string, organizationId string, createRequest *WorkspaceCreateRequest) (*Workspace, error) {
	b, err := json.Marshal(createRequest)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}

	request, _ := http.NewRequest("POST", urlBase+organizationId+"/workspaces", bytes.NewBuffer(b))
	decoded := new(Workspace)
	err = getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}
	return decoded, nil
}

func DeleteWorkspace(apiKey string, organizationId string, workspaceId string) error {
	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", urlBase+organizationId+"/workspaces/"+workspaceId, nil)
	request.Header.Set("Authorization", "Bearer "+apiKey)
	_, err := client.Do(request)

	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}
	return nil
}

func GetWorkspace(apiKey string, organizationId string, workspaceId string) (*Workspace, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", urlBase+organizationId+"/workspaces/"+workspaceId, nil)
	request.Header.Set("Authorization", "Bearer "+apiKey)
	http_resp, err := client.Do(request)

	//TODO what if nothing comes back, add additional checks

	if err != nil {
		return nil, fmt.Errorf("Error: %s", err) //TODO improve error handling
	}
	defer http_resp.Body.Close()
	decoded := new(Workspace)
	json.NewDecoder(http_resp.Body).Decode(&decoded)

	if decoded.Message != "" {
		return nil, fmt.Errorf("%s", decoded.Message)
	}
	return decoded, nil
}

func UpdateWorkspace(apiKey string, organizationId string, workspaceId string, updateRequest *WorkspaceUpdateRequest) (*Workspace, error) {
	//TODO add more checks here
	if workspaceId == "" {
		return nil, fmt.Errorf("No Workspace ID Given.")
	}
	b, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err)
	}

	request, _ := http.NewRequest("POST", urlBase+organizationId+"/workspaces/"+workspaceId, bytes.NewBuffer(b))
	decoded := new(Workspace)
	err = getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return decoded, nil
}
