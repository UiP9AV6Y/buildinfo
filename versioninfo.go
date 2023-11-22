package buildinfo

import (
	"fmt"
)

const (
	// default value for VersionInfo.Version
	DefaultVersion = "0.0.0"
	// default value for VersionInfo.Revision
	DefaultRevision = "HEAD"
	// default value for VersionInfo.Branch
	DefaultBranch = "trunk"
)

const (
	versionConcat = "-"
)

// VersionInfo contains data about the project state in a version control system
type VersionInfo struct {
	Version  string `json:"version,omitempty"`
	Revision string `json:"revision,omitempty"`
	Branch   string `json:"branch,omitempty"`
}

// NewVersionInfo returns a VersionInfo instance with default values
func NewVersionInfo() *VersionInfo {
	return &VersionInfo{
		Version:  DefaultVersion,
		Revision: DefaultRevision,
		Branch:   DefaultBranch,
	}
}

// String returns version, branch and revision information.
func (i *VersionInfo) String() string {
	return fmt.Sprintf("(version=%s, branch=%s, revision=%s)", i.Version, i.Branch, i.Revision)
}

// Equal compares the fields of this instance to the given one
func (i *VersionInfo) Equal(o *VersionInfo) bool {
	if o == nil {
		return i == nil
	}

	return i.Version == o.Version && i.Revision == o.Revision && i.Branch == o.Branch
}

// ShortRevision returns the truncated Revision.
// If that value is less tha 8 characters long, the result
// is the Revision value itself
func (i *VersionInfo) ShortRevision() string {
	if len(i.Revision) < 8 {
		return i.Revision
	}

	return i.Revision[0:8]
}

// VersionRevision returns the Version and (short) Revision value
// concatenated by an hyphen character
func (i *VersionInfo) VersionRevision() string {
	return i.Version + versionConcat + i.ShortRevision()
}
