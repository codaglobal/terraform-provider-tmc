package tmc

import (
	"context"
	"fmt"
	"time"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClusterGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterGroupsRead,
		Schema: map[string]*schema.Schema{
			"labels": labelsSchema(),
			"ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "UID of the All Tanzu ClusterGroups",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Names of the All Tanzu ClusterGroups",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceClusterGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	var diags diag.Diagnostics

	labels := d.Get("labels").(map[string]interface{})

	res, err := client.GetAllClusterGroups(labels)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read cluster groups",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	clusterGroupNames := make([]interface{}, len(*res))
	clusterGroupIds := make([]interface{}, len(*res))

	for i, clusterGroup := range *res {
		clusterGroupNames[i] = clusterGroup.FullName.Name
		clusterGroupIds[i] = clusterGroup.Meta.UID
	}

	if err := d.Set("names", clusterGroupNames); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ids", clusterGroupIds); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().String())

	return diags
}
