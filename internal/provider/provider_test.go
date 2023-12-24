package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"astronomer": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ASTRONOMER_API_TOKEN"); v == "" {
		t.Fatal("ASTRONOMER_API_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("ORGANIZATION_ID"); v == "" {
		t.Fatal("ORGANIZATION_ID must be set for acceptance tests")
	}
}
