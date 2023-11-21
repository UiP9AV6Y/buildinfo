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
	"github.com/UiP9AV6Y/buildinfo/tools/renderer/golang"
	"github.com/UiP9AV6Y/buildinfo/tools/renderer/json"
)

// Application is a logic router implementation
type Application struct {
	Filename, ProjectDir string
	Format, Namespace    string

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
	i, err := a.buildInfo(logger)
	if err != nil {
		return err
	}

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

func (a *Application) buildInfo(logger log.Logger) (*buildinfo.BuildInfo, error) {
	level.Info(logger).Log("msg", "Parsing build information", "input", a.ProjectDir)

	vp, err := parser.ParseVersionParser(a.ProjectDir)
	if err != nil {
		return nil, err
	}

	level.Info(logger).Log("msg", "Parsing VCS information", "parser", &lazyReflect{v: vp})

	v, err := vp.ParseVersionInfo()
	if err != nil {
		return nil, err
	}

	ep, err := parser.ParseEnvironmentParser()
	if err != nil {
		return nil, err
	}

	level.Info(logger).Log("msg", "Parsing environment information", "parser", &lazyReflect{v: ep})

	e, err := ep.ParseEnvironmentInfo()
	if err != nil {
		return nil, err
	}

	return buildinfo.NewBuildInfo(v, e), nil
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
