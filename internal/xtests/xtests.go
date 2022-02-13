package xtests

import (
	"os"

	"github.com/cockroachdb/errors"
)

func SetupASDFDataDir(path string) (func(), error) {
	reset, err := envSetter(map[string]string{
		"ASDF_DATA_DIR": path,
	})
	if err != nil {
		return nil, err
	}

	return reset, nil
}

func envSetter(envs map[string]string) (func(), error) {
	originals := map[string]string{}

	for name, value := range envs {
		if val, ok := os.LookupEnv(name); ok {
			originals[name] = val
		}

		if err := os.Setenv(name, value); err != nil {
			return nil, errors.Wrap(err, "set environment variables for test")
		}
	}

	return func() {
		for name := range envs {
			original, ok := originals[name]
			if ok {
				_ = os.Setenv(name, original)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}, nil
}
