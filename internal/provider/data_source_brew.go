package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBrew() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstalled("brew", "list"),
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
