package buildinfo

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestParse(t *testing.T) {
	type testCase struct {
		have      []byte
		wantError bool
		want      *BuildInfo
	}

	testCases := map[string]testCase{
		"nil": {
			have: nil,
			want: New(),
		},
		"no data": {
			have: []byte{},
			want: New(),
		},
		"empty string": {
			have: []byte(""),
			want: New(),
		},
		"empty object": {
			have: []byte("{}"),
			want: New(),
		},
		"trailing newline": {
			have: []byte("{}\n"),
			want: New(),
		},
		"full": {
			have: []byte(`{
				"version":"1",
				"revision":"2",
				"branch":"3",
				"user":"4",
				"host":"5",
				"date":"1970-01-01T00:00:00.00000000Z"
			}`),
			want: &BuildInfo{
				VersionInfo: &VersionInfo{
					Version:  "1",
					Revision: "2",
					Branch:   "3",
				},
				EnvironmentInfo: &EnvironmentInfo{
					User: "4",
					Host: "5",
					Date: time.Unix(0, 0),
				},
			},
		},
		"malformed": {
			have:      []byte("{"),
			wantError: true,
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got, err := Parse(tc.have)

			if tc.wantError {
				assert.Assert(t, err != nil)
			} else if tc.want == nil {
				assert.Assert(t, err)
				assert.Assert(t, got == nil)
			} else {
				assert.Assert(t, err)
				assert.Assert(t, tc.want.VersionInfo.Equal(got.VersionInfo), "want=%s; got=%s", tc.want, got)
			}
		})
	}
}

func TestBuildInfoEqual(t *testing.T) {
	type testCase struct {
		haveLeft, haveRight *BuildInfo
		want                bool
	}

	testCases := map[string]testCase{
		"nil": {
			want: true,
		},
		//		"default": {
		//			haveLeft:  New(),
		//			haveRight: New(),
		//			want:      true,
		//		},
		"custom": {
			haveLeft: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
			haveRight: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
			want: true,
		},
		"left nil": {
			haveRight: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
		},
		"right nil": {
			haveLeft: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
		},
		"version nil": {
			haveLeft: NewBuildInfo(nil, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
			haveRight: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
		},
		"version mismatch": {
			haveLeft: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
			haveRight: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "0",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
		},
		"environment nil": {
			haveLeft: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, nil),
			haveRight: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
		},
		"environment mismatch": {
			haveLeft: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			}),
			haveRight: NewBuildInfo(&VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			}, &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(1, 0),
			}),
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got := tc.haveLeft.Equal(tc.haveRight)

			assert.Equal(t, got, tc.want)
		})
	}
}
