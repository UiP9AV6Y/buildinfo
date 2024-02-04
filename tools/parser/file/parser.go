package file

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/UiP9AV6Y/buildinfo"
)

const (
	Filename    = "VERSION"
	AltFilename = "VERSION.txt"
)

const (
	// we assume the information parts are separated by this
	versionConcat = "-"
)

var (
	// Error when trying to parse a directory that does not contain a version file
	ErrNoFile = fs.ErrNotExist
	// ErrMalformedVersion is the error used when parsing invalid version information.
	ErrMalformedVersion = errors.New("malformed version information")
	// ErrMalformedRevision is the error used when parsing invalid revision information.
	ErrMalformedRevision = errors.New("malformed revision information")
	// ErrMalformedBranch is the error used when parsing invalid branch information.
	ErrMalformedBranch = errors.New("malformed branch information")
)

// parser.VersionParser implementation reading information from a file
type File struct {
	file string
}

func TryParse(path string) (*File, error) {
	var file string

	file = filepath.Join(path, Filename)
	if _, err := os.Stat(file); err == nil {
		return New(file), nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	file = filepath.Join(path, AltFilename)
	if _, err := os.Stat(file); err == nil {
		return New(file), nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return nil, ErrNoFile
}

// New creates a new parser.Parser instance using the provided
// file as version information source.
func New(file string) *File {
	result := &File{
		file: file,
	}

	return result
}

// String implements the fmt.Stringer interface
func (f *File) String() string {
	return fmt.Sprintf("(file=%s)", f.file)
}

// Equal compares the fields of this instance to the given one
func (f *File) Equal(o *File) bool {
	if o == nil {
		return f == nil
	}

	return o.file == f.file
}

// ParseVersionInfo implements the parser.VersionParser interface
func (f *File) ParseVersionInfo() (*buildinfo.VersionInfo, error) {
	h, err := os.Open(f.file)
	if err != nil {
		return nil, err
	}

	defer h.Close()

	b, err := io.ReadAll(h)
	if err != nil {
		return nil, err
	}

	return ParseVersionInfo(bytes.TrimSpace(b))
}

// ParseVersionInfo extract version information from the
// provided input. it generally can been seen as the inverse
// of VersionInfo.VersionRevision()
func ParseVersionInfo(info []byte) (*buildinfo.VersionInfo, error) {
	result := buildinfo.NewVersionInfo()
	parts := bytes.Split(info, []byte(versionConcat))

	if len(parts) > 0 {
		if len(parts[0]) > 0 {
			result.Version = string(parts[0])
		} else {
			return nil, ErrMalformedVersion
		}
	} else {
		return nil, ErrMalformedVersion
	}

	if len(parts) > 1 {
		if len(parts[1]) > 0 {
			result.Revision = string(parts[1])
		} else {
			return nil, ErrMalformedRevision
		}
	} // else optional

	if len(parts) > 2 {
		if len(parts[2]) > 0 {
			result.Branch = string(parts[2])
		} else {
			return nil, ErrMalformedBranch
		}
	} // else optional

	return result, nil
}
