package apt

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
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
		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	paths := strings.Split(string(out), "\n")

	return findExecutablePath(paths)
}

func Uninstall(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "sudo", "apt-get", "-y", "remove", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

func findExecutablePath(paths []string) (string, error) {
	for _, path := range paths {
		info, err := os.Lstat(path)
		if err != nil {
			continue
		}

		// If executable by either owner, group, or other
		if !info.IsDir() && info.Mode()&0o111 != 0 {
			return path, nil
		}
	}

	return "", xerrors.ErrNotExecutable
}
