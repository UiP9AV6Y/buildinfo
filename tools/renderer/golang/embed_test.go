package golang

import (
	"testing"

	"gotest.tools/v3/assert"
	"gotest.tools/v3/golden"

	"github.com/UiP9AV6Y/buildinfo"
)

func TestRenderBuildInfo(t *testing.T) {
	type testCase struct {
		havePkg       string
		haveGenerator string
		haveArgs      []string
		haveInfo      *buildinfo.BuildInfo
		wantError     bool
		want          string
	}

	testCases := map[string]testCase{
		"simple": {
			havePkg:       "golang",
			haveGenerator: "gotest",
			haveArgs:      DefaultArgs("/mock/src", "gotest"),
			want:          "embed.golden",
		},
		"args": {
			havePkg:       "golang",
			haveGenerator: "gotest",
			haveArgs:      []string{"--verbose", "--minify", "--log.level", "debug"},
			want:          "args.golden",
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			subject := NewArgs(tc.havePkg, tc.haveGenerator, tc.haveArgs...)
			got, err := subject.RenderBuildInfo(tc.haveInfo)

			if tc.wantError {
				assert.Assert(t, err != nil)
			} else {
				assert.Assert(t, err)
				golden.Assert(t, string(got), tc.want)
			}
		})
	}
}
