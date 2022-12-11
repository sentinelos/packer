package v1alpha1

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/sentinelos/packer/pkg/types"
)

// Pkgfile describes structure of 'Pkgfile'.
type Pkgfile struct {
	Vars   types.Variables `yaml:"vars,omitempty"`
	Format string          `yaml:"format"`
}

// NewPkgfile loads Pkgfile from `[]byte` contents.
func NewPkgfile(contents []byte) (*Pkgfile, error) {
	var pkgfile Pkgfile

	if err := yaml.Unmarshal(contents, &pkgfile); err != nil {
		return nil, err
	}

	// TODO: this might be used in the future to pick correct format based on Pkgfile, leave it simple for now
	if pkgfile.Format != "v1alpha1" {
		return nil, fmt.Errorf("unsupported format: %q, supported formats: %q", pkgfile.Format, []string{"v1alpha1"})
	}

	return &pkgfile, nil
}
