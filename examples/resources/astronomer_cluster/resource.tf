resource "astronomer_workspace" "dedicated" {
  name                  = "TF Workspace - Dedicated Deployment"
  cicd_enforced_default = true
  description           = "Workspace that demos a dedicated deployment set up"
}

resource "astronomer_cluster" "aws_dedicated" {
  cloud_provider   = "AWS"
  name             = "Cluster with Updated Name"
  region           = "us-east-1"
  type             = "DEDICATED"
  vpc_subnet_range = "172.20.0.0/20"
  k8s_tags         = []
  node_pools       = []
  workspace_ids    = [astronomer_workspace.dedicated.id]
}
