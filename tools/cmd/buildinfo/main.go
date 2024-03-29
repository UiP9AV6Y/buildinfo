package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/UiP9AV6Y/buildinfo/tools/cmd/buildinfo/app"
	"github.com/UiP9AV6Y/buildinfo/tools/version"
)

func main() {
	os.Exit(run(os.Args, os.Stdout, os.Stderr))
}

func runHelp(fs *flag.FlagSet) int {
	fmt.Fprintf(fs.Output(), "Usage of %s:\n", fs.Name())
	fs.PrintDefaults()

	return 0
}

func runVersion(fs *flag.FlagSet) int {
	fmt.Fprintln(fs.Output(), version.Print(fs.Name()))

	return 0
}

func newLogger(lvl string, o io.Writer) (log.Logger, error) {
	l, err := level.Parse(lvl)
	if err != nil {
		return nil, err
	}
	logger := log.NewLogfmtLogger(o)
	logger = level.NewFilter(logger, level.Allow(l))

	return logger, nil
}

func run(a []string, o, e io.Writer) int {
	var l io.Writer
	n := filepath.Base(a[0])
	fs := flag.NewFlagSet(n, flag.ContinueOnError)
	app := app.New(n)
	help := fs.Bool("help", false, "Show the program usage and exit")
	info := fs.Bool("version", false, "Show the program version and exit")
	lvl := fs.String("log.level", "info", "Emit information about the internal processing")

	fs.StringVar(&app.Filename, "filename", os.Getenv("BUILDINFO_FILENAME"), "File path to write data to instead of STDOUT")
	fs.StringVar(&app.ProjectDir, "project-dir", os.Getenv("BUILDINFO_PROJECT_DIR"), "Project root directory to parse for version information")
	fs.StringVar(&app.Format, "generate", os.Getenv("BUILDINFO_GENERATE"), "Data generator to use for build information processing")
	fs.StringVar(&app.Namespace, "generate.namespace", os.Getenv("GOPACKAGE"), "Code namespace if output directory is not suitable/detectable")
	fs.StringVar(&app.VersionParser, "parser.version", os.Getenv("BUILDINFO_PARSER_VERSION"), "Version parser strategy to use. Valid values include git, file, and mock. If not specified, an appropriate provider will be selected")
	fs.StringVar(&app.GitExe, "git.exe", os.Getenv("BUILDINFO_GIT_EXE"), "Filesystem location for the git executable")
	fs.StringVar(&app.MockVersion, "mock.version", os.Getenv("BUILDINFO_MOCK_VERSION"), "Version value for the mock strategy")
	fs.StringVar(&app.MockRevision, "mock.revision", os.Getenv("BUILDINFO_MOCK_REVISION"), "Revision value for the mock strategy")
	fs.StringVar(&app.MockBranch, "mock.branch", os.Getenv("BUILDINFO_MOCK_BRANCH"), "Branch value for the mock strategy")
	fs.SetOutput(io.Discard) // discard any output until after parse, as it writes error messages on its own

	if err := fs.Parse(a[1:]); err != nil {
		fs.SetOutput(o)
		if errors.Is(err, flag.ErrHelp) {
			return runHelp(fs)
		} else {
			fmt.Fprintln(e, err)

			return 1
		}
	} else {
		fs.SetOutput(o)
	}

	if *info {
		return runVersion(fs)
	} else if *help {
		return runHelp(fs)
	}

	if app.Stdout() {
		l = e
	} else {
		l = o
	}

	logger, err := newLogger(*lvl, l)
	if err != nil {
		fmt.Fprintln(e, err)
		return 1
	}

	if err := app.Run(logger); err != nil {
		fmt.Fprintln(e, err)
		return 1
	}

	return 0
}
