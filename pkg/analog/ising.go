package analog

// IsingModel represents the Ising Hamiltonian for quantum computation.
type IsingModel struct {
	Spins         int                // Number of spins (qubits) in the model
	Interactions  map[[2]int]float64 // J_ij values representing interaction strengths between spins
	ExternalField map[int]float64    // h_i values representing external fields on spins
}
