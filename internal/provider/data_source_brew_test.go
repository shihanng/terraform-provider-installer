//go:build brew

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func TestAccDataSourceBrew(t *testing.T) {
	t.Parallel()

	t.Run("data.installer_brew", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceBrew,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.installer_brew.test", "name", "sl"),
						resource.TestMatchResourceAttr("data.installer_brew.test", "path", regexp.MustCompile(`[\w\./]+bin/sl$`)),
					),
				},
			},
		})
	})

	t.Run("data.installer_brew error", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccDataSourceBrewError,
					ExpectError: regexp.MustCompile(xerrors.ErrNotInstalled.Error()),
				},
			},
		})
	})
}

const testAccDataSourceBrew = `
data "installer_brew" "test" {
  name = "sl"
}
`

const testAccDataSourceBrewError = `
data "installer_brew" "test" {
  name = "ls"
}
`
