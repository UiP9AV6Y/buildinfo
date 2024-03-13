package buildinfo

import (
	"fmt"
	"time"
)

const (
	// default value for EnvironmentInfo.User
	DefaultUser = "unknown"
	// default value for EnvironmentInfo.Host
	DefaultHost = "localhost"
)

const (
	userConcat = "@"
)

// EnvironmentInfo contains data about the build context of a project
type EnvironmentInfo struct {
	User string    `json:"user,omitempty"`
	Host string    `json:"host,omitempty"`
	Date time.Time `json:"date,omitempty"`
}

// NewEnvironmentInfo returns a EnvironmentInfo instance with default values
func NewEnvironmentInfo() *EnvironmentInfo {
	return &EnvironmentInfo{
		User: DefaultUser,
		Host: DefaultHost,
		Date: time.Now(),
	}
}

// String returns build user, -host, and -date information.
func (i *EnvironmentInfo) String() string {
	return fmt.Sprintf("(user=%s, host=%s, date=%s)", i.User, i.Host, i.Date)
}

// Clone creates an independant copy of itself.
func (i *EnvironmentInfo) Clone() *EnvironmentInfo {
	i2 := *i

	return &i2
}

// Equal compares the fields of this instance to the given one
func (i *EnvironmentInfo) Equal(o *EnvironmentInfo) bool {
	if i == nil || o == nil {
		return i == nil && o == nil
	}

	return i.User == o.User && i.Host == o.Host && i.Date == o.Date
}

// UserHost returns the User and Host value
// concatenated by an @ character
func (i *EnvironmentInfo) UserHost() string {
	return i.User + userConcat + i.Host
}
