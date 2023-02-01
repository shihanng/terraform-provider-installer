package brew

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

var ErrFormulaNotFound = errors.New("formula not found")

func Install(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "brew", "install", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	if strings.Contains(string(out), "No available formula with the name") {
		return ErrFormulaNotFound
	}

	return nil
}

func FindInstalled(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "brew", "list", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "No available formula with the name") {
			return "", xerrors.ErrNotInstalled
		}

		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	paths := strings.Split(string(out), "\n")

	return system.FindExecutablePath(paths) // nolint:wrapcheck
}

func Uninstall(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "brew", "uninstall", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}
