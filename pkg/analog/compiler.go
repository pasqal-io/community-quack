package analog

// CompileMIS compiles a graph into an Ising model for the Maximum Independent Set problem.
// The parameter h is the external field (typically negative to encourage spins in +1 state),
// and j is the interaction strength (typically positive to penalize adjacent spins both in +1 state).
func CompileMIS(graph Graph, h, j float64) IsingModel {
	// Initialize the Ising model
	ising := IsingModel{
		Spins:         graph.Vertices,
		Interactions:  make(map[[2]int]float64),
		ExternalField: make(map[int]float64),
	}

	// Set external field h_i for each spin to encourage +1 state (maximize set size)
	for i := 0; i < graph.Vertices; i++ {
		ising.ExternalField[i] = h
	}

	// Set interaction strengths J_ij for edges to penalize adjacent spins both in +1 state
	for _, edge := range graph.Edges {
		// Ensure edge is stored in a consistent order (smaller index first)
		var key [2]int
		if edge[0] < edge[1] {
			key = [2]int{edge[0], edge[1]}
		} else {
			key = [2]int{edge[1], edge[0]}
		}
		// Check for duplicate edges
		if _, exists := ising.Interactions[key]; exists {
			return IsingModel{} // Return empty model to indicate error
		}
		ising.Interactions[key] = j
	}

	return ising
}
