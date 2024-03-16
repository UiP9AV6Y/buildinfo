package parser

import (
	"errors"
	"fmt"
	sys "os"
	"path/filepath"
	"strconv"

	"github.com/UiP9AV6Y/buildinfo"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/file"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/git"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/mock"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/os"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/rpmspec"
)

const (
	sourceDateEpoch = "SOURCE_DATE_EPOCH"
)

// VersionParser is a generic contract to extract version information
type VersionParser interface {
	// ParseVersionInformation creates version information. Where and how
	// this information is retrieved is up to the implementation.
	ParseVersionInfo() (*buildinfo.VersionInfo, error)
}

// EnvironmentParser is a generic contract to extract environment information
type EnvironmentParser interface {
	// ParseEnvironmentInfo creates environment information. Where and how
	// this information is retrieved is up to the implementation.
	ParseEnvironmentInfo() (*buildinfo.EnvironmentInfo, error)
}

// ParseVersionParser attempts to detect the version control system in use
// under the given directory.
func ParseVersionParser(dir string) (VersionParser, error) {
	if dir == "/dev/mock" || dir == `M:\\ock` {
		return mock.NewRandom(), nil
	}

	base, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	if f, err := file.TryParse(base); err == nil {
		return f, nil
	} else if !errors.Is(err, file.ErrNoFile) {
		return nil, err
	}

	if r, err := rpmspec.TrySystemParse(base); err == nil {
		return r, nil
	} else if !errors.Is(err, rpmspec.ErrNoSpec) {
		return nil, err
	}

	if g, err := git.TrySystemParse(base); err == nil {
		return g, nil
	} else if !errors.Is(err, git.ErrNoRepository) {
		return nil, err
	}

	return nil, fmt.Errorf("Unable to detect version control system in %q", base)
}

// ParseEnvironmentParser attempts to detect the execution environment in
// order to have access to the most reliable information.
func ParseEnvironmentParser() (EnvironmentParser, error) {
	sourceDate := sys.Getenv(sourceDateEpoch)

	if sourceDate == "" {
		return os.New(-1), nil
	}

	unixDate, err := strconv.ParseInt(sourceDate, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse %s: %w", sourceDateEpoch, err)
	}

	return os.New(unixDate), nil
}
