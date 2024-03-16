package rpmspec

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/buildinfo"
)

func mockRPMSPECBin() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	mockPath := wd + "/testdata"

	os.Setenv("PATH", mockPath)

	return mockPath + "/rpmspec-mock.sh", nil
}

func TestTryParse(t *testing.T) {
	type testCase struct {
		haveCmd, havePath string
		wantError         bool
		want              *RPMSpec
	}

	rpmspecBin, err := mockRPMSPECBin()
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]testCase{
		"not in PATH": {
			haveCmd:   "rpmspec-notexists",
			wantError: true,
		},
		"no spec file": {
			haveCmd:   "rpmspec-mock.sh",
			havePath:  "testdata/nospec",
			wantError: true,
		},
		"broken spec file": {
			haveCmd:  "rpmspec-mock.sh",
			havePath: "testdata/broken",
			want:     New(rpmspecBin, "testdata/broken/broken.spec"),
		},
		"macro spec file": {
			haveCmd:  "rpmspec-mock.sh",
			havePath: "testdata/macro",
			want:     New(rpmspecBin, "testdata/macro/macro.spec"),
		},
		"minimal spec file": {
			haveCmd:  "rpmspec-mock.sh",
			havePath: "testdata/minimal",
			want:     New(rpmspecBin, "testdata/minimal/minimal.spec"),
		},
		"multiple spec files": {
			haveCmd:  "rpmspec-mock.sh",
			havePath: "testdata/multiple",
			want:     New(rpmspecBin, "testdata/multiple/multiple.spec"),
		},
		"relative bin": {
			haveCmd:  "rpmspec-mock.sh",
			havePath: "testdata/minimal",
			want:     New(rpmspecBin, "testdata/minimal/minimal.spec"),
		},
		"absolute bin": {
			haveCmd:  rpmspecBin,
			havePath: "testdata/minimal",
			want:     New(rpmspecBin, "testdata/minimal/minimal.spec"),
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got, err := TryParse(tc.haveCmd, tc.havePath)

			if tc.wantError {
				assert.Assert(t, err != nil)
			} else {
				assert.Assert(t, err)
				assert.Assert(t, tc.want.Equal(got), "want=%s; got=%s", tc.want, got)
			}
		})
	}
}

func TestParseVersionInfo(t *testing.T) {
	type testCase struct {
		have      *RPMSpec
		wantError bool
		want      *buildinfo.VersionInfo
	}

	rpmspecBin, err := mockRPMSPECBin()
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]testCase{
		"broken": {
			have:      New(rpmspecBin, "testdata/broken/broken.spec"),
			wantError: true,
		},
		"nospec": {
			have:      New(rpmspecBin, "testdata/nospec/nospec.spec"),
			wantError: true,
		},
		"macro": {
			have: New(rpmspecBin, "testdata/macro/macro.spec"),
			want: &buildinfo.VersionInfo{
				Version:  "1.2.3~19701230gitd5a3191",
				Revision: "1.rhel",
				Branch:   "trunk",
			},
		},
		"minimal": {
			have: New(rpmspecBin, "testdata/minimal/minimal.spec"),
			want: &buildinfo.VersionInfo{
				Version:  "1.0",
				Revision: "1",
				Branch:   "trunk",
			},
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got, err := tc.have.ParseVersionInfo()

			if tc.wantError {
				assert.Assert(t, err != nil)
			} else {
				assert.Assert(t, err)
				assert.Assert(t, tc.want.Equal(got), "want=%s; got=%s", tc.want, got)
			}
		})
	}
}
