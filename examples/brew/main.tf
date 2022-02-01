terraform {
  required_version = "~> 1.1.4"
  required_providers {
    installer = {
      source  = "registry.terraform.io/shihanng/installer"
      version = "~> 0.0.1"
    }
  }
}

resource "installer_brew" "test" {
  name = "sl"
}

output "resource_test" {
  value = installer_brew.test
}

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
