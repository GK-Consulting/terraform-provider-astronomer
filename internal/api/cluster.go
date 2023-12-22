package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	CloudProviderAws   = "AWS"
	CloudProviderAzure = "AZURE"
	CloudProviderGcp   = "GCP"
)

const (
	ClusterStatusCreating     = "CREATING"
	ClusterStatusCreated      = "CREATED"
	ClusterStatusCreateFailed = "CREATE_FAILED"
	ClusterStatusUpdating     = "UPDATING"
)

const (
	ClusterTypeDedicated = "DEDICATED"
	ClusterTypeHybrid    = "HYBRID"
)

type ClusterListResponse struct {
	Clusters   []ClusterResponse `json:"clusters"`
	Limit      int               `json:"limit"`
	Offset     int               `json:"offset"`
	TotalCount int               `json:"totalCount"`
}

type ClusterMetadataResponse struct {
	ExternalIPs   []string `json:"externalIPs"`
	OidcIssuerUrl []string `json:"oidcIssuerUrl"`
}

type ClusterK8sTags struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NodePoolResponse struct {
	CloudProvider          string   `json:"cloudProvider"`
	ClusterId              string   `json:"clusterId"`
	CreatedAt              string   `json:"createdAt"`
	Id                     string   `json:"id"`
	IsDefault              bool     `json:"isDefault"`
	MaxNodeCount           int      `json:"maxNodeCount"`
	Name                   string   `json:"name"`
	NodeInstanceType       string   `json:"nodeInstanceType"`
	SupportedAstroMachines []string `json:"supportedAstroMachines"`
	UpdatedAt              string   `json:"updatedAt"`
}

type NodePoolRequest struct {
	IsDefault        bool   `json:"isDefault"`
	MaxNodeCount     int    `json:"maxNodeCount"`
	Name             string `json:"name"`
	NodeInstanceType string `json:"nodeInstanceType"`
}

type ClusterResponse struct {
	CloudProvider       string                  `json:"cloudProvider"`
	CreatedAt           string                  `json:"createdAt"`
	DbInstanceType      string                  `json:"dbInstanceType"`
	Id                  string                  `json:"id"`
	IsLimited           bool                    `json:"isLimited"`
	Metadata            ClusterMetadataResponse `json:"metadata"`
	Name                string                  `json:"name"`
	NodePools           []NodePoolResponse      `json:"nodePools"`
	OrganizationId      string                  `json:"organizationId"`
	PodSubnetRange      string                  `json:"podSubnetRange"`
	ProviderAccount     string                  `json:"providerAccount"`
	Region              string                  `json:"region"`
	ServicePeeringRange string                  `json:"servicePeeringRange"`
	ServiceSubnetRange  string                  `json:"serviceSubnetRange"`
	Status              string                  `json:"status"`
	Tags                []ClusterK8sTags        `json:"tags"`
	TenantId            string                  `json:"tenantId"`
	Type                string                  `json:"type"`
	UpdatedAt           string                  `json:"updatedAt"`
	VpcSubnetRange      string                  `json:"vpcSubnetRange"`
	WorkspaceIds        []string                `json:"workspaceIds"`
}

type ClusterCreateRequest struct {
	CloudProvider   string            `json:"cloudProvider"`
	DbInstanceType  string            `json:"dbInstanceType"`
	K8sTags         []ClusterK8sTags  `json:"k8sTags"`
	Name            string            `json:"name"`
	NodePools       []NodePoolRequest `json:"nodePools"`
	ProviderAccount string            `json:"providerAccount"`
	Region          string            `json:"region"`
	Type            string            `json:"type"`
	VpcSubnetRange  string            `json:"vpcSubnetRange"`
	WorkspaceIds    []string          `json:"workspaceIds"`
}

type ClusterUpdateRequest struct {
	DbInstanceType string            `json:"dbInstanceType"`
	K8sTags        []ClusterK8sTags  `json:"k8sTags"`
	Name           string            `json:"name"`
	NodePools      []NodePoolRequest `json:"nodePools"`
	WorkspaceIds   []string          `json:"workspaceIds"`
}

type ClusterDeleteResponse struct{}

func GetCluster(apiKey string, organizationId string, clusterId string) (*ClusterResponse, error) {
	request, _ := http.NewRequest("GET", urlBase+organizationId+"/clusters/"+clusterId, nil)
	decoded := new(ClusterResponse)
	err := getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}
	return decoded, nil
}

func CreateCluster(apiKey string, organizationId string, createRequest *ClusterCreateRequest) (*ClusterResponse, error) {
	b, err := json.Marshal(createRequest)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	request, _ := http.NewRequest("POST", urlBase+organizationId+"/clusters", bytes.NewBuffer(b))
	decoded := new(ClusterResponse)

	err = getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return decoded, nil
}

func UpdateCluster(apiKey string, organizationId string, clusterId string, updateRequest *ClusterUpdateRequest) (*ClusterResponse, error) {
	b, err := json.Marshal(updateRequest)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	request, _ := http.NewRequest("POST", urlBase+organizationId+"/clusters/"+clusterId, bytes.NewBuffer(b))
	decoded := new(ClusterResponse)

	err = getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return decoded, nil
}

func DeleteCluster(apiKey string, organizationId string, clusterId string) error {
	request, _ := http.NewRequest("DELETE", urlBase+organizationId+"/clusters/"+clusterId, nil)
	_, httpErr := makeAuthorizedRequest(request, apiKey)

	if httpErr != nil {
		return fmt.Errorf("Delete Error: %s", httpErr)
	}
	return nil
}
