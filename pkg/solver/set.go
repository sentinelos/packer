package solver

import (
	"io"

	"github.com/emicklei/dot"
)

// PackageSet is a list of PackageNodes.
type PackageSet []*PackageNode

// DumpDot dumps nodes and deps in dot format.
func (set PackageSet) DumpDot(w io.Writer) {
	g := dot.NewGraph(dot.Directed)

	for _, node := range set {
		node.DumpDot(g)
	}

	g.Write(w)
}
