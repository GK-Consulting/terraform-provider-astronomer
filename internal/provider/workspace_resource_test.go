package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccWorkspaceResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("astronomer_workspace.test", "name", "one"),
					resource.TestCheckResourceAttr("astronomer_workspace.test", "cicd_enforced_default", "true"),
					resource.TestCheckResourceAttr("astronomer_workspace.test", "description", "TestAcc"),
				),
			},
			{
				ResourceName:      "astronomer_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWorkspaceResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("astronomer_workspace.test", "name", "two"),
				),
			},
		},
	})
}

func testAccWorkspaceResourceConfig(name string) string {
	orgId := os.Getenv("ORGANIZATION_ID")
	return fmt.Sprintf(`
provider "astronomer" {
	organization_id = %[1]q
}
resource "astronomer_workspace" "test" {
  name = %[2]q
  cicd_enforced_default = true
  description = "TestAcc"
}
`, orgId, name)
}
