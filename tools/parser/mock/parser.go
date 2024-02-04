package mock

import (
	"fmt"

	"github.com/UiP9AV6Y/buildinfo"
)

// parser.VersionParser implementation which simulates a version control system
type Mock struct {
	info *buildinfo.VersionInfo
}

// TryParse creates a parser instance with version info populated
// from the given input. only non-empty values are actually
// stored in the resulting output.
func TryParse(version, revision, branch string) (*Mock, error) {
	result := buildinfo.NewVersionInfo()

	if version != "" {
		result.Version = version
	}

	if revision != "" {
		result.Revision = revision
	}

	if branch != "" {
		result.Branch = branch
	}

	return New(result), nil
}

// NewRandom creates a new mock parser instance which produces
// a new buildinfo.BuildInfo instance with every call to Mock#ParseBuildInfo()
func NewRandom() *Mock {
	return New(nil)
}

// New creates a new mock parser instance which returns the given
// buildinfo.VersionInfo instance with every call to Mock#ParseVersionInfo()
func New(info *buildinfo.VersionInfo) *Mock {
	result := &Mock{
		info: info,
	}

	return result
}

// String implements the fmt.Stringer interface
func (m *Mock) String() string {
	return fmt.Sprintf("(info=%s)", m.info)
}

// Equal compares the fields of this instance to the given one
func (m *Mock) Equal(o *Mock) bool {
	if o == nil {
		return m == nil
	}

	if o.info == nil {
		return m.info == nil
	}

	return o.info.Equal(m.info)
}

// ParseVersionInfo implements the parser.VersionParser interface
func (m *Mock) ParseVersionInfo() (*buildinfo.VersionInfo, error) {
	if m.info != nil {
		return m.info, nil
	}

	return buildinfo.NewVersionInfo(), nil
}
