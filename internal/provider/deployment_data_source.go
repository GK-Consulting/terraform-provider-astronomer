package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/openglshaders/astronomer-api/v2"
)

var _ datasource.DataSource = &DeploymentDataSource{}

func NewDeploymentDataSource() datasource.DataSource {
	return &DeploymentDataSource{}
}

type DeploymentDataSource struct {
	token          string
	organizationId string
}

type DeploymentDataSourceModel struct {
	AirflowVersion types.String `tfsdk:"airflow_version"`
	CloudProvider  types.String `tfsdk:"cloud_provider"`
	ClusterId      types.String `tfsdk:"cluster_id"`
	ClusterName    types.String `tfsdk:"cluster_name"`
	Id             types.String `tfsdk:"id"`
	IsCicdEnforced types.Bool   `tfsdk:"is_cicd_enforced"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
}

func (d *DeploymentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deployment"
}

func (d *DeploymentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Astronomer Deployment Resource",

		Attributes: map[string]schema.Attribute{
			"airflow_version": schema.StringAttribute{
				MarkdownDescription: "The Deployment's Astro Runtime version.",
				Computed:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider for the Deployment's cluster. Optional if `ClusterId` is specified.",
				Computed:            true,
			},
			"cluster_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the cluster to which the Deployment will be created in. Optional if cloud provider and region is specified.",
				Computed:            true,
			},
			"cluster_name": schema.StringAttribute{
				MarkdownDescription: "Cluster Name",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The Deployment's Identifier",
				Required:            true,
			},
			"is_cicd_enforced": schema.BoolAttribute{
				MarkdownDescription: "Whether the Deployment requires that all deploys are made through CI/CD.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Deployment's name.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Deployment's description.",
				Computed:            true,
			},
		},
	}
}

func (d *DeploymentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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
}

func (d *DeploymentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeploymentDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	decoded, err := api.GetDeployment(d.token, d.organizationId, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("ERROR: %s", err.Error()))
		return
	}

	data.AirflowVersion = types.StringValue(decoded.AirflowVersion)
	data.CloudProvider = types.StringValue(decoded.CloudProvider)
	data.ClusterId = types.StringValue(decoded.ClusterId)
	data.ClusterName = types.StringValue(decoded.ClusterName)
	data.Description = types.StringValue(decoded.Description)
	data.Id = types.StringValue(decoded.Id)
	data.IsCicdEnforced = types.BoolValue(decoded.IsCicdEnforced)
	data.Name = types.StringValue(decoded.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
