package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	api "github.com/openglshaders/astronomer-api/v2"
)

var _ datasource.DataSource = &OrgDataSource{}

func NewOrgDataSource() datasource.DataSource {
	return &OrgDataSource{}
}

type OrgDataSource struct {
	token          string
	organizationId string
}

type OrgDataSourceModel struct {
	BillingEmail types.String `tfsdk:"billing_email"`
	CreatedAt    types.String `tfsdk:"created_at"`
	// CreatedBy      BasicSubjectProfileModel `tfsdk:"created_by"`
	ID             types.String         `tfsdk:"id"`
	IsScimEnabled  types.Bool           `tfsdk:"is_scim_enabled"`
	ManagedDomains []ManagedDomainModel `tfsdk:"managed_domains"`
	Name           types.String         `tfsdk:"name"`
	PaymentMethod  types.String         `tfsdk:"payment_method"`
	Product        types.String         `tfsdk:"product"`
	Status         types.String         `tfsdk:"status"`
	SupportPlan    types.String         `tfsdk:"support_plan"`
	TrialExpiresAt types.String         `tfsdk:"trial_expires_at"`
	UpdatedAt      types.String         `tfsdk:"updated_at"`
	// UpdatedBy      BasicSubjectProfileModel `tfsdk:"updated_by"`
}

type BasicSubjectProfileModel struct {
	APITokenName types.String `tfsdk:"api_token_name"`
	AvatarUrl    types.String `tfsdk:"avatar_url"`
	FullName     types.String `tfsdk:"full_name"`
	ID           types.String `tfsdk:"id"`
	SubjectType  types.String `tfsdk:"subject_type"`
	Username     types.String `tfsdk:"username"`
}

type ManagedDomainModel struct {
	CreatedAt      types.String   `tfsdk:"created_at"`
	EnforcedLogins []types.String `tfsdk:"enforced_logins"`
	ID             types.String   `tfsdk:"id"`
	Name           types.String   `tfsdk:"name"`
	Status         types.String   `tfsdk:"status"`
	UpdatedAt      types.String   `tfsdk:"updated_at"`
}

func (d *OrgDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization"
}

func (d *OrgDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Astronomer Organization Resource",

		Attributes: map[string]schema.Attribute{

			"billing_email": schema.StringAttribute{
				MarkdownDescription: "Billing email on file for the organization.",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Timestamped string of when this organization was created",
				Computed:            true,
			},
			// "created_by": schema.ObjectAttribute{
			// 	AttributeTypes: map[string]attr.Type{
			// 		"api_token_name": types.StringType,
			// 		"avatar_url":     types.StringType,
			// 		"full_name":      types.StringType,
			// 		"id":             types.StringType,
			// 		"subject_type":   types.StringType,
			// 		"username":       types.StringType,
			// 	},
			// 	MarkdownDescription: "Who created this organization.",
			// 	Computed:            true,
			// },
			"id": schema.StringAttribute{
				MarkdownDescription: "Organization's unique identifier",
				Required:            true,
			},
			"is_scim_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether or not scim is enabled",
				Computed:            true,
			},
			"managed_domains": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_at": schema.StringAttribute{
							Required: true,
						},
						"enforced_logins": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"id": schema.BoolAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
				MarkdownDescription: "List of managed domains (nested)",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Organization's name",
				Computed:            true,
			},
			"payment_method": schema.StringAttribute{
				MarkdownDescription: "Payment method (if set)",
				Optional:            true,
			},
			"product": schema.StringAttribute{
				MarkdownDescription: "Type of astro product (e.g. hosted or hybrid)",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Status of the organization",
				Computed:            true,
			},
			"support_plan": schema.StringAttribute{
				MarkdownDescription: "Type of support plan the organization has",
				Computed:            true,
			},
			"trial_expires_at": schema.StringAttribute{
				MarkdownDescription: "When the trial expires, if organization is in a trial",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Last time the organization was updated",
				Computed:            true,
			},
			// "updated_by": schema.ObjectAttribute{
			// 	AttributeTypes: map[string]attr.Type{
			// 		"api_token_name": types.StringType,
			// 		"avatar_url":     types.StringType,
			// 		"full_name":      types.StringType,
			// 		"id":             types.StringType,
			// 		"subject_type":   types.StringType,
			// 		"username":       types.StringType,
			// 	},
			// 	MarkdownDescription: "Who last updated this organization.",
			// 	Computed:            true,
			// },
		},
	}
}

func (d *OrgDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrgDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	decoded, err := api.GetOrg(d.token, d.organizationId)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("ERROR: %s", err.Error()))
		return
	}

	data.BillingEmail = types.StringValue(decoded.BillingEmail)
	data.CreatedAt = types.StringValue(decoded.CreatedAt)
	// data.CreatedBy = types.Object
	data.ID = types.StringValue(decoded.ID)
	data.IsScimEnabled = types.BoolValue(decoded.IsScimEnabled)
	data.ManagedDomains = loadManagedDomainsFromResponse(decoded)
	data.Name = types.StringValue(decoded.Name)
	data.PaymentMethod = types.StringValue(decoded.PaymentMethod)
	data.Product = types.StringValue(decoded.Product)
	data.Status = types.StringValue(decoded.Status)
	data.SupportPlan = types.StringValue(decoded.SupportPlan)
	data.TrialExpiresAt = types.StringValue(decoded.TrialExpiresAt)
	data.UpdatedAt = types.StringValue(decoded.UpdatedAt)
	// UpdatedBy      BasicSubjectProfile `tfsdk:"updated_by"`

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func loadManagedDomainsFromResponse(org *api.OrgResponse) []ManagedDomainModel {
	var managedDomains []ManagedDomainModel
	for _, value := range org.ManagedDomains {
		managedDomains = append(managedDomains, ManagedDomainModel{
			CreatedAt:      types.StringValue(value.CreatedAt),
			EnforcedLogins: loadEnforcedLoginsFromValues(value.EnforcedLogins),
			ID:             types.StringValue(value.ID),
			Name:           types.StringValue(value.Name),
			Status:         types.StringValue(value.Status),
			UpdatedAt:      types.StringValue(value.UpdatedAt),
		})
	}
	return managedDomains
}

func loadEnforcedLoginsFromValues(values []string) []types.String {
	var enforcedLogins []types.String = []types.String{}
	for _, value := range values {
		enforcedLogins = append(enforcedLogins, types.StringValue(value))
	}

	return enforcedLogins
}
