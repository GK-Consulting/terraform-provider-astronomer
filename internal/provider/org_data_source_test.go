package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOrgDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testOrgDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.astronomer_organization.test", "billing_email"),
					resource.TestCheckResourceAttrSet("data.astronomer_organization.test", "id"),
					resource.TestCheckResourceAttrSet("data.astronomer_organization.test", "name"),
				),
			},
		},
	})
}

func testOrgDataSourceConfig() string {
	orgId := os.Getenv("ORGANIZATION_ID")
	return fmt.Sprintf(`
provider "astronomer" {
	organization_id = %[1]q
}

data "astronomer_organization" "test" {
  id = %[1]q
}
`, orgId)
}
