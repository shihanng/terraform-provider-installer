//go:build brew

// nolint:dupl
package provider_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/shihanng/terraform-provider-installer/internal/brew"
)

func TestAccResourceBrewBasic(t *testing.T) { // nolint:tparallel
	t.Parallel()

	t.Run("resource.installer_brew", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckBrewDestroy,
			Steps: []resource.TestStep{
				{
					Config: readTFFile("./testdata/resources/brew/resources_brew_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_brew.basic"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_brew_tap", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckBrewDestroy,
			Steps: []resource.TestStep{
				{
					Config: readTFFile("./testdata/resources/brew/resources_brew_tap.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_brew.tap"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_brew_cask", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckBrewDestroy,
			Steps: []resource.TestStep{
				{
					Config: readTFFile("./testdata/resources/brew/resources_brew_cask.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_brew.cask"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_brew_cask_fqn", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckBrewDestroy,
			Steps: []resource.TestStep{
				{
					Config: readTFFile("./testdata/resources/brew/resources_brew_cask_fqn.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_brew.cask_fqn"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_brew_cask_treat_as_formula", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckBrewDestroy,
			Steps: []resource.TestStep{
				{
					Config: readTFFile("./testdata/resources/brew/resources_brew_treat_as_formula.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_brew.treat_as_formula"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_brew error", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckBrewDestroy,
			Steps: []resource.TestStep{
				{
					Config:      testAccResourceBrewBasicError,
					ExpectError: regexp.MustCompile("formula not found"),
				},
			},
		})
	})
}

func testAccCheckBrewDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "installer_brew" {
			continue
		}

		path := resource.Primary.Attributes["path"]
		name := resource.Primary.Attributes["name"]

		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			uninstallErr := brew.Uninstall(context.Background(), name)

			return errors.CombineErrors(err, uninstallErr)
		}
	}

	return nil
}

const testAccResourceBrewBasicError = `
resource "installer_brew" "test" {
  name = "abc"
}
`

func readTFFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}
