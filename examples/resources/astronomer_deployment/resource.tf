resource "astronomer_deployment" "standard_deployment" {
  astro_runtime_version   = "9.1.0"
  cloud_provider          = "AWS"
  default_task_pod_cpu    = "0.5"
  default_task_pod_memory = "1Gi"
  description             = "A Standard Deployment"
  executor                = "CELERY"
  is_dag_deploy_enabled   = true
  is_cicd_enforced        = true
  is_high_availability    = true
  name                    = "Test Deployment TF"
  region                  = "us-east-1"
  resource_quota_cpu      = "160"
  resource_quota_memory   = "320Gi"
  scheduler_size          = "MEDIUM"
  type                    = "STANDARD"
  workspace_id            = astronomer_workspace.complete_setup.id
  worker_queues = [
    {
      astro_machine : "A5",
      is_default : true,
      max_worker_count : 1,
      min_worker_count : 1,
      name : "default",
      worker_concurrency : 1,
    },
  ]
  environment_variables = [
    {
      is_secret : true,
      key : "AWS_ACCESS_SECRET_KEY",
      value : "SECRET_VALUE",
    },
    {
      is_secret : false,
      key : "AWS_ACCESS_KEY_ID",
      value : "NOT_SECRET",
    },
  ]
}