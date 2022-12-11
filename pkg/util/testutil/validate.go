package testutil

import (
	"os/exec"
	"testing"
)

// ValidateRunner runs packer validate.
type ValidateRunner struct {
	CommandRunner
}

// Run implements Run interface.
func (runner ValidateRunner) Run(t *testing.T) {
	cmd := exec.Command("packer", "validate")

	runner.run(t, cmd, "packer validate")
}
