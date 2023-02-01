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
		Description: "`installer_script` manages an application using a custom script.\n\n" +
			"Adding an `installer_script` resource means that Terraform will install " +
			"application in the `path` by running the `install_script` when creating the resource.",
		CreateContext: resourceScriptCreate,
		ReadContext:   resourceScriptRead,
		DeleteContext: resourceScriptDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Internal ID of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"path": {
				Description: "is the location of the application installed by the install script. " +
					"If the application does not exist at path, the resource is considered not exist by Terraform",
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"install_script": {
				Description: "is the script that will be called by Terraform when executing `terraform plan/apply`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"uninstall_script": {
				Description: "is the script that will be called by Terraform when executing `terraform destroy`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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

	tflog.Debug(ctx, "run", map[string]interface{}{
		"install_script":  installScript,
		"uninstallScript": uninstallScript,
	})

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

	tflog.Debug(ctx, "run", map[string]interface{}{
		"uninstallScript": uninstallScript,
	})

	if err := script.Run(ctx, uninstallScript); err != nil {
		return xerrors.ToDiags(err)
	}

	data.SetId("")

	return diags
}
