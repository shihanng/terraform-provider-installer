package provider_test

import (
	"context"
	"io/ioutil"
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

func TestAccResourceASDFPluginBasic(t *testing.T) { // nolint:dupl,paralleltest
	// Cannot run parallel test because it involves setting environmen variables
	reset, err := xtests.SetupASDFDataDir(t.TempDir())
	assert.NilError(t, err)

	t.Cleanup(reset)

	t.Run("resource.installer_asdf_plugin", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckASDFPluginDestroy,
			Steps: []resource.TestStep{
				{
					Config: mustReadFile("../../examples/resources/installer_asdf_plugin/resource.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_asdf_plugin.this"),
					),
				},
			},
		})
	})

	t.Run("resource.installer_asdf_plugin error", func(t *testing.T) { // nolint:paralleltest // due to locking
		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			CheckDestroy:      testAccCheckASDFPluginDestroy,
			Steps: []resource.TestStep{
				{
					Config:      testAccResourceASDFPluginBasicError,
					ExpectError: regexp.MustCompile("plugin (.+) not found in repository"),
				},
			},
		})
	})
}

func testAccCheckASDFPluginDestroy(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "installer_asdf_plugin" {
			continue
		}

		name := resource.Primary.Attributes["name"]
		ctx := context.Background()

		if _, err := asdf.FindAddedPlugin(ctx, name); !errors.Is(err, xerrors.ErrNotInstalled) {
			removeErr := asdf.RemovePlugin(ctx, name)

			return errors.CombineErrors(err, removeErr)
		}
	}

	return nil
}

const testAccResourceASDFPluginBasicError = `
resource "installer_asdf_plugin" "test" {
  name = "abc"
}
`

func mustReadFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(content)
}
