package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeploymentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDeploymentResourceConfig("TestDeployment"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("astronomer_deployment.test", "name", "TestDeployment"),
					resource.TestCheckResourceAttr("astronomer_deployment.test", "astro_runtime_version", "9.1.0"),
					resource.TestCheckResourceAttr("astronomer_deployment.test", "is_high_availability", "true"),
					resource.TestCheckResourceAttr("astronomer_deployment.test", "is_cicd_enforced", "true"),
					resource.TestCheckResourceAttr("astronomer_deployment.test", "is_dag_deploy_enforced", "true"),
				),
			},
			{
				ResourceName:            "astronomer_deployment.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"astro_runtime_version"},
			},
			{
				Config: testDeploymentResourceConfig("TestDeploymentUpdate"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("astronomer_deployment.test", "name", "TestDeploymentUpdate"),
				),
			},
		},
	})
}

func testDeploymentResourceConfig(name string) string {
	orgId := os.Getenv("ORGANIZATION_ID")
	return fmt.Sprintf(`
provider "astronomer" {
	organization_id = %[1]q
}

resource "astronomer_workspace" "test" {
	name = "TestDeploymentWorkspace"
	cicd_enforced_default = true
	description = "TestAccDataSource"
}

resource "astronomer_deployment" "test" {
	astro_runtime_version = "9.1.0"
	cloud_provider = "AWS"
	default_task_pod_cpu = "0.5"
	default_task_pod_memory = "1Gi"
	description = "A Standard Deployment"
	executor = "CELERY"
	is_dag_deploy_enforced = true
	is_cicd_enforced = true
	is_high_availability = true
	name = %[2]q
	region = "us-east-1"
	resource_quota_cpu = "160"
	resource_quota_memory = "320Gi"
	scheduler_size = "MEDIUM"
	type = "STANDARD"
	workspace_id = astronomer_workspace.test.id
	worker_queues = [
		{
		astro_machine:      "A5",
		is_default:         true,
		max_worker_count:    1,
		min_worker_count:    1,
		name:              "default",
		worker_concurrency: 1,
		},
	]
}
`, orgId, name)
}
