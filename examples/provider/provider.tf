provider "installer" {}

locals {
  apps = ["tmux"]
}

resource "installer_brew" "this" {
  for_each = toset(local.apps)
  name     = each.key
}
