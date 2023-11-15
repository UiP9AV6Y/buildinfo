package git

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/buildinfo"
)

func mockGitBin() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	mockPath := wd + "/testdata"

	os.Setenv("PATH", mockPath)

	return mockPath + "/git-mock.sh", nil
}

func TestTryParse(t *testing.T) {
	type testCase struct {
		haveCmd, havePath string
		wantError         bool
		want              *Git
	}

	gitBin, err := mockGitBin()
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]testCase{
		"not in PATH": testCase{
			haveCmd:   "git-notexists",
			wantError: true,
		},
		"no git repo": testCase{
			haveCmd:   "git-mock.sh",
			havePath:  "/mock/NOT_GIT_REPO",
			wantError: true,
		},
		"git error": testCase{
			haveCmd:   "git-mock.sh",
			havePath:  "/mock/FAIL",
			wantError: true,
		},
		"relative bin": testCase{
			haveCmd:  "git-mock.sh",
			havePath: "/mock/SHOW_TOPLEVEL",
			want:     New(gitBin, "/mock/src"),
		},
		"absolute bin": testCase{
			haveCmd:  gitBin,
			havePath: "/mock/SHOW_TOPLEVEL",
			want:     New(gitBin, "/mock/src"),
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
		have      *Git
		wantError bool
		want      *buildinfo.VersionInfo
	}

	gitBin, err := mockGitBin()
	if err != nil {
		t.Fatal(err)
	}

	testCases := map[string]testCase{
		"all parsed": testCase{
			have: New(gitBin, "/mock/PARSE_ALL"),
			want: &buildinfo.VersionInfo{
				Version:  "1.23.456",
				Revision: "deadbeefcafe",
				Branch:   "test_mock",
			},
		},
		"no tag": testCase{
			have: New(gitBin, "/mock/PARSE_TAG_FAIL"),
			want: &buildinfo.VersionInfo{
				Version:  "0.0.0",
				Revision: "deadbeefcafe",
				Branch:   "test_mock",
			},
		},
		"no rev": testCase{
			have:      New(gitBin, "/mock/PARSE_REV_FAIL"),
			wantError: true,
		},
		"no branch": testCase{
			have:      New(gitBin, "/mock/PARSE_BRANCH_FAIL"),
			wantError: true,
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
