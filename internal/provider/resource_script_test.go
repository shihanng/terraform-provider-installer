package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceScriptBasic(t *testing.T) {
	t.Parallel()

	t.Run("resource.installer_script", func(t *testing.T) {
		t.Parallel()

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { testAccPreCheck(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccResourceScriptBasic,
					Check: resource.ComposeTestCheckFunc(
						testAccCheckResourceExists("installer_script.test"),
					),
				},
			},
		})
	})
}

const testAccResourceScriptBasic = `
resource "installer_script" "test" {
  path           = "/tmp/installer-myapp-test"
  install_script = <<-EOF
  /bin/bash

  touch /tmp/installer-myapp-test
  chmod +x /tmp/installer-myapp-test
  exit 0
  EOF

  uninstall_script = <<-EOF
  /bin/bash

  rm -f /tmp/installer-myapp-test
  exit 0
  EOF
}
`
