package file

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/buildinfo"
)

func TestTryParse(t *testing.T) {
	type testCase struct {
		havePath  string
		wantError bool
		want      *File
	}

	testCases := map[string]testCase{
		"not found": testCase{
			havePath:  "testdata/not_found",
			wantError: true,
		},
		"unsupported": testCase{
			havePath:  "testdata/unsupported",
			wantError: true,
		},
		"VERSION.txt": testCase{
			havePath: "testdata/VERSION.txt",
			want:     New("testdata/VERSION.txt/VERSION.txt"),
		},
		"VERSION": testCase{
			havePath: "testdata/VERSION",
			want:     New("testdata/VERSION/VERSION"),
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got, err := TryParse(tc.havePath)

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
		have      *File
		wantError bool
		want      *buildinfo.VersionInfo
	}

	testCases := map[string]testCase{
		"not found": testCase{
			have:      New("testdata/variants/not_found"),
			wantError: true,
		},
		"empty": testCase{
			have:      New("testdata/variants/empty"),
			wantError: true,
		},
		"full": testCase{
			have: New("testdata/variants/full"),
			want: &buildinfo.VersionInfo{
				Version:  "999.99.9",
				Revision: "EOL",
				Branch:   "archive",
			},
		},
		"minimal": testCase{
			have: New("testdata/variants/minimal"),
			want: &buildinfo.VersionInfo{
				Version:  "1.2.3",
				Revision: "HEAD",
				Branch:   "trunk",
			},
		},
		"newline": testCase{
			have: New("testdata/variants/newline"),
			want: &buildinfo.VersionInfo{
				Version:  "2.4.5",
				Revision: "HEAD",
				Branch:   "trunk",
			},
		},
		"revision": testCase{
			have: New("testdata/variants/revision"),
			want: &buildinfo.VersionInfo{
				Version:  "1.2",
				Revision: "3",
				Branch:   "trunk",
			},
		},
		"spaced": testCase{
			have: New("testdata/variants/spaced"),
			want: &buildinfo.VersionInfo{
				Version:  "0.1.2",
				Revision: "HEAD",
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
