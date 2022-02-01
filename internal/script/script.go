package script

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/cockroachdb/errors"
)

const execPerm = 0o700

func Run(ctx context.Context, script string) error {
	tmpfile, err := ioutil.TempFile("", "installer")
	if err != nil {
		return errors.Wrap(err, "create temp file for script")
	}

	defer os.Remove(tmpfile.Name())

	if err := os.Chmod(tmpfile.Name(), execPerm); err != nil {
		return errors.Wrap(err, "change temp file permission to 0700")
	}

	if _, err := tmpfile.WriteString(script); err != nil {
		return errors.Wrap(err, "copy script to temp file")
	}

	if err := tmpfile.Close(); err != nil {
		return errors.Wrap(err, "close temp file")
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", script)

	if out, err := cmd.CombinedOutput(); err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

func IsInstalled(path string) (bool, error) {
	if _, err := exec.LookPath(path); err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return false, nil
		}

		return false, errors.Wrap(err, "check if path is installed")
	}

	return true, nil
}
