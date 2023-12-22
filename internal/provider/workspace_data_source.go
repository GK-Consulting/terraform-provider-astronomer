package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/openglshaders/astronomer-api/v2"
)

var _ datasource.DataSource = &WorkspaceDataSource{}

func NewWorkspaceDataSource() datasource.DataSource {
	return &WorkspaceDataSource{}
}

type WorkspaceDataSource struct {
	token          string
	organizationId string
}

type WorkspaceDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	CicdEnforcedDefault types.Bool   `tfsdk:"cicd_enforced_default"`
	Description         types.String `tfsdk:"description"`
	Name                types.String `tfsdk:"name"`
}

func (d *WorkspaceDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (d *WorkspaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Astronomer Workspace Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Workspace's identifier.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Workspace's name",
				Computed:            true,
			},
			"cicd_enforced_default": schema.BoolAttribute{
				MarkdownDescription: "Whether new Deployments enforce CI/CD deploys by default.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Workspace's description",
				Computed:            true,
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

	d.token = provider.Token
	d.organizationId = provider.OrganizationId
}

func (d *WorkspaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data WorkspaceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	decoded, err := api.GetWorkspace(d.token, d.organizationId, data.Id.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf(err.Error()))
		return
	}

	data.Id = types.StringValue(decoded.Id)
	data.CicdEnforcedDefault = types.BoolValue(decoded.CicdEnforcedDefault)
	data.Description = types.StringValue(decoded.Description)
	data.Name = types.StringValue(decoded.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
