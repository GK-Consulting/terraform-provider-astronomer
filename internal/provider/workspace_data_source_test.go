package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestWorkspaceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testWorkspaceDataSourceConfig("Data Source Test Workspace"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.astronomer_workspace.test", "name", "Data Source Test Workspace"),
					resource.TestCheckResourceAttr("data.astronomer_workspace.test", "cicd_enforced_default", "true"),
					resource.TestCheckResourceAttr("data.astronomer_workspace.test", "description", "TestAccDataSource"),
				),
			},
		},
	})
}

func testWorkspaceDataSourceConfig(workspaceName string) string {
	orgId := os.Getenv("ORGANIZATION_ID")
	return fmt.Sprintf(`
terraform {
	required_providers {
		astronomer = {
			source = "registry.terraform.io/gk-consulting/astronomer"
		}
	}
}
provider "astronomer" {
	organization_id = %[1]q
}


resource "astronomer_workspace" "test" {
	name = %[2]q
	cicd_enforced_default = true
	description = "TestAccDataSource"
}

data "astronomer_workspace" "test" {
	id = astronomer_workspace.test.id
}
`, orgId, workspaceName)
}
