# terraform-provider-setupenv

**terraform-provider-setupenv** is a [Terraform](https://www.terraform.io/) provider to set up a development environment machine.

## Development

### Code quality

We use [pre-commit](https://pre-commit.com/) to maintain the code quality of this project. Refer to [.pre-commit-config.yaml](./.pre-commit-config.yaml) for the list of linters that we are using. Refer to [this page](https://pre-commit.com/#install) to install pre-commit and the git hook script.

```
pre-commit install
```

### Testing with development version

You added a new feature or fixed a bug in **terraform-provider-setupenv**. Now you want to test it directly with your Terraform configurations in your local machine. Here is what you can do.

1. Run `make install`. This command also installs the provider in `~/.terraform.d/plugins/registry.terraform.io/shihanng/setupenv/<version>/<os_arch>/terraform-provider-setupenv`. On macOS, use

   ```
   OS_ARCH=darwin_arm64 make install
   ```

2. Have a look at [./example](./example) for an example of Terraform configuration. You can also use the example for testing, e.g.
   ```
   terraform -chdir="./example" init
   ```

## References

- [Custom Providers](https://learn.hashicorp.com/collections/terraform/providers)
- [Terraform Provider Scaffolding](https://github.com/hashicorp/terraform-provider-scaffolding)
- [Terraform Provider Hashicups](https://github.com/hashicorp/terraform-provider-hashicups)
