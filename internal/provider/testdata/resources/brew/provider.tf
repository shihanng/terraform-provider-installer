terraform {
  required_version = ">= 1.0"
  required_providers {
    installer = {
      source  = "registry.terraform.io/shihanng/installer"
      version = "~> 0.6.0"
    }
  }
}
