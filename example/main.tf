terraform {
  required_version = "~> 1.1.4"
  required_providers {
    setupenv = {
      source  = "registry.terraform.io/shihanng/setupenv"
      version = "~> 0.0.1"
    }
  }
}

data "setupenv_apt" "test" {
  name = "dpkg"
}

output "out" {
  value = data.setupenv_apt.test
}
