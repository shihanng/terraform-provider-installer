package provider

import (
	"bytes"
	"context"
	"errors"
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

func dataSourceAptRead(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	name := data.Get("name").(string) // nolint:forcetypeassert

	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, "dpkg", "-L", name)
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

func findPathHasSuffix(paths []string, suffix string) (string, error) {
	for _, p := range paths {
		ok := strings.HasSuffix(p, suffix)
		if ok {
			return p, nil
		}
	}

	return "", errNotFound
}
