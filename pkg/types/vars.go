// Package types describes basic types which are not versioned.
package types

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v2"
)

// Variables presents generic variables for templating/environment.
type Variables map[string]string

// Merge two Variables in place.
func (v Variables) Merge(other Variables) Variables {
	for key := range other {
		v[key] = other[key]
	}

	return v
}

// Copy the Variables.
func (v Variables) Copy() Variables {
	result := make(Variables, len(v))

	for key, val := range v {
		result[key] = val
	}

	return result
}

// Load the variables from YAML with the given context.
func (v *Variables) Load(path string, context Variables) error {
	tmpl, err := template.New(filepath.Base(path)).
		Funcs(sprig.HermeticTxtFuncMap()).
		ParseFiles(path)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, context); err != nil {
		return err
	}

	return yaml.Unmarshal(buf.Bytes(), v)
}

// LoadContents the variables from byte slice with the given context.
func (v *Variables) LoadContents(contents []byte, context Variables) error {
	tmpl, err := template.New("vars").
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(string(contents))
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, context); err != nil {
		return err
	}

	return yaml.Unmarshal(buf.Bytes(), v)
}
