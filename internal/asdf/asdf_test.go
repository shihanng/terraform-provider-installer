package asdf_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/asdf"
	"github.com/shihanng/terraform-provider-installer/internal/xtests"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestASDF(t *testing.T) { //nolint:tparallel
	t.Parallel()

	reset, err := xtests.SetupASDFDataDir(t.TempDir())
	assert.NilError(t, err)

	t.Cleanup(reset)

	ctx := context.Background()

	name := "terraform-ls"
	gitURL := "https://github.com/shihanng/asdf-terraform-ls"
	version := "0.25.2"

	t.Run("run asdf plugin add", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, asdf.AddPlugin(ctx, name, gitURL, []string{}))
	})

	t.Run("run FindAddedPlugin", func(t *testing.T) { //nolint:paralleltest
		url, err := asdf.FindAddedPlugin(ctx, name)
		assert.NilError(t, err)
		assert.Equal(t, url, gitURL)
	})

	t.Run("run asdf install", func(t *testing.T) { //nolint:paralleltest
		err := asdf.Install(ctx, name, version, []string{})
		t.Log(errors.FlattenDetails(err))
		assert.NilError(t, err)
	})

	t.Run("run asdf where", func(t *testing.T) { //nolint:paralleltest
		path, err := asdf.FindInstalled(ctx, name, version)
		assert.NilError(t, err)
		assert.Assert(t, hasSuffix(path, "installs/terraform-ls/0.25.2"))
	})

	t.Run("run asdf plugin remove", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, asdf.RemovePlugin(ctx, name))
	})
}

func hasSuffix(actual, suffix string) cmp.Comparison {
	return func() cmp.Result {
		if ok := strings.HasSuffix(actual, suffix); ok {
			return cmp.ResultSuccess
		}

		return cmp.ResultFailure(fmt.Sprintf("'%s' does not contain suffix '%s'", actual, suffix))
	}
}
