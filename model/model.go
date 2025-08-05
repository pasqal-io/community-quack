package model

// Graph represents a graph for problems like Maximum Independent Set.
type Graph struct {
	Vertices int      // Number of vertices in the graph
	Edges    [][2]int // List of edges (pairs of vertex indices)
}

// IsingModel represents the Ising Hamiltonian.
type IsingModel struct {
	Spins         int                // Number of spins (vertices)
	Interactions  map[[2]int]float64 // J_ij: Interaction strengths between spins
	ExternalField map[int]float64    // h_i: External field for each spin
}
