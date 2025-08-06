package analog

// Graph represents a graph structure with vertices and edges for optimization problems.
type Graph struct {
	Vertices int      // Number of vertices in the graph
	Edges    [][2]int // List of edges, where each edge is a pair of vertex indices
}
