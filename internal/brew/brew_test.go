package brew_test

import (
	"context"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/brew"
	"gotest.tools/v3/assert"
)

func TestBrew(t *testing.T) { //nolint:tparallel
	t.Parallel()

	// reset, err := xtests.SetupASDFDataDir(t.TempDir())
	// assert.NilError(t, err)

	// t.Cleanup(reset)

	ctx := context.Background()

	name := "cowsay"

	t.Run("check brew cmd", func(t *testing.T) { //nolint:paralleltest
		want := &brew.Cmd{
			Args: []string{
				"info",
				"cowsay",
				"--formulae",
				"--json=v2",
			},
		}

		got := brew.NewCmd("info", name, brew.WithCask(false), brew.WithJSONV2())
		assert.DeepEqual(t, want, got)
	})

	t.Run("check brew cmd with cask", func(t *testing.T) { //nolint:paralleltest
		want := &brew.Cmd{
			Args: []string{
				"info",
				"cowsay",
				"--cask",
				"--json=v2",
			},
		}

		got := brew.NewCmd("info", name, brew.WithCask(true), brew.WithJSONV2())
		assert.DeepEqual(t, want, got)
	})

	t.Run("check brew install", func(t *testing.T) { //nolint:paralleltest
		err := brew.Install(ctx, brew.NewCmd("install", name, brew.WithCask(false)).Args)
		t.Log(errors.FlattenDetails(err))
		assert.NilError(t, err)
	})

	t.Run("brew install cask must fail", func(t *testing.T) { //nolint:paralleltest
		err := brew.Install(ctx, brew.NewCmd("install", name, brew.WithCask(true)).Args)
		t.Log(errors.FlattenDetails(err))
		assert.Assert(t, err != nil)
	})

	// t.Run("run asdf install", func(t *testing.T) { //nolint:paralleltest
	// 	err := asdf.Install(ctx, name, version, []string{})
	// 	t.Log(errors.FlattenDetails(err))
	// 	assert.NilError(t, err)
	// })

	// t.Run("run asdf where", func(t *testing.T) { //nolint:paralleltest
	// 	path, err := asdf.FindInstalled(ctx, name, version)
	// 	assert.NilError(t, err)
	// 	assert.Assert(t, hasSuffix(path, "installs/terraform-ls/0.25.2"))
	// })

	// t.Run("run asdf plugin remove", func(t *testing.T) { //nolint:paralleltest
	// 	assert.NilError(t, asdf.RemovePlugin(ctx, name))
	// })
}
