package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTmcScan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTmcScanCreate,
		ReadContext:   resourceTmcScanRead,
		DeleteContext: resourceTmcScanDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID for the Cluster Scan",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the scan",
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Cluster where the scans are to be performed",
			},
			"management_cluster": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an existing management cluster to be used",
			},
			"provisioner_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of an existing provisioner to be used",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Type of the scan to perform",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "lite" && v != "cis" && v != "conformance" {
						errs = append(errs, fmt.Errorf("supported scan types are lite, cis and conformance"))
					}
					return
				},
			},
		},
	}
}

func resourceTmcScanCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	clusterName := d.Get("cluster_name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)
	scanType := d.Get("type").(string)

	scan, err := client.CreateScan(clusterName, managementClusterName, provisionerName, scanType)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to create scan for the cluster",
			Detail:   fmt.Sprintf("Error creating resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId(scan.FullName.Name)
	d.Set("name", scan.FullName.Name)
	d.Set("type", scanType)

	resourceTmcClusterGroupRead(ctx, d, m)

	return nil
}
func resourceTmcScanRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	scanName := d.Get("name").(string)
	clusterName := d.Get("cluster_name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	scan, err := client.GetScan(scanName, clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read scan for the cluster",
			Detail:   fmt.Sprintf("Error reading resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.Set("name", scan.FullName.Name)

	return nil
}

func resourceTmcScanDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := m.(*tanzuclient.Client)

	scanName := d.Get("name").(string)
	clusterName := d.Get("cluster_name").(string)
	managementClusterName := d.Get("management_cluster").(string)
	provisionerName := d.Get("provisioner_name").(string)

	err := client.DeleteScan(scanName, clusterName, managementClusterName, provisionerName)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to delete scan for the cluster",
			Detail:   fmt.Sprintf("Error deleting resource %s: %s", d.Get("name"), err),
		})
		return diags
	}

	d.SetId("")

	return nil
}
