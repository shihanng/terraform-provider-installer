package script_test

import (
	"context"
	"testing"

	"github.com/shihanng/terraform-provider-installer/internal/script"
	"gotest.tools/v3/assert"
)

func TestScript(t *testing.T) { //nolint:tparallel
	t.Parallel()

	ctx := context.Background()

	t.Run("run install script", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, script.Run(ctx, testInstallScript))
	})

	t.Run("check installed", func(t *testing.T) { //nolint:paralleltest
		ok, err := script.IsInstalled("/tmp/installer-myapp")
		assert.NilError(t, err)
		assert.Equal(t, ok, true)
	})

	t.Run("run uninstall script", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, script.Run(ctx, testUninstallScript))
	})
}

const testInstallScript = `
/bin/bash

touch /tmp/installer-myapp
chmod +x /tmp/installer-myapp
exit 0
`

const testUninstallScript = `
/bin/bash

rm -f /tmp/installer-myapp
exit 0
`
