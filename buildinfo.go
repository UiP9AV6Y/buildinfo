package buildinfo

import (
	"encoding/json"
	"fmt"
	"runtime"
)

const (
	infoConcat = "; "
	// infoFmt contains the format used by BuildInfo.Print
	infoFmt = `%s, version %s (branch: %s, revision: %s)
  build user:       %s
  build host:       %s
  build date:       %s
  go version:       %s
  platform:         %s`
)

var (
	GoVersion = runtime.Version()
	GoOS      = runtime.GOOS
	GoArch    = runtime.GOARCH
)

// BuildInfo is a container for version- and environment information
type BuildInfo struct {
	*VersionInfo
	*EnvironmentInfo
}

// NewBuildInfo returns a BuildInfo instance using the provided values
func NewBuildInfo(v *VersionInfo, e *EnvironmentInfo) *BuildInfo {
	return &BuildInfo{
		VersionInfo:     v,
		EnvironmentInfo: e,
	}
}

// New returns a BuildInfo instance with default values
func New() *BuildInfo {
	return NewBuildInfo(NewVersionInfo(), NewEnvironmentInfo())
}

// Parse unmarshals the given JSON byte data into
// a BuildInfo instance
func Parse(info []byte) (*BuildInfo, error) {
	result := New()
	if err := json.Unmarshal(info, result); err != nil {
		return nil, err
	}

	return result, nil
}

// TryParse calls Parse and returns a
// default instance in case of an error
func TryParse(info []byte) *BuildInfo {
	result, err := Parse(info)
	if err != nil {
		return New()
	}

	return result
}

// MustParse calls Parse and panics on any error
func MustParse(info []byte) *BuildInfo {
	result, err := Parse(info)
	if err != nil {
		panic(err)
	}

	return result
}

// String returns the Version- and Environment information concatenated
func (i *BuildInfo) String() string {
	return i.VersionInfo.String() + infoConcat + i.EnvironmentInfo.String()
}

// Equal compares the fields of this instance to the given one
func (i *BuildInfo) Equal(o *BuildInfo) bool {
	if o == nil {
		return i == nil
	}

	return i.VersionInfo.Equal(o.VersionInfo) && i.EnvironmentInfo.Equal(o.EnvironmentInfo)
}

// JSON is a wrapper for json.Marshal using the instance as parameter
func (i *BuildInfo) JSON() ([]byte, error) {
	return json.Marshal(i)
}

// Print returns version- and environment information.
func (i *BuildInfo) Print(program string) string {
	return fmt.Sprintf(infoFmt,
		program,
		i.VersionInfo.Version,
		i.VersionInfo.Branch,
		i.VersionInfo.ShortRevision(),
		i.EnvironmentInfo.User,
		i.EnvironmentInfo.Host,
		i.EnvironmentInfo.Date,
		GoVersion,
		GoOS+"/"+GoArch,
	)
}
