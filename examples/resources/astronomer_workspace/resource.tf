resource "astronomer_workspace" "complete_setup" {
  name                  = "TF Workspace - Standard Deployment"
  cicd_enforced_default = true
  description           = "Testing Workspace"
}
