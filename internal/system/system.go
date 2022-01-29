package system

import (
	"os/exec"

	"github.com/cockroachdb/errors"
)

var errNotExecutable = errors.New("could not find executable path")

func FindExecutablePath(paths []string) (string, error) {
	for _, path := range paths {
		_, err := exec.LookPath(path)
		if err != nil {
			continue
		}

		return path, nil
	}

	return "", errNotExecutable
}
