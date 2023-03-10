// Package autodiscovery describes the interfaces which the various mechanisms
// implement.
package autodiscovery

import (
	"github.com/grafana/agent/component/discovery"
)

// Autodiscovery is the base interface for an autodiscovery mechanism.
// Implementations may also implement extension interfaces (named
// <Extension>Autodiscovery) to implement extra known behavior.
type Autodiscovery interface {
	Run() (*Result, error)
}

// Result ???
type Result struct {
	RiverConfig    string
	MetricsExport  string
	MetricsTargets []discovery.Target
	LogfileTargets []discovery.Target
}
