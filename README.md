[![Checks](https://github.com/shihanng/terraform-provider-installer/actions/workflows/checks.yml/badge.svg)](https://github.com/shihanng/terraform-provider-installer/actions/workflows/checks.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/shihanng/terraform-provider-installer)
[![GitHub license](https://img.shields.io/github/license/shihanng/terraform-provider-installer)](https://github.com/shihanng/terraform-provider-installer/blob/trunk/LICENSE)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit&logoColor=white)](https://github.com/pre-commit/pre-commit)

# terraform-provider-installer

**terraform-provider-installer** is a [Terraform](https://www.terraform.io/) provider for installing softwares via various package management tools. Currently, **terraform-provider-installer** supports

- [APT](https://ubuntu.com/server/docs/package-management)
- [Homebrew](https://brew.sh/)
- Shell script
- [asdf](https://asdf-vm.com/)

The following shows how to install **git** and **starship** through Homebrew using **terraform-provider-installer** provider. See <https://registry.terraform.io/providers/shihanng/installer/latest/docs> for complete documentation.

```tf
terraform {
  required_version = "~> 1.1.4"
  required_providers {
    installer = {
      source  = "shihanng/installer"
    }
  }
}

locals {
  apps = ["git", "starship"]
}

resource "installer_brew" "this" {
  for_each = toset(local.apps)
  name     = each.key
}
```

## Development

### Code quality

We use [pre-commit](https://pre-commit.com/) to maintain the code quality of this project. Refer to [.pre-commit-config.yaml](./.pre-commit-config.yaml) for the list of linters that we are using. Refer to [this page](https://pre-commit.com/#install) to install pre-commit and the git hook script.

```
pre-commit install
```

### Running automated tests

Run unit tests (no resources will be created/destroy) with

```
make test
```

Run acceptance tests with

```
make TESTARGS="-tags=apt" testacc
```

We must provide the value for `-tags` because some tests only runs on specific platform. Currently the valid values for `-tags` are:

- `apt` for environment that uses [APT](https://ubuntu.com/server/docs/package-management).
- `brew` for environment that uses [Homebrew](https://brew.sh/).

### Testing with development version

You added a new feature or fixed a bug in **terraform-provider-installer**. Now you want to test it directly with your Terraform configurations in your local machine. Here is what you can do.

1. Run `make install`. This command installs the provider in `/tmp/tfproviders`. We've setup Terraform in `.terraformrc` to use provider from `/tmp/tfproviders`. On macOS, use

   ```
   OS_ARCH=darwin_arm64 make install
   ```

2. Have a look at [./examples](./examples) for an example of Terraform configuration. You can also use the example for testing, e.g.
   ```
   export TF_CLI_CONFIG_FILE=$(pwd)/.terraformrc
   terraform -chdir="./examples" init
   ```

#### Tips

1. Use `export TF_LOG_PROVIDER=DEBUG` for debugging. See <https://www.terraform.io/internals/debugging>.

## References

- [Custom Providers](https://learn.hashicorp.com/collections/terraform/providers)
- [Terraform Provider Scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding)
- [Terraform Provider Hashicups](https://github.com/hashicorp/terraform-provider-hashicups)
- [Terraform Provider GitHub](https://github.com/integrations/terraform-provider-github)
