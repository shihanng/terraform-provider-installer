//go:build apt

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceApt(t *testing.T) {
	t.Parallel()

	t.Run("data.setupenv_apt", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceApt,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.setupenv_apt.test", "name", "dpkg"),
						resource.TestCheckResourceAttr("data.setupenv_apt.test", "path", "/usr/bin/dpkg"),
					),
				},
			},
		})
	})

	t.Run("data.setupenv_apt error", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccDataSourceAptError,
					ExpectError: regexp.MustCompile("is not installed"),
				},
			},
		})
	})
}

const testAccDataSourceApt = `
data "setupenv_apt" "test" {
  name = "dpkg"
}
`

const testAccDataSourceAptError = `
data "setupenv_apt" "test" {
  name = "ls"
}
`
