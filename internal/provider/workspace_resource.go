package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/openglshaders/astronomer-api/v2"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &WorkspaceResource{}

func NewWorkspaceResource() resource.Resource {
	return &WorkspaceResource{}
}

type WorkspaceResource struct {
	token          string
	organizationId string
}

type WorkspaceResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	CicdEnforcedDefault types.Bool   `tfsdk:"cicd_enforced_default"`
	Description         types.String `tfsdk:"description"`
	Name                types.String `tfsdk:"name"`
}

func (r *WorkspaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workspace"
}

func (r *WorkspaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Astronomer Workspace Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Workspace's identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The Workspace's name.",
				Required:            true,
			},
			"cicd_enforced_default": schema.BoolAttribute{
				MarkdownDescription: "Whether new Deployments enforce CI/CD deploys by default.",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The Workspace's description.",
				Optional:            true,
			},
		},
	}
}

func (r *WorkspaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	provider, ok := req.ProviderData.(*AstronomerProviderResourceDataModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *AstronomerProviderResourceDataModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.token = provider.Token
	r.organizationId = provider.OrganizationId
}

func (r *WorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkspaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceCreateRequest := &api.WorkspaceCreateRequest{
		CicdEnforcedDefault: data.CicdEnforcedDefault.ValueBool(),
		Description:         data.Description.ValueString(),
		Name:                data.Name.ValueString(),
	}

	workspace, err := api.CreateWorkspace(r.token, r.organizationId, workspaceCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace, got error: %+v\n", err))
		return
	}
	data.Id = types.StringValue(workspace.Id)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkspaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	decoded, err := api.GetWorkspace(r.token, r.organizationId, data.Id.ValueString())

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

func (r *WorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkspaceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	tflog.Debug(ctx, "")
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &api.WorkspaceUpdateRequest{
		CicdEnforcedDefault: data.CicdEnforcedDefault.ValueBool(),
		Description:         data.Description.ValueString(),
		Name:                data.Name.ValueString(),
	}
	_, err := api.UpdateWorkspace(r.token, r.organizationId, data.Id.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkspaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := api.DeleteWorkspace(r.token, r.organizationId, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
		return
	}
}

func (r *WorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
