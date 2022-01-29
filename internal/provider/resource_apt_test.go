//go:build apt

package provider_test

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
			cmd := exec.Command("sudo", "apt-get", "-y", "remove", name)

			if out, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("%s: %w", string(out), err)
			}

			return fmt.Errorf("unexpect error from stat: %w", err)
		}
	}

	return nil
}

const testAccResourceAptBasic = `
resource "installer_apt" "test" {
  name = "sl"
}
`

const testAccResourceAptBasicError = `
resource "installer_apt" "test" {
  name = "abc"
}
`
