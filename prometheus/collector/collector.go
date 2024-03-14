package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/UiP9AV6Y/buildinfo"
)

// New returns a collector that exports metrics
// using the provided data as information source.
func New(buildInfo *buildinfo.BuildInfo, program string) prometheus.Collector {
	help := "A metric with a constant '1' value labeled by version, revision, branch, goversion from which " +
		program + " was built, and the goos and goarch for the build."
	labels := prometheus.Labels{
		"version":   buildInfo.Version,
		"revision":  buildInfo.Revision,
		"branch":    buildInfo.Branch,
		"goversion": buildinfo.GoVersion,
		"goos":      buildinfo.GoOS,
		"goarch":    buildinfo.GoArch,
	}

	return prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Namespace:   program,
			Name:        "build_info",
			Help:        help,
			ConstLabels: labels,
		},
		func() float64 { return 1 },
	)
}
