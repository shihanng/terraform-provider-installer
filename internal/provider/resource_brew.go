package provider //nolint:dupl

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/brew"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
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
		Description: "`installer_brew` manages an application using [Homebrew](https://brew.sh/).\n\n" +
			"It works on systems that use Homebrew as the package management system. " +
			"Adding an `installer_brew` resource means that Terraform will ensure that " +
			"the application defined in the `name` argument is made available via brew.",
		CreateContext: resourceBrewCreate,
		ReadContext:   resourceBrewRead,
		DeleteContext: resourceBrewDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal ID of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the application that `brew` recognizes.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"path": {
				Description: "The path where the application is installed by `brew` after Terraform creates this resource.",
				Type:        schema.TypeString,
				Computed:    true,
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
		return xerrors.ToDiags(err)
	}

	data.SetId(brewID(name))

	resourceBrewRead(ctx, data, meta)

	return diag.Diagnostics{}
}

func resourceBrewRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromBrewID(data.Id())

	path, err := brew.FindInstalled(ctx, name)
	if err != nil {
		if errors.Is(err, xerrors.ErrNotInstalled) {
			data.SetId("")

			return diags
		}

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

func resourceBrewDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromBrewID(data.Id())

	if err := brew.Uninstall(ctx, name); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}
