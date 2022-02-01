package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-installer/internal/script"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

const scriptIDPrefix = "script:"

func scriptID(path string) string {
	return scriptIDPrefix + path
}

func pathFromScriptID(id string) string {
	return strings.TrimPrefix(id, scriptIDPrefix)
}

func resourceScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScriptCreate,
		ReadContext:   resourceScriptRead,
		DeleteContext: resourceScriptDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"install_script": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"uninstall_script": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceScriptCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	path := data.Get("path").(string) // nolint:forcetypeassert

	installScript := data.Get("install_script").(string) // nolint:forcetypeassert

	uninstallScript := data.Get("uninstall_script").(string) // nolint:forcetypeassert

	tflog.Debug(ctx, "run", "install_script", installScript)

	if err := script.Run(ctx, installScript); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId(scriptID(path))

	if err := data.Set("install_script", installScript); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("uninstall_script", uninstallScript); err != nil {
		return diag.FromErr(err)
	}

	resourceScriptRead(ctx, data, meta)

	return diag.Diagnostics{}
}

func resourceScriptRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	path := pathFromScriptID(data.Id())

	ok, err := script.IsInstalled(path)
	if err != nil {
		return diag.FromErr(err)
	}

	if !ok {
		data.SetId("")
	}

	return diags
}

func resourceScriptDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	uninstallScript := data.Get("uninstall_script").(string) // nolint:forcetypeassert

	tflog.Debug(ctx, "run", "uninstall_script", uninstallScript)

	if err := script.Run(ctx, uninstallScript); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}
