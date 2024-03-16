package rpmspec

import (
	"fmt"
	"io/fs"
	"os/exec"
	"path/filepath"

	"github.com/UiP9AV6Y/buildinfo"
	"github.com/UiP9AV6Y/buildinfo/tools/util"
)

const (
	// spec field to render/extract for the version
	versionMacro = "%{version}"
	// spec field to render/extract for the revision
	revisionMacro = "%{release}"
	systemCommand = "rpmspec"
)

var (
	// Error when no RPM spec file or `rpmspec` command was found
	ErrNoSpec = fs.ErrNotExist
)

// parser.VersionParser implementation rendering and parsing a RPM spec file
type RPMSpec struct {
	cmd  string
	file string
}

// TrySystemParse calls TryParse using the rpmspec command found in the PATH
func TrySystemParse(path string) (*RPMSpec, error) {
	return TryParse(systemCommand, path)
}

// TryParse attempts to parse the given directory for a RPM spec file.
// If no file was found or the given command was not found
// ErrNoRepository is returned. All other errors are a result of file
// access problems.
func TryParse(cmd, path string) (*RPMSpec, error) {
	realCmd, err := exec.LookPath(cmd)
	if err != nil {
		// unable to parse spec file without `rpmspec`
		return nil, ErrNoSpec
	}

	pattern := filepath.Join(path, "*.spec")
	haystack, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	} else if haystack == nil || len(haystack) == 0 {
		return nil, ErrNoSpec
	}

	return New(realCmd, haystack[0]), nil
}

// NewSystem creates a new parser.Parser instance using the provided
// RPM spec file. the rpmspec executable is invoked as-is,
// relying on its presence in one of the PATH directories.
func NewSystem(file string) *RPMSpec {
	return New(systemCommand, file)
}

// New creates a new parser.Parser instance using the provided
// RPM spec file. the rpmspec executable is invoked using
// the provided path.
func New(cmd, file string) *RPMSpec {
	result := &RPMSpec{
		cmd:  cmd,
		file: file,
	}

	return result
}

// String implements the fmt.Stringer interface
func (s *RPMSpec) String() string {
	return fmt.Sprintf("(cmd=%s, spec=%s)", s.cmd, s.file)
}

// Equal compares the fields of this instance to the given one
func (s *RPMSpec) Equal(o *RPMSpec) bool {
	if o == nil {
		return s == nil
	}

	return s.cmd == o.cmd && s.file == o.file
}

// ParseVersionInfo implements the parser.VersionParser interface
func (s *RPMSpec) ParseVersionInfo() (*buildinfo.VersionInfo, error) {
	result := buildinfo.NewVersionInfo()

	if version, err := s.rpmspecQuery(versionMacro); err != nil {
		return nil, err
	} else if version != "" {
		result.Version = version
	}

	if revision, err := s.rpmspecQuery(revisionMacro); err != nil {
		return nil, err
	} else if revision != "" {
		result.Revision = revision
	}

	return result, nil
}

func (s *RPMSpec) rpmspecQuery(query string) (string, error) {
	return util.RunCmd(s.cmd, []string{"--query", "--queryformat", query, s.file})
}
