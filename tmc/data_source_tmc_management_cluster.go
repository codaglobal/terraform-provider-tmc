package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTmcManagementCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcManagementClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster Group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
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
				Computed:    true,
				Description: "Optional Description for the management cluster",
			},
			"labels": labelsSchemaComputed(),
			"kubernetes_provider_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates the k8s provider Type",
			},
			"default_cluster_group": {
				Type:        schema.TypeString,
				Computed:    true,
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

func dataSourceTmcManagementClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	var diags diag.Diagnostics

	mgmtClusterName := d.Get("name").(string)

	mgmtCluster, err := client.GetMgmtCluster(mgmtClusterName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read management cluster",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("description", mgmtCluster.Meta.Description)
	if err := d.Set("labels", mgmtCluster.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster group",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("kubernetes_provider_type", mgmtCluster.Spec.KubernetesProviderType)
	d.Set("default_cluster_group", mgmtCluster.Spec.DefaultClusterGroup)
	d.Set("registration_url", mgmtCluster.Status.RegistrationURL)

	d.SetId(string(mgmtCluster.Meta.UID))

	return diags
}
