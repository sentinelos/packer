package testutil

import (
	"os"

	"gopkg.in/yaml.v2"
)

// TestManifest describes single integration test in test.yaml.
type TestManifest struct {
	Runs []RunManifest `yaml:"run"`
}

// RunManifest describes single run of integration test.
type RunManifest struct {
	Name         string  `yaml:"name"`
	Runner       string  `yaml:"runner"`
	Platform     string  `yaml:"platform"`
	Target       string  `yaml:"target"`
	Expect       string  `yaml:"expect"`
	ExpectStdout *string `yaml:"expectStdout"`
	CreateFile   string  `yaml:"createFile"`
	Template     string  `yaml:"template"`
}

// NewTestManifest loads TestManifest from test.yaml file.
func NewTestManifest(path string) (manifest TestManifest, err error) {
	var f *os.File

	f, err = os.Open(path)
	if err != nil {
		return
	}

	defer f.Close() //nolint:errcheck

	err = yaml.NewDecoder(f).Decode(&manifest)

	return
}
