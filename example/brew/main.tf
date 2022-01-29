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
