package dag_test

import (
	"testing"

	"github.com/ep0ll/nixlang-go/dag"
)

func TestGraphAcyclic(t *testing.T) {
	g := dag.New()
	_ = g.AddNode(dag.Node{ID: "a", Kind: dag.KindDerivation, Name: "a"})
	_ = g.AddNode(dag.Node{ID: "b", Kind: dag.KindDerivation, Name: "b"})
	_ = g.AddNode(dag.Node{ID: "c", Kind: dag.KindOutput, Name: "c"})
	_ = g.AddEdge("a", "b") // a depends on b
	_ = g.AddEdge("b", "c")
	if err := g.Validate(); err != nil {
		t.Fatal(err)
	}
	if g.Len() != 3 {
		t.Fatalf("Len = %d", g.Len())
	}
	nodes := g.Nodes()
	if len(nodes) != 3 || nodes[0].ID != "a" {
		t.Fatalf("Nodes = %+v", nodes)
	}
}

func TestGraphCycle(t *testing.T) {
	g := dag.New()
	_ = g.AddNode(dag.Node{ID: "a"})
	_ = g.AddNode(dag.Node{ID: "b"})
	_ = g.AddEdge("a", "b")
	_ = g.AddEdge("b", "a")
	if err := g.Validate(); err == nil {
		t.Fatal("expected cycle error")
	}
}

func TestGraphMissingEndpoint(t *testing.T) {
	g := dag.New()
	_ = g.AddNode(dag.Node{ID: "a"})
	_ = g.AddEdge("a", "missing")
	if err := g.Validate(); err == nil {
		t.Fatal("expected missing endpoint error")
	}
}

func TestGraphRejectsSelfEdge(t *testing.T) {
	g := dag.New()
	_ = g.AddNode(dag.Node{ID: "a"})
	if err := g.AddEdge("a", "a"); err == nil {
		t.Fatal("expected self-edge error")
	}
}

func TestGraphNilSafe(t *testing.T) {
	var g *dag.Graph
	if g.Len() != 0 {
		t.Fatal("nil Len")
	}
	if err := g.AddNode(dag.Node{ID: "x"}); err == nil {
		t.Fatal("expected error on nil graph")
	}
	if err := g.Validate(); err == nil {
		t.Fatal("expected error on nil Validate")
	}
}
