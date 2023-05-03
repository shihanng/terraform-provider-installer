package provider

import (
	"context"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/asdf"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

const asdfPluginIDPrefix = "asdf_plugin:"

func asdfPluginID(name string) string {
	return asdfPluginIDPrefix + name
}

func nameFromASDFPluginID(id string) string {
	return strings.TrimPrefix(id, asdfPluginIDPrefix)
}

func resourceASDFPlugin() *schema.Resource {
	return &schema.Resource{
		Description:   "`installer_asdf_plugin` manages an [asdf plugin](https://asdf-vm.com/manage/plugins.html).",
		CreateContext: resourceASDFPluginCreate,
		ReadContext:   resourceASDFPluginRead,
		DeleteContext: resourceASDFPluginDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal ID of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "is the name of the plugin.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"git_url": {
				Description: "is the Git repository's URL which will be added as the plugin if specified.",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
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

func resourceASDFPluginCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := data.Get("name").(string)      // nolint:forcetypeassert
	gitURL := data.Get("git_url").(string) // nolint:forcetypeassert
	environment := data.Get("environment").(map[string]interface{})

	env := getEnv(environment)

	if err := asdf.AddPlugin(ctx, name, gitURL, env); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId(asdfPluginID(name))

	resourceASDFPluginRead(ctx, data, meta)

	return diag.Diagnostics{}
}

func resourceASDFPluginRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromASDFPluginID(data.Id())

	gitURL, err := asdf.FindAddedPlugin(ctx, name)
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

	if err := data.Set("git_url", gitURL); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceASDFPluginDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromASDFPluginID(data.Id())

	if err := asdf.RemovePlugin(ctx, name); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}
