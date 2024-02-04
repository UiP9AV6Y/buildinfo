package golang

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/UiP9AV6Y/buildinfo"
)

const (
	// Code template for the Golang-Embed renderer
	EmbedText = `// Code generated by {{ .Generator }}. DO NOT EDIT.
package {{ .Pkg }}

import (
	_ "embed"

	"github.com/UiP9AV6Y/buildinfo"
)

//{{ "go:" }}generate {{ .Generator }} --generate buildinfo --filename {{ .Generator }}.json --project-dir {{ .Input }}
//{{ "go:" }}embed {{ .Generator }}.json
var embedInfo []byte
var buildInfo *buildinfo.BuildInfo

func init() {
	buildInfo = buildinfo.MustParse(embedInfo)
}

// Version returns the application release information
func Version() string {
	return buildInfo.VersionInfo.VersionRevision()
}

// Print returns version- and environment information.
func Print(program string) string {
	return buildInfo.Print(program)
}
`
)

// Embed is a code renderer using the internal fields as template variables.
type Embed map[string]string

// New returns a new Embed instance with the given template properties defined.
func New(pkg, input, generator string) Embed {
	result := map[string]string{
		"Pkg":       pkg,
		"Input":     input,
		"Generator": generator,
	}

	return result
}

// String implements the fmt.Stringer interface
func (e Embed) String() string {
	pairs := make([]string, 0, len(e))
	for k, v := range e {
		pairs = append(pairs, k+"="+v)
	}

	return "(" + strings.Join(pairs, ", ") + ")"
}

// RenderBuildInfo implements the renderer.BuildRenderer interface
func (e Embed) RenderBuildInfo(_ *buildinfo.BuildInfo) ([]byte, error) {
	tmpl, err := template.New("golang-embed").Parse(EmbedText)
	if err != nil {
		return nil, err
	}

	store := make([]byte, 0, len(EmbedText))
	buf := bytes.NewBuffer(store)
	err = tmpl.Execute(buf, e)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
