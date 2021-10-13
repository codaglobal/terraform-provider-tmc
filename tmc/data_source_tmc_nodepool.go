package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNodePool() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodePoolRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Nodepool",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Nodepool in the cluster",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cluster in which the nodepool is present",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Nodepool",
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
			"node_labels":  labelsSchemaComputed(),
			"cloud_labels": labelsSchemaComputed(),
			"worker_node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of worker nodes in the nodepool",
			},
			"tkg_aws": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Description: "Availability zone of the worker node",
							Computed:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the worker node",
							Computed:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Kubernetes version to be used",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNodePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)
	var diags diag.Diagnostics

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	cluster_name := d.Get("cluster_name").(string)

	nodepool, err := client.GetNodePool(npName, cluster_name, managementClusterName, provisionerName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("description", nodepool.Meta.Description)
	d.Set("worker_node_count", nodepool.Spec.WorkerNodeCount)

	if err := d.Set("cloud_labels", nodepool.Spec.CloudLabels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read Nodepool",
			Detail:   fmt.Sprintf("Error getting cloud labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	if err := d.Set("node_labels", nodepool.Spec.NodeLabels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read Nodepool",
			Detail:   fmt.Sprintf("Error getting node labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	tkgAws := make([]interface{}, 0)

	awsNodeSpec := flattenAwsNodeSpec(&nodepool.Spec.NodeTkgAws)

	tkgAws = append(tkgAws, awsNodeSpec)

	if err := d.Set("tkg_aws", tkgAws); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster",
			Detail:   fmt.Sprintf("Error getting spec for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(nodepool.Meta.UID)

	return diags
}

func flattenAwsNodeSpec(data *tanzuclient.AwsNodeSpec) map[string]interface{} {
	aws := make(map[string]interface{})

	aws["availability_zone"] = data.AvailabilityZone
	aws["instance_type"] = data.InstanceType
	aws["version"] = data.Version

	return aws
}
