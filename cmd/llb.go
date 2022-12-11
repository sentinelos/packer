package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/moby/buildkit/client/llb"
	solverpb "github.com/moby/buildkit/solver/pb"
	"github.com/spf13/cobra"

	"github.com/sentinelos/packer/pkg/convert"
	"github.com/sentinelos/packer/pkg/solver"
)

var llbCmdFlags struct {
	json bool
}

// llbCmd represents the llb command.
var llbCmd = &cobra.Command{
	Use:   "llb",
	Short: "Dump buildkit LLB for the build",
	Long: `This command parses build instructions from pkg.yaml files,
and outputs buildkit LLB to stdout. This can be used as 'packer pack ... | buildctl ...'.
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
		}

		packages, err := solver.NewPackages(&loader)
		if err != nil {
			log.Fatal(err)
		}

		graph, err := packages.Resolve(options.Target)
		if err != nil {
			log.Fatal(err)
		}

		dt, err := convert.MarshalLLB(graph, options)
		if err != nil {
			log.Fatal(err)
		}

		if llbCmdFlags.json {
			pb := dt.ToPB()
			var b []byte
			if b, err = json.MarshalIndent(pb, "", "  "); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", b)

			for _, def := range pb.Def {
				b, err = json.MarshalIndent(def, "", "  ")
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Def %s: ", b)

				op := new(solverpb.Op)
				if err = op.Unmarshal(def); err != nil {
					log.Fatal(err)
				}

				b, err = json.MarshalIndent(op, "", "  ")
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s\n", b)
			}

			return
		}

		err = llb.WriteTo(dt, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	llbCmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target image to build")
	llbCmd.MarkFlagRequired("target") //nolint:errcheck
	llbCmd.Flags().Var(&options.BuildPlatform, "build-platform", "Build platform")
	llbCmd.Flags().Var(&options.TargetPlatform, "target-platform", "Target platform")
	llbCmd.Flags().BoolVar(&llbCmdFlags.json, "json", false, "Dump as JSON for debug")
	rootCmd.AddCommand(llbCmd)
}
