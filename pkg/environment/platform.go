package environment

import (
	"fmt"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/sentinelos/packer/pkg/types"
)

// Platforms is mapping of platform ID to Platform.
var Platforms = map[string]Platform{}

// Platform definitions.
var (
	LinuxAmd64 = Platform{
		ID:           "linux/amd64",
		Arch:         "x86_64",
		Target:       "x86_64-sentinelos-linux",
		Build:        "x86_64-linux",
		Host:         "x86_64-linux",
		LLBPlatform:  llb.LinuxAmd64,
		PlatformSpec: platforms.MustParse("linux/amd64"),
	}

	LinuxArm64 = Platform{
		ID:           "linux/arm64",
		Arch:         "aarch64",
		Target:       "aarch64-sentinelos-linux",
		Build:        "aarch64-linux",
		Host:         "aarch64-linux",
		LLBPlatform:  llb.LinuxArm64,
		PlatformSpec: platforms.MustParse("linux/arm64"),
	}
)

func init() {
	for _, platform := range []Platform{
		LinuxAmd64,
		LinuxArm64,
	} {
		Platforms[platform.ID] = platform
	}
}

// Platform describes build & target platforms.
type Platform struct {
	ID           string
	Arch         string
	Target       string
	Build        string
	Host         string
	LLBPlatform  llb.ConstraintsOpt
	PlatformSpec specs.Platform
}

// BuildVariables returns build env variables.
func (p Platform) BuildVariables() types.Variables {
	return types.Variables{
		"BUILD": p.Build,
		"HOST":  p.Host,
	}
}

// TargetVariables returns target env variables.
func (p Platform) TargetVariables() types.Variables {
	return types.Variables{
		"ARCH":   p.Arch,
		"TARGET": p.Target,
	}
}

func (p Platform) String() string {
	return p.ID
}

// Set implements pflag.Value interface.
func (p *Platform) Set(id string) error {
	if _, exists := Platforms[id]; !exists {
		return fmt.Errorf("platform %q is not defined", id)
	}

	*p = Platforms[id]

	return nil
}

// Type implements pflag.Value interface.
func (p *Platform) Type() string {
	return "platform"
}
