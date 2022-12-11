package convert

import (
	"context"
	"sort"

	"github.com/moby/buildkit/client/llb"

	"github.com/sentinelos/packer/pkg/constants"
	"github.com/sentinelos/packer/pkg/environment"
	"github.com/sentinelos/packer/pkg/solver"
	"github.com/sentinelos/packer/pkg/types/v1alpha1"
)

// GraphLLB wraps PackageGraph to provide LLB conversion.
//
// GraphLLB caches common images used in the build.
type GraphLLB struct {
	*solver.PackageGraph

	Options *environment.Options

	BaseImages   map[v1alpha1.Variant]llb.State
	Checksummer  llb.State
	LocalContext llb.State

	baseImageProcessor llbProcessor
	cache              map[*solver.PackageNode]llb.State

	commonRunOptions []llb.RunOption
}

type llbProcessor func(llb.State) llb.State

// NewGraphLLB creates new GraphLLB and initializes shared images.
func NewGraphLLB(graph *solver.PackageGraph, options *environment.Options) *GraphLLB {
	result := &GraphLLB{
		PackageGraph: graph,
		Options:      options,
		cache:        make(map[*solver.PackageNode]llb.State),
	}

	if options.ProxyEnv != nil {
		result.commonRunOptions = append(result.commonRunOptions, llb.WithProxy(*options.ProxyEnv))
	}

	result.buildBaseImages()
	result.buildChecksummer()
	result.buildLocalContext()

	return result
}

func (graph *GraphLLB) buildBaseImages() {
	graph.BaseImages = make(map[v1alpha1.Variant]llb.State)

	addPkg := func(root llb.State) llb.State {
		return root.File(
			llb.Mkdir(pkgDir, constants.DefaultDirMode),
			llb.WithCustomNamef("%smkdir %s", graph.Options.CommonPrefix, pkgDir),
		).Dir(pkgDir)
	}

	addEnv := func(root llb.State) llb.State {
		vars := graph.Options.GetVariables()
		keys := make([]string, 0, len(vars))

		for key := range vars {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			root = root.AddEnv(key, vars[key])
		}

		return root
	}

	graph.baseImageProcessor = func(root llb.State) llb.State {
		return addEnv(addPkg(root))
	}

	graph.BaseImages[v1alpha1.Alpine] = graph.baseImageProcessor(llb.Image(
		constants.DefaultBaseImage,
		llb.WithCustomName(graph.Options.CommonPrefix+"base"),
	).Run(
		append(graph.commonRunOptions,
			llb.Shlex("apk --no-cache --update add bash"),
			llb.WithCustomName(graph.Options.CommonPrefix+"base-apkinstall"),
		)...,
	).Run(
		append(graph.commonRunOptions,
			llb.Args([]string{"ln", "-svf", "/bin/bash", "/bin/sh"}),
			llb.WithCustomName(graph.Options.CommonPrefix+"base-symlink"),
		)...,
	).Root())

	graph.BaseImages[v1alpha1.Scratch] = graph.baseImageProcessor(llb.Scratch())
}

func (graph *GraphLLB) buildChecksummer() {
	graph.Checksummer = llb.Image(
		constants.DefaultBaseImage,
		llb.WithCustomName(graph.Options.CommonPrefix+"cksum"),
	).Run(
		append(graph.commonRunOptions,
			llb.Shlex("apk --no-cache --update add coreutils"),
			llb.WithCustomName(graph.Options.CommonPrefix+"cksum-apkinstall"),
		)...,
	).Root()
}

func (graph *GraphLLB) buildLocalContext() {
	graph.LocalContext = llb.Local(
		"context",
		llb.ExcludePatterns(
			[]string{
				"**/.*",
				"**/" + constants.PkgYaml,
				"**/" + constants.VarsYaml,
			},
		),
		llb.WithCustomName(graph.Options.CommonPrefix+"context"),
	)
}

// Build converts package graph to LLB.
func (graph *GraphLLB) Build() (llb.State, error) {
	return NewNodeLLB(graph.Root, graph).Build()
}

// Marshal returns marshaled LLB.
func (graph *GraphLLB) Marshal() (*llb.Definition, error) {
	out, err := graph.Build()
	if err != nil {
		return nil, err
	}

	out = out.SetMarshalDefaults(graph.Options.BuildPlatform.LLBPlatform)

	return out.Marshal(context.TODO())
}
