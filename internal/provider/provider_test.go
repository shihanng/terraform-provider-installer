package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/shihanng/terraform-provider-setupenv/internal/provider"
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){ //nolint:deadcode,unused,varcheck,gochecknoglobals
	"scaffolding": func() (*schema.Provider, error) { //nolint:unparam
		return provider.New("dev")(), nil
	},
}

func TestProvider(t *testing.T) {
	t.Parallel()

	if err := provider.New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) { //nolint:deadcode,unused,thelper
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}
