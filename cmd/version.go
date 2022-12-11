package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sentinelos/packer/pkg/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints packer version.",
	Long:  `Prints packer version.`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, _ []string) {
		line := fmt.Sprintf("%s version %s (%s)", version.Name, version.Tag, version.SHA)
		fmt.Println(line)
	},
}
