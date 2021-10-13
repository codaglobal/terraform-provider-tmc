package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster",
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource version of the Cluster",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Cluster",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Cluster",
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
			"cluster_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cluster group",
			},
			"labels": labelsSchema(),
			"tkg_aws": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Details of Cluster hosted on AWS",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "Region of the AWS Cluster",
							Required:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Kubernetes version to be used",
							Required:    true,
						},
						"credential_name": {
							Type:        schema.TypeString,
							Description: "Kubernetes version of the AWS Cluster",
							Required:    true,
						},
						"availability_zones": {
							Type:        schema.TypeList,
							Description: "Availability zones of the control plane node",
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"instance_type": {
							Type:        schema.TypeString,
							Description: "Instance type used to deploy the control plane node",
							Required:    true,
						},
						"vpc_cidrblock": {
							Type:        schema.TypeString,
							Description: "CIDR block used by the Cluster's VPC",
							Required:    true,
						},
						"pod_cidrblock": {
							Type:        schema.TypeString,
							Description: "CIDR block used by the Cluster's Pods",
							Optional:    true,
							Default:     "192.168.0.0/16",
						},
						"service_cidrblock": {
							Type:        schema.TypeString,
							Description: "CIDR block used by the Cluster's Services",
							Optional:    true,
							Default:     "10.96.0.0/12",
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Description: "Name of the SSH Keypair used in the AWS Cluster",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)
	tkg_aws := d.Get("tkg_aws").([]interface{})

	cluster, err := client.CreateCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, labels, tkg_aws[0].(map[string]interface{}))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Create Cluster Failed",
			Detail:   err.Error(),
		})
		return diags
	}

	d.SetId(cluster.Meta.UID)

	resourceClusterRead(ctx, d, m)

	return diags
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	cluster, err := client.GetCluster(clusterName, managementClusterName, provisionerName)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("resource_version", cluster.Meta.ResourceVersion)
	d.Set("description", cluster.Meta.Description)
	d.Set("cluster_group", cluster.Spec.ClusterGroupName)

	if err := d.Set("labels", cluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	tkgAws := make([]interface{}, 0)

	awsData := flattenAwsData(cluster)

	tkgAws = append(tkgAws, awsData)

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

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})
	cluster_group := d.Get("cluster_group").(string)
	tkg_aws := d.Get("tkg_aws").([]interface{})
	resourceVersion := d.Get("resource_version").(string)

	if d.HasChange("labels") || d.HasChange("cluster_group") {
		_, err := client.UpdateCluster(clusterName, managementClusterName, provisionerName, cluster_group, description, resourceVersion, labels, tkg_aws[0].(map[string]interface{}))
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Update Cluster Failed",
				Detail:   err.Error(),
			})
			return diags
		}

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceClusterRead(ctx, d, m)

}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	err := client.DeleteCluster(clusterName, managementClusterName, provisionerName)
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
