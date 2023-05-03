package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/asdf"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

const asdfIDPrefix = "asdf_plugin:"

func asdfID(name, version string) string {
	return asdfIDPrefix + name + ":" + version
}

func fromASDFID(id string) (name, version string) {
	splitted := strings.Split(id, ":")

	return splitted[1], splitted[2]
}

func resourceASDF() *schema.Resource {
	return &schema.Resource{
		Description: "`installer_asdf` manages a specify version of application using " +
			"[asdf](https://asdf-vm.com/).",
		CreateContext: resourceASDFCreate,
		ReadContext:   resourceASDFRead,
		DeleteContext: resourceASDFDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal ID of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "is the name of the plugin. See `installer_asdf_plugin`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"version": {
				Description: "is the version of the plugin that asdf should install.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"path": {
				Description: "is the path of the application installed by asdf after Terraform creates the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"environment": {
				Description: "are the environment variables set during the installation.",
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceASDFCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := data.Get("name").(string)       // nolint:forcetypeassert
	version := data.Get("version").(string) // nolint:forcetypeassert
	environment := data.Get("environment").(map[string]interface{})

	env := getEnv(environment)

	if err := asdf.Install(ctx, name, version, env); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId(asdfID(name, version))

	resourceASDFRead(ctx, data, meta)

	return diag.Diagnostics{}
}

func resourceASDFRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name, version := fromASDFID(data.Id())

	path, err := asdf.FindInstalled(ctx, name, version)
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

	if err := data.Set("version", version); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("path", path); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceASDFDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name, version := fromASDFID(data.Id())

	if err := asdf.Uninstall(ctx, name, version); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}

func getEnv(environment map[string]interface{}) []string {
	env := make([]string, 0, len(environment))

	for key, value := range environment {
		env = append(env, fmt.Sprintf("%s=%v", key, value))
	}

	return env
}
