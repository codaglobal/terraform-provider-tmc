package tmc

import (
	"context"
	"fmt"

	"github.com/codaglobal/terraform-provider-tmc/tanzuclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTmcAwsDataProtectionCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTmcAwsDataProtectionCredentialRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID of the Tanzu Aws Data Protection Credential",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Tanzu Aws Data Protection Credential",
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Spec defines the specification of the desired Account Credential",
				Elem: &schema.Resource{
					Schema: awsDataCredentialSpecFields(false, true),
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTmcAwsDataProtectionCredentialRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*tanzuclient.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	credential, err := client.GetAwsDataProtectionCredential(d.Get("name").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("spec", credential.Spec); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Failed to read credential specs",
			Detail:   fmt.Sprintf("Error fetching spec for resource %s: %s", d.Get("name"), err),
		})
		return diags
	}
	d.SetId(string(credential.MetaData.UID))
	d.Set("status", credential.Status.Phase)

	return diags
}
