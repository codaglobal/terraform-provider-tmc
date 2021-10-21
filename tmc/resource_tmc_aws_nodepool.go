package tmc

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsNodePoolCreate,
		ReadContext:   resourceAwsNodePoolRead,
		UpdateContext: resourceAwsNodePoolUpdate,
		DeleteContext: resourceAwsNodePoolDelete,
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
	}
}

func resourceAwsNodePoolCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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

	awsNodeSpec := &tanzuclient.AwsNodeSpec{
		AvailabilityZone: d.Get("availability_zone").(string),
		Version:          d.Get("version").(string),
		InstanceType:     d.Get("instance_type").(string),
	}

	nodepool, err := client.CreateNodePool(npName, managementClusterName, provisionerName, cluster_name, description, cloud_labels, node_labels, worker_node_count, awsNodeSpec)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create Nodepool Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(nodepool.Meta.UID)

	resourceAwsNodePoolRead(ctx, d, m)

	return diags
}

func resourceAwsNodePoolRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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

	nodeCount, _ := strconv.Atoi(nodepool.Spec.WorkerNodeCount)

	d.Set("description", nodepool.Meta.Description)
	d.Set("worker_node_count", nodeCount)

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

	d.Set("availability_zone", nodepool.Spec.NodeTkgAws.AvailabilityZone)
	d.Set("instance_type", nodepool.Spec.NodeTkgAws.InstanceType)
	d.Set("version", nodepool.Spec.NodeTkgAws.Version)

	return diags
}

func resourceAwsNodePoolUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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

	awsNodeSpec := &tanzuclient.AwsNodeSpec{
		AvailabilityZone: d.Get("availability_zone").(string),
		Version:          d.Get("version").(string),
		InstanceType:     d.Get("instance_type").(string),
	}

	if d.HasChange("worker_node_count") {
		_, err := client.UpdateNodePool(npName, managementClusterName, provisionerName, cluster_name, description, cloud_labels, node_labels, worker_node_count, awsNodeSpec)
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

	return resourceAwsNodePoolRead(ctx, d, m)
}

func resourceAwsNodePoolDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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
