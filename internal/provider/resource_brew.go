// nolint:dupl
package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/brew"
)

const brewIDPrefix = "brew:"

func brewID(name string) string {
	return brewIDPrefix + name
}

func nameFromBrewID(id string) string {
	return strings.TrimPrefix(id, brewIDPrefix)
}

func resourceBrew() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBrewCreate,
		ReadContext:   resourceBrewRead,
		DeleteContext: resourceBrewDelete,
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

func resourceBrewCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := data.Get("name").(string) // nolint:forcetypeassert

	if err := brew.Install(ctx, name); err != nil {
		return brew.ToDiags(err)
	}

	data.SetId(brewID(name))

	return resourceBrewRead(ctx, data, meta)
}

func resourceBrewRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromBrewID(data.Id())

	path, err := brew.FindInstalled(ctx, name)
	if err != nil {
		return brew.ToDiags(err)
	}

	if err := data.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("path", path); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceBrewDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromBrewID(data.Id())

	if err := brew.Uninstall(ctx, name); err != nil {
		return brew.ToDiags(err)
	}

	data.SetId("")

	return diags
}
