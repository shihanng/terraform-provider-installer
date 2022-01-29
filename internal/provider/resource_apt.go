package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
	var diags diag.Diagnostics

	name := data.Get("name").(string) // nolint:forcetypeassert

	cmd := exec.CommandContext(ctx, "sudo", "apt-get", "-y", "install", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("'%s' returns %v", strings.Join(cmd.Args, " "), err.Error()),
				Detail:   string(out),
			},
		}
	}

	data.SetId(aptID(name))

	resourceAptRead(ctx, data, meta)

	return diags
}

const aptIDPrefix = "apt:"

func aptID(name string) string {
	return aptIDPrefix + name
}

func nameFromAptID(id string) string {
	return strings.TrimPrefix(id, aptIDPrefix)
}

func findPathIsExecutable(ctx context.Context, paths []string) (string, error) {
	tflog.Debug(ctx, "recevied paths", "paths", paths)

	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			tflog.Debug(ctx, "lstat error", "error", err, "path", path)

			continue
		}

		// If executable by either owner, group, or other
		if !info.IsDir() && info.Mode()&0o111 != 0 {
			return path, nil
		}
	}

	return "", errNotFound
}

func resourceAptRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromAptID(data.Id())

	cmd := exec.CommandContext(ctx, "dpkg", "-L", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("'%s' returns %v", strings.Join(cmd.Args, " "), err.Error()),
				Detail:   string(out),
			},
		}
	}

	paths := strings.Split(string(out), "\n")

	validPath, err := findPathIsExecutable(ctx, paths)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("path", validPath); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceAptDelete(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := nameFromAptID(data.Id())

	cmd := exec.CommandContext(ctx, "sudo", "apt-get", "-y", "remove", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("'%s' returns %v", strings.Join(cmd.Args, " "), err.Error()),
				Detail:   string(out),
			},
		}
	}

	data.SetId("")

	return diags
}
