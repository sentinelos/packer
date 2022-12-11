// Package update provides facilities for checking for available pkgs updates.
package update

import (
	"context"
	"fmt"
	"net/url"
)

// LatestInfo represents information about available update.
type LatestInfo struct {
	// BaseURL may contain base URL for releases.
	BaseURL string
	// LatestURL may contain URL for the latest asset.
	LatestURL string
	// HasUpdate is true if there seems to be an update available.
	HasUpdate bool
}

// Latest returns information about available update.
func Latest(ctx context.Context, source string) (*LatestInfo, error) {
	u, err := url.Parse(source)
	if err != nil {
		return nil, err
	}

	switch u.Host {
	case "github.com":
		return newGitHub(gitHubTokenFromEnv()).Latest(ctx, source) //nolint:contextcheck

	default:
		return nil, fmt.Errorf("unhandled host %q", u.Host)
	}
}
