package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ provider.Provider = &AstronomerProvider{}

type AstronomerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type AstronomerProviderModel struct {
	Token          types.String `tfsdk:"token"`
	OrganizationId types.String `tfsdk:"organization_id"`
}

type AstronomerProviderResourceDataModel struct {
	Token          string
	OrganizationId string
}

type AstronomerProviderDataSourceDataModel struct {
	Token          string
	OrganizationId string
}

func (p *AstronomerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "astronomer"
	resp.Version = p.version
}

func (p *AstronomerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "Astronomer API Token. Can be set with an `ASTRONOMER_API_TOKEN` env var.",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Organization id this provider will operate on.",
			},
		},
	}
}

func (p *AstronomerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AstronomerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Token.IsNull() {
		data.Token = types.StringValue(os.Getenv("ASTRONOMER_API_TOKEN"))
	}

	if data.Token.ValueString() == "" {
		tflog.Error(ctx, "No api key provided - either via provider configuration or ASTRONOMER_API_TOKEN environment variable.")
		resp.Diagnostics.AddError(
			"API Key Not Found",
			"No api key provided - either via provider configuration or ASTRONOMER_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	dataSourceModel := new(AstronomerProviderDataSourceDataModel)
	dataSourceModel.Token = data.Token.ValueString()
	dataSourceModel.OrganizationId = data.OrganizationId.ValueString()

	resp.DataSourceData = dataSourceModel

	dataModel := new(AstronomerProviderResourceDataModel)
	dataModel.Token = data.Token.ValueString()
	dataModel.OrganizationId = data.OrganizationId.ValueString()

	resp.ResourceData = dataModel
}

func (p *AstronomerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewClusterResource,
		NewDeploymentResource,
		NewWorkspaceResource,
	}
}

func (p *AstronomerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
		// NewOrganizationDataSource,
		NewClusterDataSource,
		NewDeploymentDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AstronomerProvider{
			version: version,
		}
	}
}
