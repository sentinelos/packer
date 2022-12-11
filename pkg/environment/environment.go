// Package environment contains default environment contents.
package environment

import (
	"fmt"

	"github.com/sentinelos/packer/pkg/constants"
	"github.com/sentinelos/packer/pkg/types"
)

// Default returns default values for environment variables.
func Default() types.Variables {
	return types.Variables{
		"CFLAGS":    constants.CFLAGS,
		"CXXFLAGS":  constants.CXXFLAGS,
		"LDFLAGS":   constants.LDFLAGS,
		"VENDOR":    constants.Vendor,
		"ARTIFACTS": constants.Artifacts,
		"TOOLCHAIN": constants.Toolchain,
		"PATH":      fmt.Sprintf("%s/bin:%s", constants.Toolchain, constants.DefaultPath),
	}
}
