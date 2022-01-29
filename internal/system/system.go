package system

import (
	"os"

	"github.com/cockroachdb/errors"
)

var errNotExecutable = errors.New("could not find executable path")

func FindExecutablePath(paths []string) (string, error) {
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

	return "", errNotExecutable
}
