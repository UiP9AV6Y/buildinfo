package git

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os/exec"
	"strings"

	"github.com/UiP9AV6Y/buildinfo"
)

const (
	systemGit = "git"
	// git error message when show-toplevel yields no result
	errParse = "not a git repository"
)

// Error when trying to parse a directory that is not a Git repository
var ErrNoRepository = fs.ErrNotExist

// parser.VersionParser implementation using Git as data backend
type Git struct {
	cmd  string
	root string
}

// TrySystemParse calls TryParse using the Git command found in the PATH
func TrySystemParse(path string) (*Git, error) {
	return TryParse(systemGit, path)
}

// TryParse attempts to parse the given directory as Git repository
// using the provided command.
// If the given path does not seem to be a Git repository,
// ErrNoRepository is returned. All other errors are a result of file
// access problems or data corruption issues.
func TryParse(cmd, path string) (*Git, error) {
	realCmd, err := exec.LookPath(cmd)
	if err != nil {
		// might be a git repo, but we have no git commandline client
		return nil, ErrNoRepository
	}

	o, e, err := run(cmd, path, "rev-parse", "--show-toplevel")
	if strings.Contains(e, errParse) {
		// ignore the error type, as long as the error output
		// contains hints about the failure cause
		return nil, ErrNoRepository
	} else if err != nil {
		return nil, err
	}

	return New(realCmd, o), nil
}

// NewSystem creates a new parser.Parser instance using the provided
// directory as project root. the git executable is invoked as-is,
// relying on its presence in one of the PATH directories.
func NewSystem(root string) *Git {
	return New(systemGit, root)
}

// New creates a new parser.Parser instance using the provided
// directory as project root. the git executable is invoked using
// the provided path.
func New(cmd, root string) *Git {
	result := &Git{
		cmd:  cmd,
		root: root,
	}

	return result
}

// String implements the fmt.Stringer interface
func (g *Git) String() string {
	return fmt.Sprintf("(cmd=%s, root=%s)", g.cmd, g.root)
}

// Equal compares the fields of this instance to the given one
func (g *Git) Equal(o *Git) bool {
	if o == nil {
		return g == nil
	}

	return g.cmd == o.cmd && g.root == o.root
}

// ParseVersionInfo implements the parser.VersionParser interface
func (g *Git) ParseVersionInfo() (*buildinfo.VersionInfo, error) {
	result := buildinfo.NewVersionInfo()

	branch, err := g.run("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("Unable to determine current git branch: %w", err)
	} else if branch != "" {
		result.Branch = branch
	}

	revision, err := g.run("rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("Unable to determine git HEAD revision: %w", err)
	} else if revision != "" {
		result.Revision = revision
	}

	// ignore error in case the project has no tags
	version, _ := g.run("describe", "--tags", "--abbrev=0")
	if version != "" {
		result.Version = strings.TrimPrefix(version, "v")
	}

	return result, nil
}

func (g *Git) run(arg ...string) (string, error) {
	o, e, err := run(g.cmd, g.root, arg...)
	if err == nil {
		return o, nil
	} else if e != "" {
		return "", errors.New(e)
	}

	return "", err
}

func run(cmd, cwd string, arg ...string) (string, string, error) {
	argv := append([]string{"-C", cwd}, arg...)
	git := exec.Command(cmd, argv...)
	if git.Err != nil {
		return "", "", git.Err
	}

	stderr, err := git.StderrPipe()
	if err != nil {
		return "", "", err
	}

	stdout, err := git.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	if err := git.Start(); err != nil {
		return "", "", err
	}

	e, err := io.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}

	o, _ := io.ReadAll(stdout)
	if err != nil {
		return "", "", err
	}

	return strings.Trim(string(o), " \n\r"),
		strings.Trim(string(e), " \n\r"),
		git.Wait()
}
