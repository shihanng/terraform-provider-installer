package asdf

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

func AddPlugin(ctx context.Context, name, gitURL string, env []string) error {
	cmd := exec.CommandContext(ctx, "asdf", "plugin", "add", name, gitURL)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

func FindAddedPlugin(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "asdf", "plugin", "list", "--urls", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	entries := strings.Split(string(out), "\n")

	return findPlugin(name, entries)
}

func findPlugin(name string, entries []string) (string, error) {
	for _, entry := range entries {
		fields := strings.Fields(entry)

		if len(fields) != 2 { //nolint:gomnd
			// Ignore invalid format
			continue
		}

		if fields[0] == name {
			return fields[1], nil
		}
	}

	return "", xerrors.ErrNotInstalled
}

func RemovePlugin(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "asdf", "plugin", "remove", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

func Install(ctx context.Context, name, version string, env []string) error {
	cmd := exec.CommandContext(ctx, "asdf", "install", name, version)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, env...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

func FindInstalled(ctx context.Context, name, version string) (string, error) {
	cmd := exec.CommandContext(ctx, "asdf", "where", name, version)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "No such plugin:") {
			return "", xerrors.ErrNotInstalled
		}

		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return string(bytes.TrimSpace(out)), nil
}

func Uninstall(ctx context.Context, name, version string) error {
	cmd := exec.CommandContext(ctx, "asdf", "uninstall", name, version)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}
