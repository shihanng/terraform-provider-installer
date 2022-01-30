//go:build apt

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func TestAccDataSourceApt(t *testing.T) {
	t.Parallel()

	t.Run("data.installer_apt", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceApt,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.installer_apt.test", "name", "sl"),
						resource.TestCheckResourceAttr("data.installer_apt.test", "path", "/usr/games/sl"),
					),
				},
			},
		})
	})

	t.Run("data.installer_apt error", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccDataSourceAptError,
					ExpectError: regexp.MustCompile(xerrors.ErrNotInstalled.Error()),
				},
			},
		})
	})
}

const testAccDataSourceApt = `
data "installer_apt" "test" {
  name = "sl"
}
`

const testAccDataSourceAptError = `
data "installer_apt" "test" {
  name = "ls"
}
`
