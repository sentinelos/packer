package testutil

import (
	"os/exec"
	"sync"
	"testing"
)

// DockerRunner runs packer via docker buildx.
type DockerRunner struct {
	CommandRunner
	Target   string
	Platform string
}

// Run implements Run interface.
func (runner DockerRunner) Run(t *testing.T) {
	if err := IsDockerAvailable(); err != nil {
		t.Skipf("docker buildx is not available: %q", err)
	}

	args := []string{
		"buildx",
		"build",
		"-f", "./Pkgfile",
		"--target", runner.Target,
		"--build-arg", "TAG=testtag",
	}

	if runner.Platform != "" {
		args = append(args, "--platform", runner.Platform)
	}

	cmd := exec.Command("docker", append(args, ".")...)

	runner.run(t, cmd, "docker buildx")
}

var (
	dockerCheckOnce sync.Once
	//nolint:errname
	dockerCheckError error
)

// IsDockerAvailable returns nil if docker buildx is ready to use.
func IsDockerAvailable() error {
	dockerCheckOnce.Do(func() {
		dockerCheckError = exec.Command("docker", "buildx", "ls").Run()
	})

	return dockerCheckError
}
