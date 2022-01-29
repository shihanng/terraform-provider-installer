# terraform-provider-installer

**terraform-provider-installer** is a [Terraform](https://www.terraform.io/) provider for installing softwares via various package management tools. Currently, **terraform-provider-installer** supports

- [APT](https://ubuntu.com/server/docs/package-management)
- [Homebrew](https://brew.sh/)

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

1. Run `make install`. This command also installs the provider in `~/.terraform.d/plugins/registry.terraform.io/shihanng/installer/<version>/<os_arch>/terraform-provider-installer`. On macOS, use

   ```
   OS_ARCH=darwin_arm64 make install
   ```

2. Have a look at [./example](./example) for an example of Terraform configuration. You can also use the example for testing, e.g.
   ```
   terraform -chdir="./example" init
   ```

#### Tips

1. Use `export TF_LOG_PROVIDER=DEBUG` for debugging. See <https://www.terraform.io/internals/debugging>.

## References

- [Custom Providers](https://learn.hashicorp.com/collections/terraform/providers)
- [Terraform Provider Scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding)
- [Terraform Provider Hashicups](https://github.com/hashicorp/terraform-provider-hashicups)
- [Terraform Provider GitHub](https://github.com/integrations/terraform-provider-github)
