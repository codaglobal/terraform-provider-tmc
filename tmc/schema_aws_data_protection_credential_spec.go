package tmc

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func awsDataCredentialSpecFields(isUpdatable bool, isComputed bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{

		"capability": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers. Value must be a positive integer.",
		},
		"data": {
			Type:     schema.TypeList,
			Required: !isComputed,
			Computed: isComputed,
			ForceNew: isUpdatable,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"aws_credential": {
						Type:     schema.TypeList,
						Required: !isComputed,
						Computed: isComputed,
						ForceNew: isUpdatable,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"iam_role_arn": {
									Type:        schema.TypeString,
									Required:    !isComputed,
									Computed:    isComputed,
									ForceNew:    isUpdatable,
									Description: "AWS IAM Role arn to use for the Tanzu Data Protection. It should have the required permissions.",
								},
							},
						},
					},
				},
			},
		},
		"meta": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"provider": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	return s
}
