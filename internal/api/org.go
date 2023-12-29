package api

import (
	"fmt"
	"net/http"
)

type OrgListResponse struct {
	Organizations []OrgResponse `json:"organizations"`
	Limit         int           `json:"limit"`
	Offset        int           `json:"offset"`
	TotalCount    int           `json:"totalCount"`
}

type OrgResponse struct {
	BillingEmail   string                      `json:"billingEmail"`
	CreatedAt      string                      `json:"createdAt"`
	CreatedBy      BasicSubjectProfileResponse `json:"createdBy"`
	Id             string                      `json:"id"`
	IsScimEnabled  bool                        `json:"isScimEnabled"`
	ManagedDomains []ManagedDomainResponse     `json:"managedDomains"`
	Name           string                      `json:"name"`
	PaymentMethod  string                      `json:"paymentMethod"`
	Product        string                      `json:"product"`
	Status         string                      `json:"status"`
	SupportPlan    string                      `json:"supportPlan"`
	TrialExpiresAt string                      `json:"trialExpiresAt"`
	UpdatedAt      string                      `json:"updatedAt"`
	UpdatedBy      BasicSubjectProfileResponse `json:"updatedBy"`
}

type BasicSubjectProfileResponse struct {
	APITokenName string `json:"apiTokenName"`
	AvatarUrl    string `json:"avatarUrl"`
	FullName     string `json:"fullName"`
	Id           string `json:"id"`
	SubjectType  string `json:"subjectType"`
	Username     string `json:"username"`
}

type ManagedDomainResponse struct {
	CreatedAt      string   `json:"createdAt"`
	EnforcedLogins []string `json:"enforcedLogins"`
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Status         string   `json:"status"`
	UpdatedAt      string   `json:"updatedAt"`
}

func GetOrgs(apiKey string) (*OrgListResponse, error) {
	request, _ := http.NewRequest("GET", urlBase, nil)
	decoded := new(OrgListResponse)
	err := getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return decoded, nil
}

func GetOrg(apiKey string, orgId string) (*OrgResponse, error) {
	request, _ := http.NewRequest("GET", urlBase+orgId, nil)
	decoded := new(OrgResponse)
	err := getObjectFromApi(apiKey, request, &decoded)
	if err != nil {
		return nil, fmt.Errorf("%s", err)
	}

	return decoded, nil
}
