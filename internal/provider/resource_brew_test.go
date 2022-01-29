//go:build brew

package provider_test

import (
	"context"
	"fmt"
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
					Config: testAccResourceBrewBasic,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckBrewExists("installer_brew.test"),
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
					ExpectError: regexp.MustCompile("No available formula with the name"),
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

const testAccResourceBrewBasic = `
resource "installer_brew" "test" {
  name = "cowsay"
}
`

const testAccResourceBrewBasicError = `
resource "installer_brew" "test" {
  name = "abc"
}
`

func testAccCheckBrewExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("%s: %w", name, errResourceNotFound)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource '%s': %w", name, errIDNotSet)
		}

		return nil
	}
}
