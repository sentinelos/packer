// Package constants provides basic constants for the program
package constants

import (
	"os"
)

// Set of variables set during the build.
var (
	DefaultRegistry     string
	DefaultOrganization string
	Version             string
)

const (
	// DefaultBaseImage for non-scratch builds.
	// renovate: datasource=docker versioning=docker depName=alpine
	DefaultBaseImage = "docker.io/alpine:3.17"

	// DefaultDirMode is UNIX file mode for mkdir.
	DefaultDirMode os.FileMode = 0o755

	// DefaultPath is default value for PATH environment variable.
	DefaultPath = "/bin:/sbin:/usr/bin:/usr/local/bin:/usr/local/sbin:/usr/sbin"

	// PkgYaml is the filename of 'pkg.yaml'.
	PkgYaml = "pkg.yaml"

	// VarsYaml is the filename of 'vars.yaml'.
	VarsYaml = "vars.yaml"

	// Pkgfile is the filename of 'Pkgfile'.
	Pkgfile = "Pkgfile"

	CFLAGS    = "-g0 -Os"
	CXXFLAGS  = "-g0 -Os"
	LDFLAGS   = "-s"
	Vendor    = "sentinelos"
	Artifacts = "/artifacts"
	Toolchain = "/toolchain"
)
