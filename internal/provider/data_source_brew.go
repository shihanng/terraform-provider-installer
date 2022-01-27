package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBrew() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBrewRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBrewRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics { //nolint:dupl
	var diags diag.Diagnostics

	name := data.Get("name").(string) // nolint:forcetypeassert

	cmd := exec.CommandContext(ctx, "brew", "list", name)

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

	validPath, err := findPathHasSuffix(paths, "bin/"+name)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("name", name); err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("path", validPath); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(validPath)

	return diags
}
