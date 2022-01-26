package provider

import (
	"bytes"
	"context"
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

	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, "brew", "list", name)
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return diag.FromErr(err)
	}

	paths := strings.Split(out.String(), "\n")

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
