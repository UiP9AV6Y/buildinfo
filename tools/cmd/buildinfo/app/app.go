package app

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/buildinfo"
	"github.com/UiP9AV6Y/buildinfo/tools/parser"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/file"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/git"
	"github.com/UiP9AV6Y/buildinfo/tools/parser/mock"
	"github.com/UiP9AV6Y/buildinfo/tools/renderer/golang"
	"github.com/UiP9AV6Y/buildinfo/tools/renderer/json"
)

// Application is a logic router implementation
type Application struct {
	Filename, ProjectDir                  string
	Format, Namespace                     string
	VersionParser                         string
	GitExe                                string
	MockVersion, MockRevision, MockBranch string

	name string
}

// New create a new Application instance
func New(name string) *Application {
	result := &Application{
		name: name,
	}

	return result
}

// String returns the application name
func (a *Application) String() string {
	return a.name
}

// Stdin checks if the instructions request output to STDOUT
func (a *Application) Stdout() bool {
	return a.Filename == "" || a.Filename == "-"
}

// Stdin checks if the instructions request input from STDIN
func (a *Application) Stdin() bool {
	return a.ProjectDir == "" || a.ProjectDir == "-"
}

// Run is a logic switch operating on the configured Format to produce
func (a *Application) Run(logger log.Logger) error {
	switch a.Format {
	case "", "version", "buildinfo", "metadata":
		return a.GenerateBuildInfo(logger)
	case "golang-embed":
		return a.GenerateGolangEmbed(logger)
	default:
		return fmt.Errorf("Invalid generator instruction %q", a.Format)
	}
}

// GenerateBuildInfo parses the buildinfo data and renders them using JSON.
func (a *Application) GenerateBuildInfo(logger log.Logger) error {
	level.Debug(logger).Log("msg", "Parsing build information", "input", a.ProjectDir)

	v, err := a.versionInfo(logger)
	if err != nil {
		return err
	}

	e, err := a.environmentInfo(logger)
	if err != nil {
		return err
	}

	i := buildinfo.NewBuildInfo(v, e)
	r := json.NewMinified()

	return a.write(func(o string, w io.Writer) error {
		b, err := r.RenderBuildInfo(i)
		if err != nil {
			return err
		}

		level.Info(logger).Log("msg", "Writing BuildInfo data", "output", o)

		_, err = w.Write(b)
		return err
	})
}

// GenerateGolangEmbed renders Golang code with buildinfo embed instructions.
func (a *Application) GenerateGolangEmbed(logger log.Logger) error {
	level.Debug(logger).Log("msg", "Rendering Golang embed code", "name", a.name, "namespace", a.namespace())

	e := golang.New(a.namespace(), a.input(), a.name)

	return a.write(func(o string, w io.Writer) error {
		b, err := e.RenderBuildInfo(nil)
		if err != nil {
			return err
		}

		level.Info(logger).Log("msg", "Writing Golang embed code", "output", o)

		_, err = w.Write(b)
		return err
	})
}

func (a *Application) versionInfo(logger log.Logger) (*buildinfo.VersionInfo, error) {
	var vp parser.VersionParser
	var err error

	switch a.VersionParser {
	case "":
		vp, err = parser.ParseVersionParser(a.ProjectDir)
	case "file":
		vp, err = file.TryParse(a.ProjectDir)
	case "git":
		if a.GitExe != "" {
			vp, err = git.TryParse(a.GitExe, a.ProjectDir)
		} else {
			vp, err = git.TrySystemParse(a.ProjectDir)
		}
	case "mock":
		vp, err = mock.TryParse(a.MockVersion, a.MockRevision, a.MockBranch)
	default:
		err = fmt.Errorf("Invalid version parser %q", a.VersionParser)
	}

	if err != nil {
		return nil, err
	}

	level.Info(logger).Log("msg", "Parsing version information", "parser", &lazyReflect{v: vp})

	return vp.ParseVersionInfo()
}

func (a *Application) environmentInfo(logger log.Logger) (*buildinfo.EnvironmentInfo, error) {
	ep, err := parser.ParseEnvironmentParser()
	if err != nil {
		return nil, err
	}

	level.Info(logger).Log("msg", "Parsing environment information", "parser", &lazyReflect{v: ep})

	return ep.ParseEnvironmentInfo()
}

func (a *Application) write(consumer func(string, io.Writer) error) error {
	var o string
	var w io.Writer
	if a.Stdout() {
		o = "STDOUT"
		w = os.Stdout
	} else {
		p, f, err := mkdirFile(a.Filename)
		if err != nil {
			return err
		}

		b := bufio.NewWriter(f)

		o = p
		w = b
		defer func() {
			b.Flush()
			f.Close()
		}()
	}

	return consumer(o, w)
}

func (a *Application) input() string {
	if !a.Stdin() {
		return a.ProjectDir
	}

	return "."
}

func (a *Application) namespace() string {
	if a.Namespace != "" {
		return a.Namespace
	} else if a.Stdout() {
		d, err := os.Getwd()
		if err != nil {
			return a.name
		}

		return filepath.Base(d)
	}

	f, err := filepath.Abs(a.Filename)
	if err != nil {
		return a.name
	}

	return filepath.Base(filepath.Dir(f))
}

func mkdirFile(i string) (string, *os.File, error) {
	f, err := filepath.Abs(i)
	if err != nil {
		return "", nil, err
	}

	d := filepath.Dir(f)
	err = os.MkdirAll(d, 0755)
	if err != nil {
		return "", nil, err
	}

	w, err := os.Create(f)
	if err != nil {
		return "", nil, err
	}

	return f, w, nil
}
