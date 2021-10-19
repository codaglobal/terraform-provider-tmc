package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodePoolCreate,
		ReadContext:   resourceNodePoolRead,
		UpdateContext: resourceNodePoolUpdate,
		DeleteContext: resourceNodePoolDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Nodepool",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Nodepool in the cluster",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !IsValidTanzuName(v) {
						errs = append(errs, fmt.Errorf("name should contain only lowercase letters, numbers or hyphens and should begin with either an alphabet or number"))
					}
					return
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the cluster in which the nodepool is present",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the Nodepool",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the management cluster used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the provisioner",
			},
			"node_labels":  labelsSchemaImmutable(),
			"cloud_labels": labelsSchemaImmutable(),
			"worker_node_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of worker nodes in the nodepool",
			},
			"tkg_aws": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Description: "Availability zone of the worker node",
							Required:    true,
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the worker node",
							Required:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Kubernetes version to be used",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceNodePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	cloud_labels := d.Get("cloud_labels").(map[string]interface{})
	node_labels := d.Get("node_labels").(map[string]interface{})
	cluster_name := d.Get("cluster_name").(string)
	worker_node_count := d.Get("worker_node_count").(int)
	tkg_aws := d.Get("tkg_aws").([]interface{})

	nodepool, err := client.CreateNodePool(npName, managementClusterName, provisionerName, cluster_name, description, cloud_labels, node_labels, worker_node_count, tkg_aws[0].(map[string]interface{}))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create Nodepool Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(nodepool.Meta.UID)

	resourceNodePoolRead(ctx, d, m)

	return diags
}

func resourceNodePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

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
			Detail:   fmt.Sprintf("Error setting spec for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	return diags
}

func resourceNodePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	cloud_labels := d.Get("cloud_labels").(map[string]interface{})
	node_labels := d.Get("node_labels").(map[string]interface{})
	cluster_name := d.Get("cluster_name").(string)
	worker_node_count := d.Get("worker_node_count").(int)
	tkg_aws := d.Get("tkg_aws").([]interface{})

	if d.HasChange("worker_node_count") {
		_, err := client.UpdateNodePool(npName, managementClusterName, provisionerName, cluster_name, description, cloud_labels, node_labels, worker_node_count, tkg_aws[0].(map[string]interface{}))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Update Nodepool Failed",
				Detail:   err.Error(),
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceNodePoolRead(ctx, d, m)
}

func resourceNodePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	npName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	cluster_name := d.Get("cluster_name").(string)

	err := client.DeleteNodePool(npName, cluster_name, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Delete Cluster Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId("")

	return diags
}
