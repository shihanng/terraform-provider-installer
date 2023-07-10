package brew

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/shihanng/terraform-provider-installer/internal/system"
	"github.com/shihanng/terraform-provider-installer/internal/xerrors"
)

var ErrFormulaNotFound = errors.New("formula not found")

func Install(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, "brew", args...)

	out, err := cmd.CombinedOutput()
	if err == nil {
		return nil
	}

	var exitError *exec.ExitError

	if errors.As(err, &exitError) {
		if exitError.ExitCode() == 1 && strings.Contains(string(out), "No available formula with the name") {
			return ErrFormulaNotFound
		}
	}

	return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
}

// InfoV2 is the JSON payload returned by brew info --json=v2.
type InfoV2 struct {
	Formulae []struct {
		Name string
		Tap  string
	}
	Casks []struct {
		Token string
		Tap   string
	}
}

type Info struct {
	Name   string
	IsCask bool
}

func (i *InfoV2) GetInfo() Info {
	var (
		info      Info
		name, tap string
	)

	if len(i.Casks) > 0 {
		name = i.Casks[0].Token
		tap = i.Casks[0].Tap
		info.IsCask = true
	} else if len(i.Formulae) > 0 {
		name = i.Formulae[0].Name
		tap = i.Formulae[0].Tap
	}

	info.Name = filepath.Join(tap, name)

	return info
}

func GetInfo(ctx context.Context, args []string) (Info, error) {
	cmd := exec.CommandContext(ctx, "brew", args...)

	// Check stdout only as there might be warnings in stderr, which break the unmarshalling later, like:
	// "Warning: Treating dash as a formula. For the cask, use homebrew/cask/dash"
	out, err := cmd.Output()
	if err != nil {
		return Info{}, errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	var infoV2 InfoV2

	if err := json.Unmarshal(out, &infoV2); err != nil {
		return Info{}, errors.Wrap(errors.WithDetail(err, string(out)), "failed to decode InfoV2")
	}

	return infoV2.GetInfo(), nil
}

func FindInstalled(ctx context.Context, name string) (string, error) {
	cmd := exec.CommandContext(ctx, "brew", "list", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "No available formula with the name") {
			return "", xerrors.ErrNotInstalled
		}

		if strings.Contains(string(out), "No available formula or cask with the name") {
			return "", xerrors.ErrNotInstalled
		}

		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	paths := strings.Split(string(out), "\n")

	return system.FindExecutablePath(paths) // nolint:wrapcheck
}

func FindCaskPath(ctx context.Context, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, "brew", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "is not installed") {
			return "", xerrors.ErrNotInstalled
		}

		return "", errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	lines := bytes.Split(out, []byte("\n"))

	for _, line := range lines {
		if bytes.HasSuffix(line, []byte(".app")) {
			return string(line), nil
		}
	}

	return "", nil
}

func Uninstall(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "brew", "uninstall", name)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(errors.WithDetail(err, string(out)), strings.Join(cmd.Args, " "))
	}

	return nil
}

// CmdOption of BrewCmd.
type CmdOption func(*Cmd)

// Cmd holds the arguments for brew commands.
type Cmd struct {
	Args []string
}

// NewCmd constructs a new Cmd instance for brew.
func NewCmd(action, pkg string, opts ...CmdOption) *Cmd {
	cmd := &Cmd{
		Args: []string{action, pkg},
	}

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func WithCask(isCask bool) CmdOption {
	return func(c *Cmd) {
		if isCask {
			c.Args = append(c.Args, "--cask")
		}
	}
}

func WithJSONV2() CmdOption {
	return func(c *Cmd) {
		c.Args = append(c.Args, "--json=v2")
	}
}
