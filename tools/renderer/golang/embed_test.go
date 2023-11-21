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
		haveInput     string
		haveInfo      *buildinfo.BuildInfo
		wantError     bool
		want          string
	}

	testCases := map[string]testCase{
		"simple": testCase{
			havePkg:       "golang",
			haveGenerator: "gotest",
			haveInput:     "/mock/src",
			want:          "embed.golden",
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			subject := New(tc.havePkg, tc.haveInput, tc.haveGenerator)
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
