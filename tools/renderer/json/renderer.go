package json

import (
	jsenc "encoding/json"
	"fmt"
	"strings"

	"github.com/UiP9AV6Y/buildinfo"
)

// JSON is a renderer.BuildRenderer implementation emitting
// JSON-formatted data
type JSON struct {
	indent string
}

// NewMinified returns a JSON renderer instance which does not
// use any indentation for RenderBuildInfo output
func NewMinified() *JSON {
	return New(0)
}

// New returns a JSON renderer instance using the provided
// indentation for RenderBuildInfo output. If it zero or below,
// the render output will also not contain any newlines, i.e.
// be minified
func New(indent int) *JSON {
	var i string

	if indent > 0 {
		i = strings.Repeat(" ", indent)
	}

	result := &JSON{
		indent: i,
	}

	return result
}

// String implements the fmt.Stringer interface
func (j *JSON) String() string {
	return fmt.Sprintf("(indent=%q)", j.indent)
}

// RenderBuildInfo implements the renderer.BuildRenderer interface
func (j *JSON) RenderBuildInfo(info *buildinfo.BuildInfo) ([]byte, error) {
	if j.indent == "" {
		return jsenc.Marshal(info)
	}

	return jsenc.MarshalIndent(info, "", j.indent)
}
