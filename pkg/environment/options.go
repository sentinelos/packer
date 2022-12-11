package environment

import (
	"time"

	"github.com/moby/buildkit/client/llb"

	"github.com/sentinelos/packer/pkg/types"
)

// Options for packer.
type Options struct {
	BuildPlatform    Platform
	TargetPlatform   Platform
	Target           string
	CommonPrefix     string
	ProxyEnv         *llb.ProxyEnv
	SourceDateEpoch  time.Time
	CacheIDNamespace string
}

// GetVariables returns set of variables set for options.
func (options *Options) GetVariables() types.Variables {
	return Default().
		Merge(options.BuildPlatform.BuildVariables()).
		Merge(options.TargetPlatform.TargetVariables())
}
