package testutil

import (
	"os"
	"os/exec"
	"sync"
	"testing"
)

// BuildkitRunner runs packer via buildctl/buildkit.
type BuildkitRunner struct {
	CommandRunner
	Target   string
	Platform string
}

// Run implements Run interface.
func (runner BuildkitRunner) Run(t *testing.T) {
	if err := IsBuildkitAvailable(); err != nil {
		t.Skipf("buildkit is not available: %q", err)
	}

	args := append(getBuildkitGlobalFlags(),
		"build",
		"--frontend", "dockerfile.v0",
		"--local", "context=.",
		"--local", "dockerfile=.",
		"--opt", "filename=Pkgfile",
		"--opt", "target="+runner.Target,
		"--build-arg", "TAG=testtag",
	)

	if runner.Platform != "" {
		args = append(args, "--opt", "platform="+runner.Platform)
	}

	cmd := exec.Command("buildctl", args...)

	runner.run(t, cmd, "buildkit")
}

func getBuildkitGlobalFlags() []string {
	var globalOpts []string

	if buildkitHost, ok := os.LookupEnv("BUILDKIT_HOST"); ok {
		globalOpts = append(globalOpts, "--addr", buildkitHost)
	}

	return globalOpts
}

var (
	buildkitCheckOnce sync.Once
	//nolint:errname
	buildkitCheckError error
)

// IsBuildkitAvailable returns nil if buildkit is ready to use.
func IsBuildkitAvailable() error {
	buildkitCheckOnce.Do(func() {
		buildkitCheckError = exec.Command("buildctl", append(getBuildkitGlobalFlags(), "debug", "workers")...).Run()
	})

	return buildkitCheckError
}
