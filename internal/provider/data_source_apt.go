// nolint:dupl
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/apt"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func dataSourceApt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAptRead,
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

func dataSourceAptRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := data.Get("name").(string) // nolint:forcetypeassert

	path, err := apt.FindInstalled(ctx, name)
	if err != nil {
		return xerrors.ToDiags(err)
	}

	if err := data.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("path", path); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(aptID(name))

	return diags
}
