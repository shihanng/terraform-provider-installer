package brew_test

import (
	"context"
	"strings"
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
	nameCask := "tiny-player"

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

	t.Run("check brew cask install", func(t *testing.T) { //nolint:paralleltest
		err := brew.Install(ctx, brew.NewCmd("install", nameCask, brew.WithCask(true)).Args)
		t.Log(errors.FlattenDetails(err))
		assert.NilError(t, err)
	})

	t.Run("brew install cask must fail", func(t *testing.T) { //nolint:paralleltest
		err := brew.Install(ctx, brew.NewCmd("install", name, brew.WithCask(true)).Args)
		t.Log(errors.FlattenDetails(err))
		assert.Assert(t, err != nil)
	})

	t.Run("check info", func(t *testing.T) { //nolint:paralleltest
		want := brew.Info{
			Name:   "homebrew/core/cowsay",
			IsCask: false,
		}

		got, err := brew.GetInfo(ctx, brew.NewCmd("info", name, brew.WithCask(false), brew.WithJSONV2()).Args)
		t.Log(errors.FlattenDetails(err))
		assert.NilError(t, err)
		assert.DeepEqual(t, want, got)
	})

	t.Run("check info", func(t *testing.T) { //nolint:paralleltest
		want := brew.Info{
			Name:   "homebrew/core/cowsay",
			IsCask: false,
		}

		got, err := brew.GetInfo(ctx, brew.NewCmd("info", name, brew.WithCask(false), brew.WithJSONV2()).Args)
		t.Log(errors.FlattenDetails(err))
		assert.NilError(t, err)
		assert.DeepEqual(t, want, got)
	})

	t.Run("run brew list", func(t *testing.T) { //nolint:paralleltest
		path, err := brew.FindInstalled(ctx, name)
		assert.NilError(t, err)
		ok := strings.HasSuffix(path, "bin/cowsay")
		assert.Equal(t, ok, true)
	})

	t.Run("run brew list cask", func(t *testing.T) { //nolint:paralleltest
    path, err := brew.FindCaskPath(ctx, brew.NewCmd("list", nameCask, brew.WithCask(true)).Args)
		assert.NilError(t, err)
		ok := strings.HasSuffix(path, "Tiny Player.app")
		assert.Equal(t, ok, true)
	})

	t.Run("run brew uninstall", func(t *testing.T) { //nolint:paralleltest
		assert.NilError(t, brew.Uninstall(ctx, name))
	})
}
