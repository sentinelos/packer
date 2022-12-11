// Package cmd contains definitions of CLI commands.
package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/sentinelos/packer/pkg/environment"
)

const defaultPlatform = runtime.GOOS + "/" + runtime.GOARCH

var (
	pkgRoot string
	debug   bool
	options = &environment.Options{
		BuildPlatform:  environment.LinuxAmd64,
		TargetPlatform: environment.LinuxAmd64,
	}
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "packer",
	Short: "A tool to build and manage software via Pkgfile and pkg.yaml",
	Long: `packer usually works in buildkit frontend mode when it's not directly
exposed as a CLI tool. In that mode of operation packer loads root Pkgfile and
a set of pkg.yamls, processes them, builds dependency graph and outputs it
as LLB graph to buildkit backend.

packer can be also used to produce graph of dependencies between build steps and
output LLB directly which is useful for development or debugging.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
//
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Enable debug logging")
	rootCmd.PersistentFlags().StringVarP(&pkgRoot, "root", "", ".", "The path to a pkg root")

	options.BuildPlatform.Set(defaultPlatform)  //nolint:errcheck
	options.TargetPlatform.Set(defaultPlatform) //nolint:errcheck
}
