package provider

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/openglshaders/astronomer-api/v2"
)

var _ resource.Resource = &DeploymentResource{}
var _ resource.ResourceWithImportState = &DeploymentResource{}

func NewDeploymentResource() resource.Resource {
	return &DeploymentResource{}
}

type DeploymentResource struct {
	token          string
	organizationId string
}

type DeploymentResourceModel struct {
	AstroRuntimeVersion  types.String       `tfsdk:"astro_runtime_version"`
	CloudProvider        types.String       `tfsdk:"cloud_provider"`
	ClusterId            types.String       `tfsdk:"cluster_id"`
	DefaultTaskPodCpu    types.String       `tfsdk:"default_task_pod_cpu"`
	DefaultTaskPodMemory types.String       `tfsdk:"default_task_pod_memory"`
	Description          types.String       `tfsdk:"description"`
	Executor             types.String       `tfsdk:"executor"`
	Id                   types.String       `tfsdk:"id"`
	IsCicdEnforced       types.Bool         `tfsdk:"is_cicd_enforced"`
	IsDagDeployEnforced  types.Bool         `tfsdk:"is_dag_deploy_enforced"`
	IsHighAvailability   types.Bool         `tfsdk:"is_high_availability"`
	Name                 types.String       `tfsdk:"name"`
	Region               types.String       `tfsdk:"region"`
	ResourceQuotaCpu     types.String       `tfsdk:"resource_quota_cpu"`
	ResourceQuotaMemory  types.String       `tfsdk:"resource_quota_memory"`
	SchedulerSize        types.String       `tfsdk:"scheduler_size"`
	Type                 types.String       `tfsdk:"type"`
	WorkerQueues         []WorkerQueueModel `tfsdk:"worker_queues"`
	WorkloadIdentity     types.String       `tfsdk:"workload_identity"`
	WorkspaceId          types.String       `tfsdk:"workspace_id"`
}

type WorkerQueueModel struct {
	AstroMachine      types.String `tfsdk:"astro_machine"`
	Id                types.String `tfsdk:"id"`
	IsDefault         types.Bool   `tfsdk:"is_default"`
	MaxWorkerCount    types.Int64  `tfsdk:"max_worker_count"`
	MinWorkerCount    types.Int64  `tfsdk:"min_worker_count"`
	Name              types.String `tfsdk:"name"`
	WorkerConcurrency types.Int64  `tfsdk:"worker_concurrency"`
}

func (r *DeploymentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

func (r *DeploymentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "An Astro Deployment is an Airflow environment that is powered by all core Airflow components.",

		Attributes: map[string]schema.Attribute{
			"astro_runtime_version": schema.StringAttribute{
				MarkdownDescription: "Deployment's Astro Runtime version.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider for the Deployment's cluster. Optional if `ClusterId` is specified.",
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cluster_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the cluster where the Deployment will be created.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Deployment's description.",
				Optional:            true,
			},
			"default_task_pod_cpu": schema.StringAttribute{
				MarkdownDescription: "The default CPU resource usage for a worker Pod when running the Kubernetes executor or KubernetesPodOperator. Units are in number of CPU cores.",
				Required:            true,
			},
			"default_task_pod_memory": schema.StringAttribute{
				MarkdownDescription: "The default memory resource usage for a worker Pod when running the Kubernetes executor or KubernetesPodOperator. Units are in `Gi`. This value must always be twice the value of `DefaultTaskPodCpu`.",
				Required:            true,
			},
			"executor": schema.StringAttribute{
				MarkdownDescription: "The Deployment's executor type.",
				Required:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The Deployment's identifier.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_cicd_enforced": schema.BoolAttribute{
				MarkdownDescription: "Whether the Deployment requires that all deploys are made through CI/CD.",
				Required:            true,
			},
			"is_dag_deploy_enforced": schema.BoolAttribute{
				MarkdownDescription: "Whether the Deployment has DAG deploys enabled.",
				Required:            true,
			},
			"is_high_availability": schema.BoolAttribute{
				MarkdownDescription: "Whether the Deployment is configured for high availability. If `true`, multiple scheduler pods will be online.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Deployment's name.",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The region to host the Deployment in. Optional if `ClusterId` is specified.",
				Optional:            true,
			},
			"resource_quota_cpu": schema.StringAttribute{
				MarkdownDescription: "The CPU quota for worker Pods when running the Kubernetes executor or KubernetesPodOperator. If current CPU usage across all workers exceeds the quota, no new worker Pods can be scheduled. Units are in number of CPU cores.",
				Required:            true,
			},
			"resource_quota_memory": schema.StringAttribute{
				MarkdownDescription: "The memory quota for worker Pods when running the Kubernetes executor or KubernetesPodOperator. If current memory usage across all workers exceeds the quota, no new worker Pods can be scheduled. Units are in `Gi`. This value must always be twice the value of `ResourceQuotaCpu`.",
				Required:            true,
			},
			"scheduler_size": schema.StringAttribute{
				MarkdownDescription: "The size of the scheduler pod.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The type of the Deployment.",
				Required:            true,
			},
			"worker_queues": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"astro_machine": schema.StringAttribute{
							Required: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"is_default": schema.BoolAttribute{
							Required: true,
						},
						"max_worker_count": schema.Int64Attribute{
							Required: true,
						},
						"min_worker_count": schema.Int64Attribute{
							Required: true,
						},
						"name": schema.StringAttribute{
							Required: true,
						},
						"worker_concurrency": schema.Int64Attribute{
							Required: true,
						},
					},
				},
				MarkdownDescription: "The list of worker queues configured for the Deployment. Applies only when `Executor` is `CELERY`. At least 1 worker queue is needed. All Deployments need at least 1 worker queue called `default`.",
				Optional:            true,
			},
			"workload_identity": schema.StringAttribute{
				MarkdownDescription: "The Deployment's workload identity.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the workspace to which the Deployment belongs.",
				Required:            true,
			},
		},
	}
}

