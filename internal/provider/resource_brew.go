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
				Description: "Name of the application that `brew` recognizes, e.g., `homebrew/cask/alfred` for a cask, `goreleaser/tap/goreleaser` for tap. Treats a package as a formula if `cask` is not set or false",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"cask": {
				Description: "Treat name argument as cask.",
				Type:        schema.TypeBool,
				Required:    false,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
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
	cask := data.Get("cask").(bool)   // nolint:forcetypeassert

	if err := brew.Install(ctx,
		brew.NewCmd("install", name, brew.WithCask(cask)).Args); err != nil {
		return xerrors.ToDiags(err)
	}

	info, err := brew.GetInfo(ctx,
		brew.NewCmd("info", name, brew.WithCask(cask), brew.WithJSONV2()).Args)
	if err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId(brewID(info.Name))

	resourceBrewRead(ctx, data, meta)

	return diag.Diagnostics{}
}

func resourceBrewRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromBrewID(data.Id())

	info, err := brew.GetInfo(ctx,
		brew.NewCmd("info", name, brew.WithJSONV2()).Args)
	if err != nil {
		return xerrors.ToDiags(err)
	}

	var (
		path    string
		pathErr error
	)

	if info.IsCask {
		path, pathErr = brew.FindCaskPath(ctx,
			brew.NewCmd("list", name, brew.WithCask(info.IsCask)).Args)
	} else {
		path, pathErr = brew.FindInstalled(ctx, info.Name)
	}

	if pathErr != nil {
		if errors.Is(pathErr, xerrors.ErrNotInstalled) {
			data.SetId("")

			return diags
		}

		return xerrors.ToDiags(pathErr)
	}

	if err := data.Set("path", path); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("cask", info.IsCask); err != nil {
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
