package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	api "github.com/openglshaders/astronomer-api/v2"
)

var _ datasource.DataSource = &ClusterDataSource{}

func NewClusterDataSource() datasource.DataSource {
	return &ClusterDataSource{}
}

type ClusterDataSource struct {
	token          string
	organizationId string
}

func (d *ClusterDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (d *ClusterDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Astronomer Cluster Data Source",
		Attributes: map[string]schema.Attribute{
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cluster's cloud provider.",
				Computed:            true,
			},
			"db_instance_type": schema.StringAttribute{
				MarkdownDescription: "The type of database instance that is used for the cluster. Required for Hybrid clusters.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The cluster's identifier.",
				Required:            true,
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
				Computed:            true,
			},
			"metadata": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"external_ips": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
					"oidc_issuer_url": schema.StringAttribute{
						Optional: true,
					},
				},
				MarkdownDescription: "The cluster's metadata.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The cluster's name.",
				Computed:            true,
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
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The organization this cluster is associated with.",
				Computed:            true,
			},
			"pod_subnet_range": schema.StringAttribute{
				MarkdownDescription: "The subnet range for Pods. For GCP clusters only.",
				Computed:            true,
			},
			"provider_account": schema.StringAttribute{
				MarkdownDescription: "The provider account ID. Required for Hybrid clusters.",
				Computed:            true,
			},
			"service_peering_range": schema.StringAttribute{
				MarkdownDescription: "The service peering range. For GCP clusters only.",
				Computed:            true,
			},
			"service_subnet_range": schema.StringAttribute{
				MarkdownDescription: "The service subnet range. For GCP clusters only.",
				Computed:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The cluster's region.",
				Computed:            true,
			},
			"tenant_id": schema.StringAttribute{
				MarkdownDescription: "The tenant ID. For Azure clusters only.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The cluster's type.",
				Computed:            true,
			},
			"vpc_subnet_range": schema.StringAttribute{
				MarkdownDescription: "The VPC subnet range.",
				Computed:            true,
			},
			"workspace_ids": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The list of Workspaces that are authorized to the cluster.",
				Computed:            true,
			},
		},
	}
}

func (d *ClusterDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*AstronomerProviderDataSourceDataModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *AstronomerProviderDataSourceDataModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.token = provider.Token
	d.organizationId = provider.OrganizationId
}

func (d *ClusterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ClusterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterResponse, err := api.GetCluster(d.token, d.organizationId, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("ERROR: %s", err.Error()))
		return
	}

	data.CloudProvider = types.StringValue(clusterResponse.CloudProvider)
	data.DbInstanceType = types.StringValue(clusterResponse.DbInstanceType)
	data.Id = types.StringValue(clusterResponse.Id)
	data.IsLimited = types.BoolValue(clusterResponse.IsLimited)

	data.Metadata, _ = getMetadata(clusterResponse)

	data.K8sTags = createK8sTagTFStateFromRequest(clusterResponse.Tags)
	data.Name = types.StringValue(clusterResponse.Name)
	data.NodePools = createNodePoolTFStateFromRequest(clusterResponse.NodePools)
	data.OrganizationId = types.StringValue(clusterResponse.OrganizationId)
	data.PodSubnetRange = types.StringValue(clusterResponse.PodSubnetRange)
	data.ProviderAccount = types.StringValue(clusterResponse.ProviderAccount)
	data.Region = types.StringValue(clusterResponse.Region)
	data.ServicePeeringRange = types.StringValue(clusterResponse.ServicePeeringRange)
	data.ServiceSubnetRange = types.StringValue(clusterResponse.ServiceSubnetRange)
	data.TenantId = types.StringValue(clusterResponse.TenantId)
	data.Type = types.StringValue(clusterResponse.Type)
	data.VpcSubnetRange = types.StringValue(clusterResponse.VpcSubnetRange)
	data.WorkspaceIds = createTFStringListFromStrings(clusterResponse.WorkspaceIds)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func getMetadata(clusterResponse *api.ClusterResponse) (basetypes.ObjectValue, diag.Diagnostics) {
	var values []attr.Value
	for _, ip := range createTFStringListFromStrings(clusterResponse.Metadata.ExternalIPs) {
		values = append(values, ip)
	}
	strs, _ := types.ListValue(types.StringType, values)
	return types.ObjectValue(map[string]attr.Type{
		"external_ips":    types.ListType{ElemType: types.StringType},
		"oidc_issuer_url": types.StringType,
	}, map[string]attr.Value{
		"external_ips":    strs,
		"oidc_issuer_url": types.StringValue(clusterResponse.Metadata.OidcIssuerUrl),
	})
}
