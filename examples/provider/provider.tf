provider "installer" {}

locals {
  apps = ["git", "starship"]
}

resource "installer_brew" "this" {
  for_each = toset(local.apps)
  name     = each.key
}
