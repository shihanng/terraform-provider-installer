package provider_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/shihanng/terraform-provider-installer/internal/asdf"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
	"github.com/shihanng/terraform-provider-installer/internal/xtests"
	"gotest.tools/v3/assert"
)

func TestAccResourceASDFBasic(t *testing.T) { // nolint:dupl,paralleltest
	// Cannot run parallel test because it involves setting environmen variables
	reset, err := xtests.SetupASDFDataDir(t.TempDir())
	assert.NilError(t, err)

	t.Cleanup(reset)

	t.Run("resource.installer_asdf", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckASDFDestroy,
			Steps: []resource.TestStep{
				{
					Config: mustReadFile("../../examples/resources/installer_asdf/resource.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_asdf.this"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_asdf error", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckASDFDestroy,
			Steps: []resource.TestStep{
				{
					Config:      testAccResourceASDFBasicError,
					ExpectError: regexp.MustCompile("No such plugin: (.+)"),
				},
			},
		})
	})
}

func testAccCheckASDFDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "installer_asdf" {
			continue
		}

		name := resource.Primary.Attributes["name"]
		version := resource.Primary.Attributes["version"]
		ctx := context.Background()

		if _, err := asdf.FindInstalled(ctx, name, version); !errors.Is(err, xerrors.ErrNotInstalled) {
			removeErr := asdf.RemovePlugin(ctx, name)

			return errors.CombineErrors(err, removeErr)
		}
	}

	return nil
}

const testAccResourceASDFBasicError = `
resource "installer_asdf" "test" {
  name    = "abc"
  version = "v0.1.2"
}
`
