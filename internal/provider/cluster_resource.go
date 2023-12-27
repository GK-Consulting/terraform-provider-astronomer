package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/openglshaders/astronomer-api/v2"
)

var _ resource.Resource = &ClusterResource{}
var _ resource.ResourceWithImportState = &ClusterResource{}

func NewClusterResource() resource.Resource {
	return &ClusterResource{}
}

type ClusterResource struct {
	token          string
	organizationId string
}

type ClusterMetadataModel struct {
	ExternalIPs   []types.String `tfsdk:"external_ips"`
	OidcIssuerUrl types.String   `tfsdk:"oidc_issuer_url"`
}
type ClusterNodePoolModel struct {
	CloudProvider          types.String   `tfsdk:"cloud_provider"`
	ClusterId              types.String   `tfsdk:"cluster_id"`
	CreatedAt              types.String   `tfsdk:"created_at"`
	Id                     types.String   `tfsdk:"id"`
	IsDefault              types.Bool     `tfsdk:"is_default"`
	MaxNodeCount           types.Int64    `tfsdk:"max_node_count"`
	Name                   types.String   `tfsdk:"name"`
	NodeInstanceType       types.String   `tfsdk:"node_instance_type"`
	SupportedAstroMachines []types.String `tfsdk:"supported_astro_machines"`
	UpdatedAt              types.String   `tfsdk:"updated_at"`
}
type ClusterK8sTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type ClusterResourceModel struct {
	CloudProvider  types.String `tfsdk:"cloud_provider"`
	DbInstanceType types.String `tfsdk:"db_instance_type"`
	Id             types.String `tfsdk:"id"`
	IsLimited      types.Bool   `tfsdk:"is_limited"`
	// Metadata            ClusterMetadataModel   `tfsdk:"metadata"`
	K8sTags             []ClusterK8sTagModel   `tfsdk:"k8s_tags"`
	Name                types.String           `tfsdk:"name"`
	NodePools           []ClusterNodePoolModel `tfsdk:"node_pools"`
	OrganizationId      types.String           `tfsdk:"organization_id"`
	PodSubnetRange      types.String           `tfsdk:"pod_subnet_range"`
	ProviderAccount     types.String           `tfsdk:"provider_account"`
	Region              types.String           `tfsdk:"region"`
	ServicePeeringRange types.String           `tfsdk:"service_peering_range"`
	ServiceSubnetRange  types.String           `tfsdk:"service_subnet_range"`
	TenantId            types.String           `tfsdk:"tenant_id"`
	Type                types.String           `tfsdk:"type"`
	VpcSubnetRange      types.String           `tfsdk:"vpc_subnet_range"`
	WorkspaceIds        []types.String         `tfsdk:"workspace_ids"`
}

func (r *ClusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *ClusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A cluster within an organization. An Astro cluster is a Kubernetes cluster that hosts the infrastructure required to run Deployments.",
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cluster's cloud provider.",
				Required:            true,
			},
			"db_instance_type": schema.StringAttribute{
				MarkdownDescription: "The type of database instance that is used for the cluster. Required for Hybrid clusters.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The cluster's identifier.",
				Computed:            true,
			},
			"is_limited": schema.BoolAttribute{
				MarkdownDescription: "Whether the cluster is limited.",
				Computed:            true,
			},
			"k8s_tags": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "The tag's key.",
							Required:            true,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The tag's value.",
							Required:            true,
						},
					},
				},
				MarkdownDescription: "The Kubernetes tags in the cluster.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The cluster's name.",
				Required:            true,
			},
			"node_pools": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"is_default": schema.BoolAttribute{
							MarkdownDescription: "Whether the node pool is the default node pool of the cluster.",
							Optional:            true,
						},
						"max_node_count": schema.Int64Attribute{
							MarkdownDescription: "The maximum number of nodes that can be created in the node pool.",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the node pool.",
							Required:            true,
						},
						"node_instance_type": schema.StringAttribute{
							MarkdownDescription: "The type of node instance that is used for the node pool.",
							Required:            true,
						},
					},
				},
				MarkdownDescription: "The list of node pools to create in the cluster.",
				Optional:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The organization this cluster is associated with.",
				Computed:            true,
			},
			"pod_subnet_range": schema.StringAttribute{
				MarkdownDescription: "The subnet range for Pods. For GCP clusters only.",
				Optional:            true,
			},
			"provider_account": schema.StringAttribute{
				MarkdownDescription: "The provider account ID. Required for Hybrid clusters.",
				Optional:            true,
			},
			"service_peering_range": schema.StringAttribute{
				MarkdownDescription: "The service peering range. For GCP clusters only.",
				Optional:            true,
			},
			"service_subnet_range": schema.StringAttribute{
				MarkdownDescription: "The service subnet range. For GCP clusters only.",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The cluster's region.",
				Required:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant ID. For Azure clusters only.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The cluster's type.",
				Required:            true,
			},
			"vpc_subnet_range": schema.StringAttribute{
				MarkdownDescription: "The VPC subnet range.",
				Required:            true,
			},
			"workspace_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of Workspaces that are authorized to the cluster.",
				Required:            true,
			},
		},
	}
}

