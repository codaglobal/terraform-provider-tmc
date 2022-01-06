package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTmcNamespace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcNamespaceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Cluster Group",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique Name of the Namespace in your Org",
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
				Description: "Name of the cluster in which the namespace is to be created",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the Namespace",
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
			"labels": labelsSchemaComputed(),
			"workspace_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the workspace for the namespace",
			},
		},
	}
}

func dataSourceTmcNamespaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	NamespaceName := d.Get("name").(string)
	clusterName := d.Get("cluster_name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	Namespace, err := client.GetNamespace(NamespaceName, clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read namespace",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("description", Namespace.Meta.Description)
	if err := d.Set("labels", Namespace.Meta.Labels); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read namespace",
			Detail:   fmt.Sprintf("Error setting labels for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("workspace_name", Namespace.Spec.WorkspaceName)

	d.SetId(Namespace.Meta.UID)

	return nil
}
