package v1alpha1

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// Dependency on another image or stage.
type Dependency struct {
	Image   string `yaml:"image,omitempty"`
	Stage   string `yaml:"stage,omitempty"`
	From    string `yaml:"from,omitempty"`
	To      string `yaml:"to,omitempty"`
	Runtime bool   `yaml:"runtime,omitempty"`
}

// IsInternal checks whether dependency is internal to some stage.
func (d *Dependency) IsInternal() bool {
	return d.Stage != ""
}

// Src returns copy source (from dependency).
func (d *Dependency) Src() string {
	if d.From != "" {
		return d.From
	}

	return "/"
}

// Dest returns copy destination (to base).
func (d *Dependency) Dest() string {
	if d.To != "" {
		return d.To
	}

	return "/"
}

// Validate the dependency.
func (d *Dependency) Validate() error {
	if d.Image != "" && d.Stage != "" {
		return fmt.Errorf("dependency can't have both image & stage set: %q, %q", d.Image, d.Stage)
	}

	if d.Image == "" && d.Stage == "" {
		return fmt.Errorf("either image or stage should be set for the dependency")
	}

	return nil
}

// Dependencies is a list of Depency.
type Dependencies []Dependency

// Validate dependencies.
func (deps Dependencies) Validate() error {
	var multiErr *multierror.Error

	for _, dep := range deps {
		multiErr = multierror.Append(multiErr, dep.Validate())
	}

	return multiErr.ErrorOrNil()
}
