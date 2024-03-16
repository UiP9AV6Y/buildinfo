package util

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestRunCmd(t *testing.T) {
	type testCase struct {
		haveCmd   string
		haveArgv  []string
		wantError string
		want      string
	}

	testCases := map[string]testCase{
		"nothing success": {
			haveCmd:  "/bin/true",
			haveArgv: []string{},
		},
		"nothing error": {
			haveCmd:   "/bin/false",
			haveArgv:  []string{},
			wantError: "exit status 1",
		},
		"stdout": {
			haveCmd:  "echo",
			haveArgv: []string{"hello", "world"},
			want:     "hello world",
		},
		"stderr": {
			haveCmd:   "/bin/sh",
			haveArgv:  []string{"-c", "echo error: foo >&2; false"},
			wantError: "error: foo",
		},
		"both": {
			haveCmd:   "/bin/sh",
			haveArgv:  []string{"-c", "echo transaction complete; echo error: foo >&2; false"},
			wantError: "error: foo",
		},
		"ignore output": {
			haveCmd:   "/bin/sh",
			haveArgv:  []string{"-c", "echo test; false"},
			wantError: "exit status 1",
		},
		"ignore error": {
			haveCmd:  "/bin/sh",
			haveArgv: []string{"-c", "echo transaction complete; echo error: foo >&2"},
			want:     "transaction complete",
		},
		"no executable": {
			haveCmd:   "/opt/fail-inc/bin/foo-bar",
			haveArgv:  []string{},
			wantError: "fork/exec /opt/fail-inc/bin/foo-bar: no such file or directory",
		},
	}

	for ctx, tc := range testCases {
		t.Run(ctx, func(t *testing.T) {
			got, err := RunCmd(tc.haveCmd, tc.haveArgv)

			if tc.wantError != "" {
				assert.Error(t, err, tc.wantError)
			} else {
				assert.Assert(t, err)
				assert.Assert(t, tc.want == got, "want=%s; got=%s", tc.want, got)
			}
		})
	}
}
