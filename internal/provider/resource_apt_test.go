//go:build apt

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
	"github.com/shihanng/terraform-provider-installer/internal/apt"
)

func TestAccResourceAptBasic(t *testing.T) { // nolint:tparallel
	t.Parallel()

	t.Run("resource.installer_apt", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckAptDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceAptBasic,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_apt.test"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_apt error", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckAptDestroy,
			Steps: []resource.TestStep{
				{
					Config:      testAccResourceAptBasicError,
					ExpectError: regexp.MustCompile("Unable to locate package"),
				},
			},
		})
	})
}

func testAccCheckAptDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "installer_apt" {
			continue
		}

		path := resource.Primary.Attributes["path"]
		name := resource.Primary.Attributes["name"]

		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			uninstallErr := apt.Uninstall(context.Background(), name)

			return errors.CombineErrors(err, uninstallErr)
		}
	}

	return nil
}

const testAccResourceAptBasic = `
resource "installer_apt" "test" {
  name = "cowsay"
}
`

const testAccResourceAptBasicError = `
resource "installer_apt" "test" {
  name = "abc"
}
`
