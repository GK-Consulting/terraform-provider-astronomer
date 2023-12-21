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
var _ datasource.DataSource = &WorkspaceDataSource{}

func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

// WorkspaceDataSource defines the data source implementation.
type WorkspaceDataSource struct {
	client *http.Client
}

// WorkspaceDataSourceModel describes the data source data model.
type WorkspaceDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	CicdEnforcedDefault types.Bool   `tfsdk:"cicd_enforced_default"`
	Description         types.String `tfsdk:"description"`
	Name                types.String `tfsdk:"name"`
	OrganizationId      types.String `tfsdk:"organization_id"`
}

func (d *WorkspaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "astronomer_workspace"
}

func (d *WorkspaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Astronomer Workspace Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Workspace Identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Computed:            true,
			},
			"cicd_enforced_default": schema.BoolAttribute{
				MarkdownDescription: "CI CD default",
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

func (d *WorkspaceDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkspaceDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	decoded, err := api.GetWorkspace(data.OrganizationId.ValueString(), data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error"))
		return
	}

	data.Id = types.StringValue(decoded.Id)
	data.CicdEnforcedDefault = types.BoolValue(decoded.CicdEnforcedDefault)
	data.Description = types.StringValue(decoded.Description)
	data.Name = types.StringValue(decoded.Name)
	data.OrganizationId = types.StringValue(decoded.OrganizationId)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
