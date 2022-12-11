package v1alpha1

import (
	"bytes"
	"errors"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"

	"github.com/sentinelos/packer/pkg/constants"
	"github.com/sentinelos/packer/pkg/types"
)

// Pkg represents build instructions for a single package.
type Pkg struct {
	Context       types.Variables `yaml:"-"`
	Title         string          `yaml:"title"`
	Description   string          `yaml:"description,omitempty"`
	Licenses      string          `yaml:"licenses,omitempty"`
	Authors       string          `yaml:"authors,omitempty"`
	Documentation string          `yaml:"documentation,omitempty"`
	Source        string          `yaml:"source,omitempty"`
	Variant       Variant         `yaml:"variant,omitempty"`
	Shell         Shell           `yaml:"shell,omitempty"`
	BaseDir       string          `yaml:"-"`
	FileName      string          `yaml:"-"`
	Install       Install         `yaml:"install,omitempty"`
	Dependencies  Dependencies    `yaml:"dependencies,omitempty"`
	Steps         Steps           `yaml:"steps"`
	Finalize      Finalize        `yaml:"finalize,omitempty"`
}

// NewPkg loads Pkg structure from file.
func NewPkg(baseDir, fileName string, contents []byte, vars types.Variables) (*Pkg, error) {
	p := &Pkg{
		BaseDir:  baseDir,
		FileName: fileName,
		Shell:    "/bin/sh",
		Variant:  Alpine,
		Context:  vars.Copy(),
		Finalize: Finalize{
			From: constants.Artifacts,
			To:   "/",
		},
	}

	tmpl, err := template.New(constants.PkgYaml).
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, vars); err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(&buf).Decode(p); err != nil {
		return nil, err
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate the Pkg.
func (p *Pkg) Validate() error {
	var multiErr *multierror.Error

	if p.Title == "" {
		multiErr = multierror.Append(multiErr, errors.New("package title can't be empty"))
	}

	multiErr = multierror.Append(multiErr, p.Dependencies.Validate(), p.Steps.Validate())

	return multiErr.ErrorOrNil()
}
