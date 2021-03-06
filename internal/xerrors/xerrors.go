package xerrors

import (
	"github.com/cockroachdb/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var ErrNotInstalled = errors.New("not installed")

func ToDiags(err error) diag.Diagnostics {
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Error,
			Summary:  err.Error(),
			Detail:   errors.FlattenDetails(err),
		},
	}
}
