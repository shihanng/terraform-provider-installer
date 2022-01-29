// nolint:dupl
package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/apt"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

const aptIDPrefix = "apt:"

func aptID(name string) string {
	return aptIDPrefix + name
}

func nameFromAptID(id string) string {
	return strings.TrimPrefix(id, aptIDPrefix)
}

func resourceApt() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAptCreate,
		ReadContext:   resourceAptRead,
		DeleteContext: resourceAptDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAptCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := data.Get("name").(string) // nolint:forcetypeassert

	if err := apt.Install(ctx, name); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId(aptID(name))

	return resourceAptRead(ctx, data, meta)
}

func resourceAptRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromAptID(data.Id())

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

	return diags
}

func resourceAptDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromAptID(data.Id())

	if err := apt.Uninstall(ctx, name); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}
