//go:build integration

package integration_test

import (
	"testing"

	"github.com/sentinelos/packer/pkg/util/testutil"
)

func TestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in -short mode")
	}

	collection, err := testutil.CollectTests()
	if err != nil {
		t.Fatalf("error collecting tests: %v", err)
	}

	collection.Each(func(name string, f func(*testing.T)) {
		t.Run(name, f)
	})
}
