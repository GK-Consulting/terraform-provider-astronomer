package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/openglshaders/astronomer-api/v2"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DeploymentDataSource{}

func NewDeploymentDataSource() datasource.DataSource {
	return &DeploymentDataSource{}
}

type DeploymentDataSource struct {
	client *http.Client
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
	OrganizationId types.String `tfsdk:"organization_id"`
}

func (d *DeploymentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "astronomer_deployment"
}

func (d *DeploymentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Astronomer Deployment Resource",

		Attributes: map[string]schema.Attribute{
			"airflow_version": schema.StringAttribute{
				MarkdownDescription: "Airflow Version",
				Computed:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "Cloud Provider",
				Computed:            true,
			},
			"cluster_id": schema.StringAttribute{
				MarkdownDescription: "Cluster Id",
				Computed:            true,
			},
			"cluster_name": schema.StringAttribute{
				MarkdownDescription: "Cluster Name",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Deployment Identifier",
				Required:            true,
			},
			"is_cicd_enforced": schema.BoolAttribute{
				MarkdownDescription: "CI CD default",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of Workspace",
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				Required:            true,
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

	d.client = provider.client
}

func (d *DeploymentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DeploymentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	decoded, err := api.GetDeployment(data.OrganizationId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error"))
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
	data.OrganizationId = types.StringValue(decoded.OrganizationId)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
