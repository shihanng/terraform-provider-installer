//go:build brew

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDataSourceBrew(t *testing.T) {
	t.Parallel()

	t.Run("data.setupenv_brew", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccDataSourceBrew,
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.setupenv_brew.test", "name", "sl"),
						resource.TestMatchResourceAttr("data.setupenv_brew.test", "path", regexp.MustCompile(`[\w\./]+bin/sl$`)),
					),
				},
			},
		})
	})

	t.Run("data.setupenv_brew error", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config:      testAccDataSourceBrewError,
					ExpectError: regexp.MustCompile("No such keg:"),
				},
			},
		})
	})
}

const testAccDataSourceBrew = `
data "setupenv_brew" "test" {
  name = "sl"
}
`

const testAccDataSourceBrewError = `
data "setupenv_brew" "test" {
  name = "ls"
}
`
