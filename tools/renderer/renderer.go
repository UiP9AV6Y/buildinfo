package generator

import (
	"github.com/UiP9AV6Y/buildinfo"
)

// BuildRenderer is a contract for marshalling build information.
type BuildRenderer interface {
	// Render transforms the given build information into
	// a textual representation.
	RenderBuildInfo(*buildinfo.BuildInfo) ([]byte, error)
}
