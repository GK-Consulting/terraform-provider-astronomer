package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DeploymentStatusCreating    = "CREATING"
	DeploymentStatusDeploying   = "DEPLOYING"
	DeploymentStatusHealthy     = "HEALTHY"
	DeploymentStatusUnhealthy   = "UNHEALTHY"
	DeploymentStatusUnknown     = "UNKNOWN"
	DeploymentStatusHibernating = "HIBERNATING"
)

const (
	SchedulerSizeSmall  = "SMALL"
	SchedulerSizeMedium = "MEDIUM"
	SchedulerSizeLarge  = "LARGE"
)

const (
	DeploymentTypeDedicated = "DEDICATED"
	DeploymentTypeHybrid    = "HYBRID"
	DeploymentTypeStandard  = "STANDARD"
)

const (
	DeploymentExecutorCelery     = "CELERY"
	DeploymentExecutorKubernetes = "KUBERNETES"
)

type DeploymentListResponse struct {
	Deployments []DeploymentResponse `json:"deployments"`
	Limit       int                  `json:"limit"`
	Offset      int                  `json:"offset"`
	TotalCount  int                  `json:"totalCount"`
}

type EnvironmentVariableResponse struct {
	IsSecret  bool   `json:"isSecret"`
	Key       string `json:"key"`
	UpdatedAt string `json:"updatedAt"`
	Value     string `json:"value"`
}

