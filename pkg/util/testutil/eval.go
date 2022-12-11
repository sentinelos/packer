package testutil

import (
	"os/exec"
	"testing"
)

// EvalRunner runs packer eval.
type EvalRunner struct {
	CommandRunner

	Target   string
	Template string
}

// Run implements Run interface.
func (runner EvalRunner) Run(t *testing.T) {
	cmd := exec.Command("packer", "eval", "--target", runner.Target, runner.Template)

	runner.run(t, cmd, "packer eval")
}
