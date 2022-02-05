terraform {
  required_version = "~> 1.1.4"
  required_providers {
    installer = {
      source  = "registry.terraform.io/shihanng/installer"
      version = "~> 0.0.1"
    }
  }
}

data "installer_apt" "test" {
  name = "dpkg"
}

resource "installer_apt" "test" {
  name = "sl"
}

output "data_test" {
  value = data.installer_apt.test
}

output "resource_test" {
  value = installer_apt.test
}
