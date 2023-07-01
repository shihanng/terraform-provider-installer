terraform {
  required_version = "~> 1.1.4"
  required_providers {
    installer = {
      source  = "shihanng/installer"
      version = "~> 0.0.1"
    }
  }
}

provider "installer" {
}

locals {
  apps = ["git", "starship"]
}

resource "installer_brew" "this" {
  for_each = toset(local.apps)
  name     = each.key
}
