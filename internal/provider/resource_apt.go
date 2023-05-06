package provider //nolint:dupl

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"
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
		Description: "`installer_apt` manages an application using [APT](https://en.wikipedia.org/wiki/APT_(software)).\n\n" +
			"It works on systems that use APT as the package management system. " +
			"Adding an `installer_apt` resource means that Terraform will ensure that " +
			"the application defined in the `name` argument is made available via APT.",
		CreateContext: resourceAptCreate,
		ReadContext:   resourceAptRead,
		DeleteContext: resourceAptDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal ID of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the application that `apt-get` recognizes. Specify a version of a package by following the package name with an equal sign and the version, e.g., `vim=2:8.2.3995-1ubuntu2.7`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"path": {
				Description: "The path where the application is installed by `apt-get` after Terraform creates this resource.",
				Type:        schema.TypeString,
				Computed:    true,
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

	resourceAptRead(ctx, data, meta)

	return diag.Diagnostics{}
}

func resourceAptRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromAptID(data.Id())

	path, err := apt.FindInstalled(ctx, name)
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

func resourceAptDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromAptID(data.Id())

	if err := apt.Uninstall(ctx, name); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}
