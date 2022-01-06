package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTmcNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmcNamespaceCreate,
		ReadContext:   resourceTmcNamespaceRead,
		DeleteContext: resourceTmcNamespaceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Namespace",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
				ForceNew:    true,
				Description: "Name of the cluster in which the namespace is to be created",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Description of the Namespace",
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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Ignore changes to the creator label added automatically added by TMC and
					// also ignore changes when the labels field itself is deleted when updating
					return k == "labels.tmc.cloud.vmware.com/creator" || k == "labels.%" || k == "labels.tmc.cloud.vmware.com/managed" || k == "labels.tmc.cloud.vmware.com/workspace"
				},
			},
			"workspace_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the workspace for the namespace",
			},
		},
	}
}

func resourceTmcNamespaceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	NamespaceName := d.Get("name").(string)

	opts := &tanzuclient.NamespaceOpts{
		Description:       d.Get("description").(string),
		Labels:            d.Get("labels").(map[string]interface{}),
		ManagementCluster: d.Get("management_cluster").(string),
		ProvisionerName:   d.Get("provisioner_name").(string),
		ClusterName:       d.Get("cluster_name").(string),
		WorkspaceName:     d.Get("workspace_name").(string),
	}

	Namespace, err := client.CreateNamespace(NamespaceName, *opts)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create namespace",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(Namespace.Meta.UID)

	return resourceTmcNamespaceRead(ctx, d, m)
}

func resourceTmcNamespaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

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

	return diags
}

func resourceTmcNamespaceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	NamespaceName := d.Get("name").(string)
	clusterName := d.Get("cluster_name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	err := client.DeleteNamespace(NamespaceName, clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete namespace",
			Detail:   fmt.Sprintf("Error deleting resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return diags
}
