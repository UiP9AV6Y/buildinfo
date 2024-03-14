package collector

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"gotest.tools/v3/assert"

	"github.com/UiP9AV6Y/buildinfo"
)

func TestNew(t *testing.T) {
	golabels := fmt.Sprintf("goarch=%q,goos=%q,goversion=%q",
		runtime.GOARCH, runtime.GOOS, runtime.Version())
	vi := &buildinfo.VersionInfo{
		Version:  "1.2.3",
		Revision: "HEAD",
		Branch:   "trunk",
	}
	ei := &buildinfo.EnvironmentInfo{
		User: "root",
		Host: "localhost",
		Date: time.Unix(0, 0),
	}
	bi := buildinfo.NewBuildInfo(vi, ei)
	got := New(bi, "test")
	want := strings.NewReader(`# HELP test_build_info A metric with a constant '1' value labeled by version, revision, branch, goversion from which test was built, and the goos and goarch for the build.
# TYPE test_build_info gauge
test_build_info{branch="trunk",` + golabels + `,revision="HEAD",version="1.2.3"} 1
`)

	assert.NilError(t, testutil.CollectAndCompare(got, want))
}
