//go:build release

package cmd

import (
	"github.com/bytebase/bytebase/backend/common"
	"github.com/bytebase/bytebase/backend/component/config"
)

func activeProfile(dataDir string) *config.Profile {
	p := getBaseProfile(dataDir)
	p.Mode = common.ReleaseModeProd
	// DISABLED: Metrics collection disabled to prevent communication with external servers.
	// Metric connection key is not set to prevent any metrics collection.
	// Original key was: "so9lLwj5zLjH09sxNabsyVNYSsAHn68F"
	// p.MetricConnectionKey = "" // Always keep empty to disable metrics
	return p
}
