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
	info := GetInfo(name)

	cmd := exec.CommandContext(ctx, "dpkg", "-L", info.Name) //nolint:gosec

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "is not installed") {
			return "", xerrors.ErrNotInstalled
		}

		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	if info.Version != "" {
		cmd := exec.CommandContext(ctx, "dpkg", "-s", info.Name) //nolint:gosec

		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
		}

		installedVersion := ExtractVersion(string(out))
		if info.Version != installedVersion {
			return "", xerrors.ErrNotInstalled
		}
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

// Info contains the package name and the package version.
type Info struct {
	Name    string
	Version string
}

const nPart = 2

// GetInfo splits the name and version from string name=version and put the
// values into Info.
func GetInfo(original string) Info {
	var info Info

	splitted := strings.SplitN(original, "=", nPart)

	info.Name = splitted[0]

	if len(splitted) == nPart {
		info.Version = splitted[1]
	}

	return info
}

// ExtractVersion extracts version value from the output of dpkg -s <package>.
func ExtractVersion(input string) string {
	for _, line := range strings.Split(input, "\n") {
		if strings.HasPrefix(line, "Version: ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Version: "))
		}
	}

	return ""
}
