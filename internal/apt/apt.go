package apt

import (
	"context"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func Install(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "sudo", "apt-get", "-y", "install", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

func FindInstalled(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "dpkg", "-L", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "is not installed") {
			return "", xerrors.ErrNotInstalled
		}

		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	paths := strings.Split(string(out), "\n")

	return system.FindExecutablePath(paths) // nolint:wrapcheck
}

func Uninstall(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "sudo", "apt-get", "-y", "remove", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}