func (r *ClusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func createK8sTagRequestFromTFState(data ClusterResourceModel) []api.ClusterK8sTags {
	var k8sTags []api.ClusterK8sTags
	for _, value := range data.K8sTags {
		k8sTags = append(k8sTags, api.ClusterK8sTags{
			Key:   value.Key.ValueString(),
			Value: value.Value.ValueString(),
		})
	}
	return k8sTags
}

func createK8sTagTFStateFromRequest(tags []api.ClusterK8sTags) []ClusterK8sTagModel {
	var k8sTags []ClusterK8sTagModel
	for _, value := range tags {
		k8sTags = append(k8sTags, ClusterK8sTagModel{
			Key:   types.StringValue(value.Key),
			Value: types.StringValue(value.Value),
		})
	}
	return k8sTags
}

func createNodePoolRequestFromTFState(data ClusterResourceModel) []api.NodePoolRequest {
	var nodePools []api.NodePoolRequest
	for _, value := range data.NodePools {
		nodePools = append(nodePools, api.NodePoolRequest{
			IsDefault:        value.IsDefault.ValueBool(),
			MaxNodeCount:     int(value.MaxNodeCount.ValueInt64()),
			Name:             value.Name.ValueString(),
			NodeInstanceType: value.NodeInstanceType.ValueString(),
		})
	}
	return nodePools
}

func createNodePoolTFStateFromRequest(pools []api.NodePoolResponse) []ClusterNodePoolModel {
	var nodePools []ClusterNodePoolModel
	for _, value := range pools {
		nodePools = append(nodePools, ClusterNodePoolModel{
			CloudProvider:          types.StringValue(value.CloudProvider),
			ClusterId:              types.StringValue(value.ClusterId),
			CreatedAt:              types.StringValue(value.CreatedAt),
			Id:                     types.StringValue(value.Id),
			IsDefault:              types.BoolValue(value.IsDefault),
			MaxNodeCount:           (types.Int64Value(int64(value.MaxNodeCount))),
			Name:                   types.StringValue(value.Name),
			NodeInstanceType:       types.StringValue(value.NodeInstanceType),
			SupportedAstroMachines: createTFStringListFromStrings(value.SupportedAstroMachines),
			UpdatedAt:              types.StringValue(value.UpdatedAt),
		})
	}
	return nodePools
}

func createStringListFromTFState(stringList []types.String) []string {
	var strings []string
	for _, value := range stringList {
		strings = append(strings, value.ValueString())
	}
	return strings
}

func createTFStringListFromStrings(stringList []string) []types.String {
	var strings []types.String
	for _, value := range stringList {
		strings = append(strings, types.StringValue(value))
	}
	return strings
}

func (r *ClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterCreateRequest := &api.ClusterCreateRequest{
		CloudProvider:   data.CloudProvider.ValueString(),
		DbInstanceType:  data.DbInstanceType.ValueString(),
		K8sTags:         createK8sTagRequestFromTFState(data),
		Name:            data.Name.ValueString(),
		NodePools:       createNodePoolRequestFromTFState(data),
		ProviderAccount: data.ProviderAccount.ValueString(),
		Region:          data.Region.ValueString(),
		Type:            data.Type.ValueString(),
		VpcSubnetRange:  data.VpcSubnetRange.ValueString(),
		WorkspaceIds:    createStringListFromTFState(data.WorkspaceIds),
	}

	createResponse, err := api.CreateCluster(r.token, r.organizationId, clusterCreateRequest)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}

	for createResponse.Status != api.ClusterStatusCreated {
		createResponse, _ = api.GetCluster(r.token, r.organizationId, createResponse.Id)
		time.Sleep(1 * time.Second)
	}

	data.CloudProvider = types.StringValue(createResponse.CloudProvider)
	data.DbInstanceType = types.StringValue(createResponse.DbInstanceType)
	data.Id = types.StringValue(createResponse.Id)
	data.IsLimited = types.BoolValue(createResponse.IsLimited)
	// data.Metadata = ClusterMetadataModel{
	// 	ExternalIPs:   createTFStringListFromStrings(createResponse.Metadata.ExternalIPs),
	// 	OidcIssuerUrl: types.StringValue(createResponse.Metadata.OidcIssuerUrl),
	// }
	data.K8sTags = createK8sTagTFStateFromRequest(createResponse.Tags)
	data.Name = types.StringValue(createResponse.Name)
	data.NodePools = createNodePoolTFStateFromRequest(createResponse.NodePools)
	data.OrganizationId = types.StringValue(createResponse.OrganizationId)
	data.PodSubnetRange = types.StringValue(createResponse.PodSubnetRange)
	data.ProviderAccount = types.StringValue(createResponse.ProviderAccount)
	data.Region = types.StringValue(createResponse.Region)
	data.ServicePeeringRange = types.StringValue(createResponse.ServicePeeringRange)
	data.ServiceSubnetRange = types.StringValue(createResponse.ServiceSubnetRange)
	data.TenantId = types.StringValue(createResponse.TenantId)
	data.Type = types.StringValue(createResponse.Type)
	data.VpcSubnetRange = types.StringValue(createResponse.VpcSubnetRange)
	data.WorkspaceIds = createTFStringListFromStrings(createResponse.WorkspaceIds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ClusterResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	// deployment, err := api.GetDeployment(r.token, r.organizationId, data.Id.ValueString())

	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	// 	return
	// }

	// // data.AstroRuntimeVersion = types.StringValue(deployment.Astro)
	// data.CloudProvider = types.StringValue(strings.ToUpper(deployment.CloudProvider))
	// data.Id = types.StringValue(deployment.Id)
	// data.DefaultTaskPodCpu = types.StringValue(deployment.DefaultTaskPodCpu)
	// data.DefaultTaskPodMemory = types.StringValue(deployment.DefaultTaskPodMemory)
	// data.Description = types.StringValue(deployment.Description)
	// data.Executor = types.StringValue(deployment.Executor)
	// data.IsCicdEnforced = types.BoolValue(deployment.IsCicdEnforced)
	// data.IsDagDeployEnforced = types.BoolValue(deployment.IsDagDeployEnabled) //TODO check names on this
	// data.IsHighAvailability = types.BoolValue(deployment.IsHighAvailability)

	// // data.DbInstanceType = types.StringValue(deployResponse.CloudProvider)
	// // data.IsLimited
	// // data.Metadata
	// data.Name = types.StringValue(deployment.Name)
	// // data.Node = types.StringValue(deployResponse.Name)
	// // data.PodSubnetRange = types.StringValue(deployResponse.OrganizationId)
	// // data.ProviderAccount = types.StringValue(deployResponse.OrganizationId)
	// data.Region = types.StringValue(deployment.Region)
	// data.ResourceQuotaCpu = types.StringValue(deployment.ResourceQuotaCpu)
	// data.ResourceQuotaMemory = types.StringValue(deployment.ResourceQuotaMemory)
	// data.SchedulerSize = types.StringValue(deployment.SchedulerSize)

	// workerQueues := loadWorkerQueuesFromResponse(deployment)
	// data.WorkerQueues = workerQueues

	// // data.ServicePeeringRange
	// // data.ServiceSubnetRange
	// // data.Tags
	// // data.TenantId
	// data.Type = types.StringValue(deployment.Type)
	// // data.VpcSubnetRange
	// data.WorkspaceId = types.StringValue(deployment.WorkspaceId)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ClusterResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// workerQueues := loadWorkerQueuesFromTFState(data)
	// deploymentUpdateRequest := &api.DeploymentUpdateRequest{
	// 	//TODO add Contact Emails
	// 	DefaultTaskPodCpu:    data.DefaultTaskPodCpu.ValueString(),
	// 	DefaultTaskPodMemory: data.DefaultTaskPodMemory.ValueString(),
	// 	Description:          data.Description.ValueString(),
	// 	EnvironmentVariables: []api.EnvironmentVariableRequest{}, // TODO finish up
	// 	Executor:             data.Executor.ValueString(),
	// 	IsCicdEnforced:       data.IsCicdEnforced.ValueBool(),
	// 	IsDagDeployEnabled:   data.IsDagDeployEnforced.ValueBool(),
	// 	IsHighAvailability:   data.IsHighAvailability.ValueBool(),
	// 	Name:                 data.Name.ValueString(),
	// 	ResourceQuotaCpu:     data.ResourceQuotaCpu.ValueString(),
	// 	ResourceQuotaMemory:  data.ResourceQuotaMemory.ValueString(),
	// 	SchedulerSize:        data.SchedulerSize.ValueString(),
	// 	Type:                 data.Type.ValueString(),
	// 	WorkerQueues:         workerQueues,
	// 	// WorkloadIdentity: data.WorkloadIdentity, // TODO
	// 	WorkspaceId: data.WorkspaceId.ValueString(),
	// }

	// deployResponse, err := api.UpdateDeployment(r.token, r.organizationId, data.Id.ValueString(), deploymentUpdateRequest)
	// log.Println(deployResponse)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	// 	return
	// }

	// data.WorkerQueues = loadWorkerQueuesFromResponse(deployResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ClusterResourceModel

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

func (r *ClusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