func (r *DeploymentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*AstronomerProviderResourceDataModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *AstronomerProviderResourceDataModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.token = provider.Token
	r.organizationId = provider.OrganizationId
}

func (r *DeploymentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DeploymentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if data.CloudProvider.ValueString() == "" && data.ClusterId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Validation Error",
			"cluster_id or cloud_provider must be specified",
		)
	}
	if data.Executor.ValueString() == "CELERY" && len(data.WorkerQueues) == 0 {
		resp.Diagnostics.AddError(
			"Validation Error",
			"Must provide at least one default worker queue when using CELERY executor.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	workerQueues := loadWorkerQueuesFromTFState(data)

	//TODO add remaining to model
	deploymentCreateRequest := &api.DeploymentCreateRequest{
		AstroRuntimeVersion:  data.AstroRuntimeVersion.ValueString(),
		ClusterId:            data.ClusterId.ValueString(),
		CloudProvider:        data.CloudProvider.ValueString(),
		DefaultTaskPodCpu:    data.DefaultTaskPodCpu.ValueString(),
		DefaultTaskPodMemory: data.DefaultTaskPodMemory.ValueString(),
		Description:          data.Description.ValueString(),
		Executor:             data.Executor.ValueString(),
		IsCicdEnforced:       data.IsCicdEnforced.ValueBool(),
		IsDagDeployEnabled:   data.IsDagDeployEnforced.ValueBool(),
		IsHighAvailability:   data.IsHighAvailability.ValueBool(),
		Name:                 data.Name.ValueString(),
		Region:               data.Region.ValueString(),
		ResourceQuotaCpu:     data.ResourceQuotaCpu.ValueString(),
		ResourceQuotaMemory:  data.ResourceQuotaMemory.ValueString(),
		// Scheduler: data.Sch,
		SchedulerSize: data.SchedulerSize.ValueString(),
		// TaskPodNodePoolId: data.Task,
		Type:         data.Type.ValueString(),
		WorkerQueues: workerQueues,
		WorkspaceId:  data.WorkspaceId.ValueString(),
	}

	deployResponse, err := api.CreateDeployment(r.token, r.organizationId, deploymentCreateRequest)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	for deployResponse.Status != api.DeploymentStatusHealthy {
		deployResponse, _ = api.GetDeployment(r.token, r.organizationId, deployResponse.Id)
		time.Sleep(1 * time.Second)
	}

	// TODO fill out the rest
	data.CloudProvider = types.StringValue(strings.ToUpper(deployResponse.CloudProvider))
	if data.ClusterId.ValueString() != "" {
		data.ClusterId = types.StringValue(deployResponse.ClusterId)
	}
	// data.DbInstanceType = types.StringValue(deployResponse.CloudProvider)
	data.Id = types.StringValue(deployResponse.Id)
	// data.IsLimited
	// data.Metadata
	data.Name = types.StringValue(deployResponse.Name)
	// data.Node = types.StringValue(deployResponse.Name)
	// data.PodSubnetRange = types.StringValue(deployResponse.OrganizationId)
	// data.ProviderAccount = types.StringValue(deployResponse.OrganizationId)
	if data.Region.ValueString() != "" {
		data.Region = types.StringValue(deployResponse.Region)
	}
	// data.ServicePeeringRange
	// data.ServiceSubnetRange
	// data.Tags
	// data.TenantId
	data.Type = types.StringValue(deployResponse.Type)
	// data.VpcSubnetRange
	workerQueuesDeployment := loadWorkerQueuesFromResponse(deployResponse)
	data.WorkerQueues = workerQueuesDeployment
	data.WorkloadIdentity = types.StringValue(deployResponse.WorkloadIdentity)
	data.WorkspaceId = types.StringValue(deployResponse.WorkspaceId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DeploymentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	deployment, err := api.GetDeployment(r.token, r.organizationId, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	// data.AstroRuntimeVersion = types.StringValue(deployment.Astro)
	data.CloudProvider = types.StringValue(strings.ToUpper(deployment.CloudProvider))
	if data.ClusterId.ValueString() != "" {
		data.ClusterId = types.StringValue(deployment.ClusterId)
	}
	data.Id = types.StringValue(deployment.Id)
	data.DefaultTaskPodCpu = types.StringValue(deployment.DefaultTaskPodCpu)
	data.DefaultTaskPodMemory = types.StringValue(deployment.DefaultTaskPodMemory)
	data.Description = types.StringValue(deployment.Description)
	data.Executor = types.StringValue(deployment.Executor)
	data.IsCicdEnforced = types.BoolValue(deployment.IsCicdEnforced)
	data.IsDagDeployEnforced = types.BoolValue(deployment.IsDagDeployEnabled) //TODO check names on this
	data.IsHighAvailability = types.BoolValue(deployment.IsHighAvailability)

	// data.DbInstanceType = types.StringValue(deployResponse.CloudProvider)
	// data.IsLimited
	// data.Metadata
	data.Name = types.StringValue(deployment.Name)
	// data.Node = types.StringValue(deployResponse.Name)
	// data.PodSubnetRange = types.StringValue(deployResponse.OrganizationId)
	// data.ProviderAccount = types.StringValue(deployResponse.OrganizationId)
	if data.Region.ValueString() != "" || deployment.Region != "" {
		data.Region = types.StringValue(deployment.Region)
	}
	data.ResourceQuotaCpu = types.StringValue(deployment.ResourceQuotaCpu)
	data.ResourceQuotaMemory = types.StringValue(deployment.ResourceQuotaMemory)
	data.SchedulerSize = types.StringValue(deployment.SchedulerSize)

	workerQueues := loadWorkerQueuesFromResponse(deployment)
	data.WorkerQueues = workerQueues

	// data.ServicePeeringRange
	// data.ServiceSubnetRange
	// data.Tags
	// data.TenantId
	data.Type = types.StringValue(deployment.Type)
	// data.VpcSubnetRange
	data.WorkloadIdentity = types.StringValue(deployment.WorkloadIdentity)
	data.WorkspaceId = types.StringValue(deployment.WorkspaceId)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func loadWorkerQueuesFromTFState(data DeploymentResourceModel) []api.WorkerQueue {
	var workerQueues []api.WorkerQueue
	for _, value := range data.WorkerQueues {
		workerQueues = append(workerQueues, api.WorkerQueue{
			AstroMachine:      value.AstroMachine.ValueString(),
			Id:                value.Id.ValueString(),
			IsDefault:         value.IsDefault.ValueBool(),
			MaxWorkerCount:    int(value.MaxWorkerCount.ValueInt64()),
			MinWorkerCount:    int(value.MinWorkerCount.ValueInt64()),
			Name:              value.Name.ValueString(),
			WorkerConcurrency: int(value.WorkerConcurrency.ValueInt64()),
		})
	}
	return workerQueues
}

func loadWorkerQueuesFromResponse(deployment *api.DeploymentResponse) []WorkerQueueModel {
	var workerQueues []WorkerQueueModel
	for _, value := range deployment.WorkerQueues {
		workerQueues = append(workerQueues, WorkerQueueModel{
			AstroMachine:      types.StringValue(value.AstroMachine),
			Id:                types.StringValue(value.Id),
			IsDefault:         types.BoolValue(value.IsDefault),
			MaxWorkerCount:    types.Int64Value(int64(value.MaxWorkerCount)),
			MinWorkerCount:    types.Int64Value(int64(value.MinWorkerCount)),
			Name:              types.StringValue(value.Name),
			WorkerConcurrency: types.Int64Value(int64(value.WorkerConcurrency)),
		})
	}
	return workerQueues
}

func (r *DeploymentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DeploymentResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	workerQueues := loadWorkerQueuesFromTFState(data)
	deploymentUpdateRequest := &api.DeploymentUpdateRequest{
		//TODO add Contact Emails
		DefaultTaskPodCpu:    data.DefaultTaskPodCpu.ValueString(),
		DefaultTaskPodMemory: data.DefaultTaskPodMemory.ValueString(),
		Description:          data.Description.ValueString(),
		EnvironmentVariables: []api.EnvironmentVariableRequest{}, // TODO finish up
		Executor:             data.Executor.ValueString(),
		IsCicdEnforced:       data.IsCicdEnforced.ValueBool(),
		IsDagDeployEnabled:   data.IsDagDeployEnforced.ValueBool(),
		IsHighAvailability:   data.IsHighAvailability.ValueBool(),
		Name:                 data.Name.ValueString(),
		ResourceQuotaCpu:     data.ResourceQuotaCpu.ValueString(),
		ResourceQuotaMemory:  data.ResourceQuotaMemory.ValueString(),
		SchedulerSize:        data.SchedulerSize.ValueString(),
		Type:                 data.Type.ValueString(),
		WorkerQueues:         workerQueues,
		// WorkloadIdentity: data.WorkloadIdentity, // TODO
		WorkspaceId: data.WorkspaceId.ValueString(),
	}

	deployResponse, err := api.UpdateDeployment(r.token, r.organizationId, data.Id.ValueString(), deploymentUpdateRequest)
	log.Println(deployResponse)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
		return
	}

	data.WorkerQueues = loadWorkerQueuesFromResponse(deployResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DeploymentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DeploymentResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := api.DeleteDeployment(r.token, r.organizationId, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
		return
	}
}

func (r *DeploymentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
