package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AstronomerProvider satisfies various provider interfaces.
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
	client         *http.Client
}

type AstronomerProviderDataSourceDataModel struct {
	Token          string
	OrganizationId string
	client         *http.Client
}

func (p *AstronomerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "astronomer"
	resp.Version = p.version
}

func (p *AstronomerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required:            false,
				Sensitive:           true,
				MarkdownDescription: "Astronomer API Token",
			},
			"organization_id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Example provider attribute",
			},
		},
	}
}

func (p *AstronomerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AstronomerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// if not token {
	// 	token := os.EnvDefaultFunc("ASTRONOMER_API_TOKEN", nil),
	// }

	client := http.DefaultClient
	dataSourceModel := new(AstronomerProviderDataSourceDataModel)
	dataSourceModel.client = client
	dataSourceModel.Token = data.Token.ValueString()

	resp.DataSourceData = dataSourceModel

	dataModel := new(AstronomerProviderResourceDataModel)
	dataModel.client = client
	dataModel.Token = data.Token.ValueString()

	resp.ResourceData = dataModel
}

func (p *AstronomerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeploymentResource,
		NewWorkspaceResource,
	}
}

func (p *AstronomerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkspaceDataSource,
		NewOrganizationDataSource,
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