type EnvironmentVariableRequest struct {
	IsSecret bool   `json:"isSecret"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type DeploymentResponse struct {
	AirflowVersion           string                        `json:"airflowVersion"`
	CloudProvider            string                        `json:"cloudProvider"`
	ClusterId                string                        `json:"clusterId"`
	ClusterName              string                        `json:"clusterName"`
	ContactEmails            []string                      `json:"contactEmails"`
	CreatedAt                string                        `json:"createdAt"`
	CreatedBy                User                          `json:"createdBy"`
	DagTarballVersion        string                        `json:"dagTarballVersion"`
	DefaultTaskPodCpu        string                        `json:"defaultTaskPodCpu"`
	DefaultTaskPodMemory     string                        `json:"defaultTaskPodMemory"`
	Description              string                        `json:"description"`
	EnvironmentVariables     []EnvironmentVariableResponse `json:"environmentVariables"`
	Executor                 string                        `json:"executor"`
	ExternalIPs              []string                      `json:"externalIPs"`
	Id                       string                        `json:"id"`
	ImageRepository          string                        `json:"imageRepository"`
	ImageTag                 string                        `json:"imageTag"`
	ImageVersion             string                        `json:"imageVersion"`
	IsCicdEnforced           bool                          `json:"isCicdEnforced"`
	IsDagDeployEnabled       bool                          `json:"isDagDeployEnabled"`
	IsHighAvailability       bool                          `json:"isHighAvailability"`
	Name                     string                        `json:"name"`
	Namespace                string                        `json:"namespace"`
	OidcIssuerUrl            string                        `json:"oidcIssuerUrl"`
	OrganizationId           string                        `json:"organizationId"`
	Region                   string                        `json:"region"`
	ResourceQuotaCpu         string                        `json:"resourceQuotaCpu"`
	ResourceQuotaMemory      string                        `json:"resourceQuotaMemory"`
	RuntimeVersion           string                        `json:"runtimeVersion"`
	SchedulerAu              int                           `json:"schedulerAu"`
	SchedulerCpu             string                        `json:"schedulerCpu"`
	SchedulerMemory          string                        `json:"schedulerMemory"`
	SchedulerReplicas        int                           `json:"schedulerReplicas"`
	SchedulerSize            string                        `json:"schedulerSize"`
	Status                   string                        `json:"status"`
	StatusReason             string                        `json:"statusReason"`
	TaskPodNodePoolId        string                        `json:"taskPodNodePoolId"`
	Type                     string                        `json:"type"`
	UpdatedAt                string                        `json:"updatedAt"`
	UpdatedBy                User                          `json:"updatedBy"`
	WebServerAirflowApiUrl   string                        `json:"webServerAirflowApiUrl"`
	WebServerCpu             string                        `json:"webServerCpu"`
	WebServerIngressHostname string                        `json:"webServerIngressHostname"`
	WebServerMemory          string                        `json:"webServerMemory"`
	WebServerReplicas        int                           `json:"webServerReplicas"`
	WebServerUrl             string                        `json:"webServerUrl"`
	WorkerQueues             []WorkerQueue                 `json:"workerQueues"`
	WorkloadIdentity         string                        `json:"workloadIdentity"`
	WorkspaceId              string                        `json:"workspaceId"`
	WorkspaceName            string                        `json:"workspaceName"`
}

type WorkerQueue struct {
	AstroMachine      string `json:"astroMachine"`
	Id                string `json:"id"`
	IsDefault         bool   `json:"isDefault"`
	MaxWorkerCount    int    `json:"maxWorkerCount"`
	MinWorkerCount    int    `json:"minWorkerCount"`
	Name              string `json:"name"`
	NodePoolId        string `json:"nodePoolId"`
	PodCpu            string `json:"podCpu"`
	PodMemory         string `json:"podMemory"`
	WorkerConcurrency int    `json:"workerConcurrency"`
}

type SchedulerRequest struct {
	Au       string `json:"au"`
	Replicas int    `json:"replicas"`
}

type DeploymentCreateRequest struct {
	AstroRuntimeVersion  string           `json:"astroRuntimeVersion"`
	CloudProvider        string           `json:"cloudProvider,omitempty"`
	ClusterId            string           `json:"clusterId,omitempty"`
	DefaultTaskPodCpu    string           `json:"defaultTaskPodCpu"`
	DefaultTaskPodMemory string           `json:"defaultTaskPodMemory"`
	Description          string           `json:"description"`
	Executor             string           `json:"executor"`
	IsCicdEnforced       bool             `json:"isCicdEnforced"`
	IsDagDeployEnabled   bool             `json:"isDagDeployEnabled"`
	IsHighAvailability   bool             `json:"isHighAvailability"`
	Name                 string           `json:"name"`
	Region               string           `json:"region,omitempty"`
	ResourceQuotaCpu     string           `json:"resourceQuotaCpu"`
	ResourceQuotaMemory  string           `json:"resourceQuotaMemory"`
	Scheduler            SchedulerRequest `json:"scheduler"`
	SchedulerSize        string           `json:"schedulerSize"`
	TaskPodNodePoolId    string           `json:"taskPodNodePoolId"`
	Type                 string           `json:"type"`
	WorkerQueues         []WorkerQueue    `json:"workerQueues"`
	WorkspaceId          string           `json:"workspaceId"`
}

type DeploymentUpdateRequest struct {
	ContactEmails        []string                     `json:"contactEmails"`
	DefaultTaskPodCpu    string                       `json:"defaultTaskPodCpu"`
	DefaultTaskPodMemory string                       `json:"defaultTaskPodMemory"`
	Description          string                       `json:"description"`
	EnvironmentVariables []EnvironmentVariableRequest `json:"environmentVariables"`
	Executor             string                       `json:"executor"`
	IsCicdEnforced       bool                         `json:"isCicdEnforced"`
	IsDagDeployEnabled   bool                         `json:"isDagDeployEnabled"`
	IsHighAvailability   bool                         `json:"isHighAvailability"`
	Name                 string                       `json:"name"`
	ResourceQuotaCpu     string                       `json:"resourceQuotaCpu"`
	ResourceQuotaMemory  string                       `json:"resourceQuotaMemory"`
	Scheduler            SchedulerRequest             `json:"scheduler"`
	SchedulerSize        string                       `json:"schedulerSize"`
	TaskPodNodePoolId    string                       `json:"taskPodNodePoolId"`
	Type                 string                       `json:"type,omitempty"`
	WorkerQueues         []WorkerQueue                `json:"workerQueues"`
	WorkloadIdentity     string                       `json:"workloadIdentity"`
	WorkspaceId          string                       `json:"workspaceId"`
}

type DeploymentDeleteResponse struct{}

func DeleteDeployment(apiKey string, organizationId string, deploymentId string) error {
	request, _ := http.NewRequest("DELETE", urlBase+organizationId+"/deployments/"+deploymentId, nil)
	_, httpErr := makeAuthorizedRequest(request, apiKey)

	if httpErr != nil {
		return fmt.Errorf("Delete Error: %s", httpErr)
	}
	return nil
}

func GetDeployment(apiKey string, organizationId string, deploymentId string) (*DeploymentResponse, error) {
	request, _ := http.NewRequest("GET", urlBase+organizationId+"/deployments/"+deploymentId, nil)
	decoded := new(DeploymentResponse)
	err := getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return decoded, nil
}

func CreateDeployment(apiKey string, organizationId string, createRequest *DeploymentCreateRequest) (*DeploymentResponse, error) {
	b, err := json.Marshal(createRequest)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err) //TODO improve error handling
	}

	request, _ := http.NewRequest("POST", urlBase+organizationId+"/deployments", bytes.NewBuffer(b))
	decoded := new(DeploymentResponse)
	err = getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return decoded, nil
}

func UpdateDeployment(apiKey string, organizationId string, deploymentId string, updateRequest *DeploymentUpdateRequest) (*DeploymentResponse, error) {
	//TODO add validation etc here
	//TODO consolidate marshalling code
	if deploymentId == "" {
		return nil, fmt.Errorf("No Deployment ID Given.")
	}
	b, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, fmt.Errorf("Error: %s", err) //TODO improve error handling
	}

	request, _ := http.NewRequest("POST", urlBase+organizationId+"/deployments/"+deploymentId, bytes.NewBuffer(b))
	request.Header.Set("Content-Type", "application/json")

	decoded := new(DeploymentResponse)
	err = getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return decoded, nil
}
