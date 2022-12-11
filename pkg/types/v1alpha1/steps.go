package v1alpha1

import (
	"errors"

	"github.com/hashicorp/go-multierror"
)

// Steps is a collection of Step.
type Steps []Step

// Validate steps.
func (steps Steps) Validate() error {
	var multiErr *multierror.Error

	if len(steps) == 0 {
		multiErr = multierror.Append(multiErr, errors.New("steps are missing, this is going to lead to empty build"))
	}

	for _, step := range steps {
		multiErr = multierror.Append(multiErr, step.Validate())
	}

	return multiErr.ErrorOrNil()
}

// Step describes a single build step.
//
// Steps are executed sequentially, each step runs in its own
// empty temporary directory.
type Step struct {
	Env       Environment  `yaml:"env,omitempty"`
	CachePath string       `yaml:"cache,omitempty"`
	TmpDir    string       `yaml:"-"`
	Sources   Sources      `yaml:"sources,omitempty"`
	Prepare   Instructions `yaml:"prepare,omitempty"`
	Build     Instructions `yaml:"build,omitempty"`
	Install   Instructions `yaml:"install,omitempty"`
	Test      Instructions `yaml:"test,omitempty"`
}

// Validate the step.
func (step *Step) Validate() error {
	return step.Sources.Validate()
}
