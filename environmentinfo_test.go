package buildinfo

import (
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestEnvironmentInfoEqual(t *testing.T) {
	type testCase struct {
		haveLeft, haveRight *EnvironmentInfo
		want                bool
	}

	testCases := map[string]testCase{
		"nil": {
			want: true,
		},
		//		"default": {
		//			haveLeft:  NewEnvironmentInfo(),
		//			haveRight: NewEnvironmentInfo(),
		//			want:      true,
		//		},
		"custom": {
			haveLeft: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
			haveRight: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
			want: true,
		},
		"left nil": {
			haveRight: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
		},
		"right nil": {
			haveLeft: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
		},
		"user mismatch": {
			haveLeft: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
			haveRight: &EnvironmentInfo{
				User: "0",
				Host: "2",
				Date: time.Unix(0, 0),
			},
		},
		"host mismatch": {
			haveLeft: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
			haveRight: &EnvironmentInfo{
				User: "1",
				Host: "0",
				Date: time.Unix(0, 0),
			},
		},
		"date mismatch": {
			haveLeft: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(0, 0),
			},
			haveRight: &EnvironmentInfo{
				User: "1",
				Host: "2",
				Date: time.Unix(1, 0),
			},
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got := tc.haveLeft.Equal(tc.haveRight)

			assert.Equal(t, got, tc.want)
		})
	}
}
