package provider

import (
	"context"
	"fmt"
	"net/http"

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
	client *http.Client
}

type WorkspaceResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	CicdEnforcedDefault types.Bool   `tfsdk:"cicd_enforced_default"`
	Description         types.String `tfsdk:"description"`
	Name                types.String `tfsdk:"name"`
	OrganizationId      types.String `tfsdk:"organization_id"`
}

func (d *WorkspaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "astronomer_workspace"
}

func (d *WorkspaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Astronomer Workspace Resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Workspace Identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name",
				Required:            true,
			},
			"cicd_enforced_default": schema.BoolAttribute{
				MarkdownDescription: "CI CD default",
				Optional:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of Workspace",
				Optional:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				Required:            true,
			},
		},
	}
}

func (d *WorkspaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.client = provider.client
}

func (r *WorkspaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkspaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	workspaceCreateRequest := &api.WorkspaceCreateRequest{
		CicdEnforcedDefault: data.CicdEnforcedDefault.ValueBool(),
		Description:         data.Description.ValueString(),
		Name:                data.Name.ValueString(),
	}

	workspace, err := api.CreateWorkspace(data.OrganizationId.ValueString(), workspaceCreateRequest)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create workspace, got error: %+v\n", err))
		return
	}
	data.Id = types.StringValue(workspace.Id)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, fmt.Sprintf("Unable to create workspace, got error: %s", workspace.Id))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *WorkspaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkspaceResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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

func (r *WorkspaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkspaceResourceModel

	// Read Terraform plan data into the model
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
	_, err := api.UpdateWorkspace(data.OrganizationId.ValueString(), data.Id.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkspaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkspaceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := api.DeleteWorkspace(data.OrganizationId.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
		return
	}
}

// TODO figure out what this does
func (r *WorkspaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
