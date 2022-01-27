package provider

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var errNotFound = errors.New("could not find path with suffix")

func dataSourceApt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAptRead,
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

func dataSourceAptRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics { //nolint:dupl
	var diags diag.Diagnostics

	name := data.Get("name").(string) // nolint:forcetypeassert

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

func findPathHasSuffix(paths []string, suffix string) (string, error) {
	for _, p := range paths {
		ok := strings.HasSuffix(p, suffix)
		if ok {
			return p, nil
		}
	}

	return "", errNotFound
}
