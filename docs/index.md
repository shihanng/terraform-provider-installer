---
layout: ""
page_title: "Provider: Installer"
description: |-
  Install/uninstall apps in the local machine using the Installer provider.
---

# Installer Provider

The Installer provider provides resources to manage apps (install/uninstall)
in your local machine in a declarative manner. It currently supports systems that use

- [APT](https://ubuntu.com/server/docs/package-management)
- [Homebrew](https://brew.sh/)

It also supports shell script via the `installer_script` resource.

## Example Usage

There is no configuration on the provider level.
The following shows how to ensure the system has git and starship installed via Homebrew.

```terraform
provider "installer" {}

locals {
  apps = ["tmux"]
}

resource "installer_brew" "this" {
  for_each = toset(local.apps)
  name     = each.key
}
```
