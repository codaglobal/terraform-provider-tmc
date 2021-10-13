package tmc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// labelsSchema returns the schema to use for labels.
//
func labelsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return k == "labels.tmc.cloud.vmware.com/creator"
		},
	}
}

func labelsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
	}
}
