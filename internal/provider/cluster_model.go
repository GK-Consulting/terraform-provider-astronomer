package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ClusterModel struct {
	CloudProvider       types.String           `tfsdk:"cloud_provider"`
	DbInstanceType      types.String           `tfsdk:"db_instance_type"`
	Id                  types.String           `tfsdk:"id"`
	IsLimited           types.Bool             `tfsdk:"is_limited"`
	Metadata            types.Object           `tfsdk:"metadata"`
	K8sTags             []ClusterK8sTagModel   `tfsdk:"k8s_tags"`
	Name                types.String           `tfsdk:"name"`
	NodePools           []ClusterNodePoolModel `tfsdk:"node_pools"`
	OrganizationId      types.String           `tfsdk:"organization_id"`
	PodSubnetRange      types.String           `tfsdk:"pod_subnet_range"`
	ProviderAccount     types.String           `tfsdk:"provider_account"`
	Region              types.String           `tfsdk:"region"`
	ServicePeeringRange types.String           `tfsdk:"service_peering_range"`
	ServiceSubnetRange  types.String           `tfsdk:"service_subnet_range"`
	TenantId            types.String           `tfsdk:"tenant_id"`
	Type                types.String           `tfsdk:"type"`
	VpcSubnetRange      types.String           `tfsdk:"vpc_subnet_range"`
	WorkspaceIds        []types.String         `tfsdk:"workspace_ids"`
}

type ClusterNodePoolModel struct {
	CloudProvider          types.String   `tfsdk:"cloud_provider"`
	ClusterId              types.String   `tfsdk:"cluster_id"`
	CreatedAt              types.String   `tfsdk:"created_at"`
	Id                     types.String   `tfsdk:"id"`
	IsDefault              types.Bool     `tfsdk:"is_default"`
	MaxNodeCount           types.Int64    `tfsdk:"max_node_count"`
	Name                   types.String   `tfsdk:"name"`
	NodeInstanceType       types.String   `tfsdk:"node_instance_type"`
	SupportedAstroMachines []types.String `tfsdk:"supported_astro_machines"`
	UpdatedAt              types.String   `tfsdk:"updated_at"`
}
type ClusterK8sTagModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}
