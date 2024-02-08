package buildinfo

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestVersionInfoEqual(t *testing.T) {
	type testCase struct {
		haveLeft, haveRight *VersionInfo
		want                bool
	}

	testCases := map[string]testCase{
		"nil": {
			want: true,
		},
		"default": {
			haveLeft:  NewVersionInfo(),
			haveRight: NewVersionInfo(),
			want:      true,
		},
		"custom": {
			haveLeft: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
			haveRight: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
			want: true,
		},
		"left nil": {
			haveRight: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
		},
		"right nil": {
			haveLeft: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
		},
		"version mismatch": {
			haveLeft: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
			haveRight: &VersionInfo{
				Version:  "0",
				Revision: "2",
				Branch:   "3",
			},
		},
		"revision mismatch": {
			haveLeft: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
			haveRight: &VersionInfo{
				Version:  "1",
				Revision: "0",
				Branch:   "3",
			},
		},
		"branch mismatch": {
			haveLeft: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "3",
			},
			haveRight: &VersionInfo{
				Version:  "1",
				Revision: "2",
				Branch:   "0",
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

func TestVersionInfoShortRevision(t *testing.T) {
	type testCase struct {
		have *VersionInfo
		want string
	}

	testCases := map[string]testCase{
		"empty": {
			have: &VersionInfo{},
		},
		"three": {
			have: &VersionInfo{
				Revision: "123",
			},
			want: "123",
		},
		"eight": {
			have: &VersionInfo{
				Revision: "12345678",
			},
			want: "12345678",
		},
		"ten": {
			have: &VersionInfo{
				Revision: "0123456789",
			},
			want: "01234567",
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got := tc.have.ShortRevision()

			assert.Equal(t, got, tc.want)
		})
	}
}
