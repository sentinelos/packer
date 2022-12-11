package solver

import (
	"github.com/sentinelos/packer/pkg/types/v1alpha1"
)

// LoadResult is a result of PackageLoader.Load function.
type LoadResult struct {
	Pkgfile *v1alpha1.Pkgfile
	Pkgs    []*v1alpha1.Pkg
}

// PackageLoader implements some way to fetch collection of Pkgs.
type PackageLoader interface {
	Load() (*LoadResult, error)
}
