package v1alpha1

// Finalize is a set of COPY instructions to finalize the build.
type Finalize struct {
	From   string            `yaml:"from,omitempty"`
	To     string            `yaml:"to,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
}
