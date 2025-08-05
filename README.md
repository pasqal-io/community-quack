Analog Neutral Atom Quantum Compiler in Golang

This repository contains a Go-based compiler for analog neutral atom quantum computers. The compiler is designed to solve various optimization problems by mapping them to the Ising Hamiltonian, a fundamental model in statistical mechanics. This approach leverages the unique properties of neutral atoms and their interactions to find the ground state of the Hamiltonian, which corresponds to the solution of the original problem.

The Physics: Analog Neutral Atom Quantum Computing

Analog quantum computing is a paradigm where the quantum system's evolution is directly analogous to the problem being solved. In our case, we use a system of neutral atoms to simulate the behavior of the Ising model.

Neutral Atoms as Qubits

In this quantum computing architecture, individual neutral atoms, such as Rubidium-87, are used as qubits. These atoms are cooled to ultra-low temperatures and trapped in a two-dimensional array using a grid of focused laser beams called optical tweezers.

The two quantum states of the qubit, |0⟩ and |1⟩, are represented by two different energy levels of the atom. Typically, these are a ground state and a highly excited Rydberg state.

The Rydberg Blockade
The key to performing computations with neutral atoms is the Rydberg blockade. When an atom is excited to a Rydberg state, its electron is in a high-energy orbit, far from the nucleus. This makes the atom much larger and highly interactive with its neighbors.

If two atoms are close to each other (within a certain "blockade radius"), the energy required to excite the second atom to a Rydberg state is significantly shifted due to the strong interaction with the first. This prevents the second atom from being excited by the same laser pulse.

This phenomenon is the foundation for creating entanglement and performing conditional logic, which are essential for quantum computation.

The Mathematics: The Ising Hamiltonian
The Ising model is a mathematical model of ferromagnetism in statistical mechanics. It describes a collection of "spins" that can be in one of two states (+1 or -1) and interact with their neighbors. The energy of a particular configuration of spins is given by the Ising Hamiltonian:

Mapping Problems to the Ising Hamiltonian
A wide range of NP-hard optimization problems can be mapped to the problem of finding the ground state of an Ising Hamiltonian. This is done by carefully constructing the J_ij and h_i terms such that the lowest energy configuration of the Ising model corresponds to the optimal solution of the original problem.

Example: Maximum Independent Set (MIS)
A classic example is the Maximum Independent Set problem. An independent set in a graph is a set of vertices where no two vertices are connected by an edge. The MIS problem is to find the largest such set.

To map this to an Ising model, we can represent each vertex in the graph with a spin. Let's define the state of a spin as follows:

We can then construct the Hamiltonian with the following constraints:

Maximize the size of the set: We want to encourage as many spins as possible to be +1. This can be achieved with a negative external field (h < 0) that lowers the energy for spins in the +1 state.

Enforce the independence constraint: If two vertices i and j are connected by an edge, they cannot both be in the independent set. We can enforce this with a large, positive interaction strength (J_ij > 0) for all connected vertices. This penalizes configurations where connected spins are both +1.

By finding the ground state of this constructed Hamiltonian, we can identify the spin configuration that corresponds to the Maximum Independent Set of the original graph.

The Golang Implementation
This project provides a Go library to facilitate the process of compiling optimization problems into the Ising Hamiltonian format suitable for analog neutral atom quantum computers.

Key Data Structures
The core of the library revolves around these main data structures:

// Represents a graph for problems like Maximum Independent Set.
type Graph struct {
    Vertices int
    Edges    [][2]int
}

// Represents the Ising Hamiltonian.
type IsingModel struct {
    Spins         int
    Interactions  map[[2]int]float64 // J_ij values
    ExternalField map[int]float64      // h_i values
}

The Compiler's Role
The primary function of the compiler is to take a problem-specific input (like a Graph) and generate the corresponding IsingModel.

// CompileMIS takes a graph and returns the corresponding IsingModel for the
// Maximum Independent Set problem.
func CompileMIS(graph Graph, h float64, J float64) IsingModel {
    // ... implementation ...
}

Usage
Here is a simple example of how to use the library to compile a Maximum Independent Set problem:

package main

import (
    "fmt"
    "github.com/your-username/analog-quantum-compiler-go"
)

func main() {
    // Define a simple graph (a square).
    g := analog.Graph{
        Vertices: 4,
        Edges: [][2]int{
            {0, 1},
            {1, 2},
            {2, 3},
            {3, 0},
        },
    }

    // Compile the graph to an Ising model.
    isingModel := analog.CompileMIS(g, -1.0, 2.0)

    // The isingModel can now be sent to an analog quantum computer
    // for solving.
    fmt.Printf("Ising Model: %+v\n", isingModel)
}

Getting Started
To get started with this project, clone the repository and install the dependencies:

git clone https://github.com/your-username/analog-quantum-compiler-go.git
cd analog-quantum-compiler-go
go mod tidy

Contributing
Contributions are welcome! Please feel free to submit a pull request or open an issue.

License
This project is licensed under the MIT License. See the LICENSE file for details.
