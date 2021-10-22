package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAwsCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAwsClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Cluster",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Cluster",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the management cluster used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the provisioner",
			},
			"cluster_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the cluster group",
			},
			"labels": labelsSchemaComputed(),
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the AWS Cluster",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Provisioner credential used to create the cluster",
				Computed:    true,
			},
			"credential_name": {
				Type:        schema.TypeString,
				Description: "Kubernetes version of the AWS Cluster",
				Computed:    true,
			},
			"availability_zones": {
				Type:        schema.TypeList,
				Description: "Availability zones of the control plane node",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"instance_type": {
				Type:        schema.TypeString,
				Description: "Instance type used to deploy the control plane node",
				Computed:    true,
			},
			"vpc_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's VPC",
				Computed:    true,
			},
			"ssh_key": {
				Type:        schema.TypeString,
				Description: "Name of the SSH Keypair used in the AWS Cluster",
				Computed:    true,
			},
			"pod_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Pods",
				Computed:    true,
			},
			"service_cidrblock": {
				Type:        schema.TypeString,
				Description: "CIDR block used by the Cluster's Services",
				Computed:    true,
			},
		},
	}
}

func dataSourceAwsClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	var diags diag.Diagnostics

	cluster, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("description", cluster.Meta.Description)
	d.Set("cluster_group", cluster.Spec.ClusterGroupName)

	if err := d.Set("labels", cluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("availability_zones", cluster.Spec.TkgAws.Topology.ControlPlane.AvailabilityZones)
	d.Set("instance_type", cluster.Spec.TkgAws.Topology.ControlPlane.InstanceType)
	d.Set("vpc_cidrblock", cluster.Spec.TkgAws.Settings.Network.Provider.Vpc.CidrBlock)
	d.Set("region", cluster.Spec.TkgAws.Distribution.Region)
	d.Set("credential_name", cluster.Spec.TkgAws.Distribution.ProvisionerCredentialName)
	d.Set("version", cluster.Spec.TkgAws.Distribution.Version)
	d.Set("ssh_key", cluster.Spec.TkgAws.Settings.Security.SshKey)
	d.Set("pod_cidrblock", cluster.Spec.TkgAws.Settings.Network.ClusterNetwork.Pods[0].CidrBlocks)
	d.Set("service_cidrblock", cluster.Spec.TkgAws.Settings.Network.ClusterNetwork.Services[0].CidrBlocks)

	d.SetId(string(cluster.Meta.UID))

	return diags
}
