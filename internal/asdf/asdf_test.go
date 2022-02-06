package asdf_test

import (
	"context"
	"os"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/asdf"
	"gotest.tools/v3/assert"
)

func TestASDF(t *testing.T) { //nolint:tparallel
	t.Parallel()

	reset := envSetter(map[string]string{
		"ASDF_DATA_DIR": "/tmp/tfi_asdf",
	})

	t.Cleanup(reset)

	ctx := context.Background()

	name := "terraform-ls"
	gitURL := "https://github.com/shihanng/asdf-terraform-ls"
	version := "0.25.2"

	t.Run("run asdf plugin add", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, asdf.AddPlugin(ctx, name, gitURL))
	})

	t.Run("run FindAddedPlugin", func(t *testing.T) { //nolint:paralleltest
		url, err := asdf.FindAddedPlugin(ctx, name)
		assert.NilError(t, err)
		assert.Equal(t, url, gitURL)
	})

	t.Run("run asdf install", func(t *testing.T) { //nolint:paralleltest
		err := asdf.Install(ctx, name, version)
		t.Log(errors.FlattenDetails(err))
		assert.NilError(t, err)
	})

	t.Run("run asdf where", func(t *testing.T) { //nolint:paralleltest
		path, err := asdf.FindInstalled(ctx, name, version)
		assert.NilError(t, err)
		assert.Equal(t, path, "/tmp/tfi_asdf/installs/terraform-ls/0.25.2")
	})

	t.Run("run asdf plugin remove", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, asdf.RemovePlugin(ctx, name))
	})
}

func envSetter(envs map[string]string) (closer func()) {
	originals := map[string]string{}

	for name, value := range envs {
		if val, ok := os.LookupEnv(name); ok {
			originals[name] = val
		}

		_ = os.Setenv(name, value)
	}

	return func() {
		for name := range envs {
			original, ok := originals[name]
			if ok {
				_ = os.Setenv(name, original)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}
