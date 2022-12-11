package cmd

import (
	"context"
	"log"

	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"

	"github.com/sentinelos/packer/pkg/pkgfile"
)

// frontendCmd represents the frontend command.
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Buildkit frontend for Pkgfile",
	Long: `This command implements buildkit frontend.

To activate, put following line as the first line of Pkgfile:

# syntax = ghcr.io/sentinelos/packer:<version>

Run with:

  nerdctl build -f ./Pkgfile --target <target> .
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := grpcclient.RunFromEnvironment(
			appcontext.Context(),
			func(ctx context.Context, c client.Client) (*client.Result, error) {
				return pkgfile.Build(ctx, c, options)
			},
		); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(frontendCmd)
}
