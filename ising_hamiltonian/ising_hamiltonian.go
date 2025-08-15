package isinghamiltonian

import (
	"fmt"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/mat"
)

func Hamiltonian(n int, h float64) {
	powN := int(math.Pow(2, float64(n)))
	qc := make([]*QuantumCircuit, 2*n-1)

	// Creating the quantum circuits that are used in the calculation of the Hamiltonian based on the number of qubits
	for i := 0; i < 2*n-1; i++ {
		qr := NewQuantumRegister(n)
		qc[i] = NewQuantumCircuit(qr) // create quantum circuits for each factor of the Hamiltonian
		if i <= n-2 {                 // for the first sum of the Hamiltonian
			qc[i].Z(i)     // value of current spin
			qc[i].Z(i + 1) // and value of its neighboring spin
		} else { // for the second sum of the Hamiltonian
			qc[i].X(2*n - 2 - i) // 2*n-2 gives the proper index since counting starts at 0
		}
	}

	// Run each circuit in the simulator
	simulator := NewUnitarySimulator()
	result := make([]*Result, 2*n-1)
	unitary := make([]*mat.Dense, 2*n-1)
	HamiltonianMatrix := mat.NewDense(powN, powN, nil)

	// Get the results for each circuit in unitary form
	for i := 0; i < 2*n-1; i++ {
		result[i] = simulator.Execute(qc[i])
		unitary[i] = result[i].GetUnitary()

		// And calculate the Hamiltonian matrix according to the formula
		if i <= n-2 {
			HamiltonianMatrix.Add(HamiltonianMatrix, unitary[i].Scale(-1, unitary[i]))
		} else {
			HamiltonianMatrix.Add(HamiltonianMatrix, unitary[i].Scale(-h, unitary[i]))
		}
	}

	fmt.Printf("The %d x %d Hamiltonian Matrix is:\n", powN, powN)
	matPrint(HamiltonianMatrix)

	// Now that we have the Hamiltonian
	var eigenvalues []complex128
	var eigenvectors *mat.Dense
	eigenvalues, eigenvectors = mat.Eigen(HamiltonianMatrix)

	fmt.Println("Eigenvectors")
	matPrint(eigenvectors)
	fmt.Println("Eigenvalues")
	for _, val := range eigenvalues {
		fmt.Println(val)
	}

	minimum := real(eigenvalues[0])
	minSpot := 0
	for i := 1; i < powN; i++ {
		if real(eigenvalues[i]) < minimum {
			minSpot = i
			minimum = real(eigenvalues[i])
		}
	}
	fmt.Println(minSpot)

	groundstate := mat.Col(nil, minSpot, eigenvectors)
	probability := make([]float64, powN)
	for i := 0; i < powN; i++ {
		probability[i] = real(groundstate[i] * cmplx.Conj(groundstate[i]))
	}

	fmt.Printf("The probability for each of the %d base states is:\n", powN)
	for _, prob := range probability {
		fmt.Println(prob)
	}
	fmt.Printf("The probabilities for each of the %d base states add up to: %.2f\n", powN, sum(probability))
}

func matPrint(X mat.Matrix) {
	fc := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fc)
}

func sum(slice []float64) float64 {
	total := 0.0
	for _, v := range slice {
		total += v
	}
	return total
}
