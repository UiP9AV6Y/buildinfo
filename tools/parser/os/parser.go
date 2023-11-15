package os

import (
	"fmt"
	sys "os"
	"os/user"
	"time"

	"github.com/UiP9AV6Y/buildinfo"
)

const (
	// Static value to use for reproducible builds
	ReproduciblePlaceholder = "reproducible"
)

// parser.EnvironmentParser implementation which retrieves information using
// operating system utilities
type OperatingSystem struct {
	buildTimestamp int64
}

// New creates a new environment parser using the golang stdlib.
// The buildTimestamp determines the calculation behaviour of
// the resulting buildinfo.EnvironmentInfo. a positive value
// indicates the intend to create reproducible output, in which
// case both the buildinfo.EnvironmentInfo#Host and
// buildinfo.EnvironmentInfo#User are set to #ReproduciblePlaceholder;
// otherwise the field value are calculated using stdlib functionality.
func New(buildTimestamp int64) *OperatingSystem {
	result := &OperatingSystem{
		buildTimestamp: buildTimestamp,
	}

	return result
}

// String implements the fmt.Stringer interface
func (s *OperatingSystem) String() string {
	return fmt.Sprintf("(buildTimestamp=%d)", s.buildTimestamp)
}

// Equal compares the fields of this instance to the given one
func (s *OperatingSystem) Equal(o *OperatingSystem) bool {
	if o == nil {
		return s == nil
	}

	return s.buildTimestamp == o.buildTimestamp
}

// ParseEnvironmentInfo implements the parser.EnvironmentParser interface
func (s *OperatingSystem) ParseEnvironmentInfo() (*buildinfo.EnvironmentInfo, error) {
	result := &buildinfo.EnvironmentInfo{
		User: ReproduciblePlaceholder,
		Host: ReproduciblePlaceholder,
	}

	if s.buildTimestamp >= 0 {
		result.Date = time.Unix(s.buildTimestamp, 0)

		return result, nil
	}

	result.Date = time.Now()

	hostname, err := sys.Hostname()
	if err != nil {
		return nil, err
	}

	result.Host = hostname

	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	result.User = usr.Username

	return result, nil
}
