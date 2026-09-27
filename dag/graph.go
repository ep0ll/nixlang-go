package dag

import (
	"fmt"
	"sort"
)

// NodeID uniquely identifies a node in the graph (typically a store path
// string or derivation name).
type NodeID string

// NodeKind classifies a graph node.
type NodeKind int

const (
	KindUnknown NodeKind = iota
	KindDerivation
	KindOutput
	KindInputSrc
)

func (k NodeKind) String() string {
	switch k {
	case KindDerivation:
		return "derivation"
	case KindOutput:
		return "output"
	case KindInputSrc:
		return "input-src"
	default:
		return "unknown"
	}
}

// Node is a single vertex in the derivation graph.
type Node struct {
	ID   NodeID
	Kind NodeKind
	Name string
	// Meta holds optional attributes (system, output names, etc.).
	Meta map[string]string
}

// Edge is a directed dependency from From → To
// (From depends on To / To is an input of From).
type Edge struct {
	From NodeID
	To   NodeID
}

// Graph is a DAG of store/derivation nodes.
type Graph struct {
	nodes map[NodeID]*Node
	edges []Edge
}

// New returns an empty graph.
func New() *Graph {
	return &Graph{nodes: make(map[NodeID]*Node)}
}

// AddNode inserts or replaces a node. ID must be non-empty.
func (g *Graph) AddNode(n Node) error {
	if g == nil {
		return fmt.Errorf("dag: nil graph")
	}
	if n.ID == "" {
		return fmt.Errorf("dag: empty node ID")
	}
	cp := n
	if n.Meta != nil {
		cp.Meta = make(map[string]string, len(n.Meta))
		for k, v := range n.Meta {
			cp.Meta[k] = v
		}
	}
	g.nodes[n.ID] = &cp
	return nil
}

// AddEdge records a dependency edge. Both endpoints should already exist;
// missing endpoints are allowed and can be validated later via Validate.
func (g *Graph) AddEdge(from, to NodeID) error {
	if g == nil {
		return fmt.Errorf("dag: nil graph")
	}
	if from == "" || to == "" {
		return fmt.Errorf("dag: empty edge endpoint")
	}
	if from == to {
		return fmt.Errorf("dag: self-edge on %s", from)
	}
	g.edges = append(g.edges, Edge{From: from, To: to})
	return nil
}

// Node returns a node by ID, or nil.
func (g *Graph) Node(id NodeID) *Node {
	if g == nil {
		return nil
	}
	return g.nodes[id]
}

// Nodes returns all nodes sorted by ID.
func (g *Graph) Nodes() []*Node {
	if g == nil {
		return nil
	}
	ids := make([]NodeID, 0, len(g.nodes))
	for id := range g.nodes {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	out := make([]*Node, len(ids))
	for i, id := range ids {
		out[i] = g.nodes[id]
	}
	return out
}

// Edges returns a copy of the edge list.
func (g *Graph) Edges() []Edge {
	if g == nil {
		return nil
	}
	out := make([]Edge, len(g.edges))
	copy(out, g.edges)
	return out
}

// Validate checks that every edge endpoint exists and that the graph is acyclic.
func (g *Graph) Validate() error {
	if g == nil {
		return fmt.Errorf("dag: nil graph")
	}
	for _, e := range g.edges {
		if _, ok := g.nodes[e.From]; !ok {
			return fmt.Errorf("dag: edge from unknown node %q", e.From)
		}
		if _, ok := g.nodes[e.To]; !ok {
			return fmt.Errorf("dag: edge to unknown node %q", e.To)
		}
	}
	// Kahn's algorithm for cycle detection.
	in := make(map[NodeID]int, len(g.nodes))
	for id := range g.nodes {
		in[id] = 0
	}
	children := make(map[NodeID][]NodeID, len(g.nodes))
	for _, e := range g.edges {
		in[e.From]++
		children[e.To] = append(children[e.To], e.From)
	}
	var q []NodeID
	for id, c := range in {
		if c == 0 {
			q = append(q, id)
		}
	}
	seen := 0
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		seen++
		for _, ch := range children[n] {
			in[ch]--
			if in[ch] == 0 {
				q = append(q, ch)
			}
		}
	}
	if seen != len(g.nodes) {
		return fmt.Errorf("dag: cycle detected (%d/%d nodes reachable in topo order)", seen, len(g.nodes))
	}
	return nil
}

// Len returns the number of nodes.
func (g *Graph) Len() int {
	if g == nil {
		return 0
	}
	return len(g.nodes)
}
