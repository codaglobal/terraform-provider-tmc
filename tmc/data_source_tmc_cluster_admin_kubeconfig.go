package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAdminKubeConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAdminKubeConfigRead,
		Schema: map[string]*schema.Schema{
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the cluster in which the nodepool is present",
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
			"kubeconfig": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Base64 encoded string of the kubeconfig for the cluster",
			},
		},
	}
}

func dataSourceAdminKubeConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	var diags diag.Diagnostics

	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	clusterName := d.Get("cluster_name").(string)

	adminKubeconfig, err := client.GetAdminKubeConfig(clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to get admin kubeconfig",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("kubeconfig", adminKubeconfig.KubeConfig)
	d.SetId(string(time.Now().String()))

	return diags
}
