package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTmcManagementCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmcManagementClusterCreate,
		ReadContext:   resourceTmcManagementClusterRead,
		DeleteContext: resourceTmcManagementClusterDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Cluster Group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Unique Name of the Tanzu Management Cluster in your Org",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !IsValidTanzuName(v) {
						errs = append(errs, fmt.Errorf("name should contain only lowercase letters, numbers or hyphens and should begin with either an alphabet or number"))
					}
					return
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optional Description for the management cluster",
			},
			"labels": labelsSchemaImmutable(),
			"kubernetes_provider_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Indicates the k8s provider Type",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !isValidK8sProviderType(v) {
						errs = append(errs, fmt.Errorf("invalid kubernetes_provider_type specified. it can be one of tkg, tkgservice, tkghosted or other"))
					}
					return
				},
			},
			"default_cluster_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Default cluster group for the workload clusters",
			},
			"registration_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL to fetch the TMC registration YAML.",
			},
		},
	}
}

func resourceTmcManagementClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	mgmtClusterName := d.Get("name").(string)
	k8sProviderType := d.Get("kubernetes_provider_type").(string)
	defaultCg := d.Get("default_cluster_group").(string)
	description := d.Get("description").(string)
	labels := d.Get("labels").(map[string]interface{})

	mgmtCluster, err := client.CreateMgmtCluster(mgmtClusterName, defaultCg, k8sProviderType, description, labels)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create Management cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(mgmtCluster.Meta.UID)

	return resourceTmcManagementClusterRead(ctx, d, m)
}

func resourceTmcManagementClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics
	client := m.(*tanzuclient.Client)

	mgmtClusterName := d.Get("name").(string)

	mgmtCluster, err := client.GetMgmtCluster(mgmtClusterName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read Management cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("description", mgmtCluster.Meta.Description)
	if err := d.Set("labels", mgmtCluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read AWS cluster",
			Detail:   fmt.Sprintf("Error getting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("kubernetes_provider_type", mgmtCluster.Spec.KubernetesProviderType)
	d.Set("default_cluster_group", mgmtCluster.Spec.DefaultClusterGroup)
	d.Set("registration_url", mgmtCluster.Status.RegistrationURL)

	return nil
}

func resourceTmcManagementClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	mgmtClusterName := d.Get("name").(string)

	err := client.DeleteMgmtCluster(mgmtClusterName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete management cluster",
			Detail:   fmt.Sprintf("Error deleting resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return nil
}
