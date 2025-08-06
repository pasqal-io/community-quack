package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"quack/pkg/analog"
)

// GraphInput represents the JSON structure for graph input.
type GraphInput struct {
	Vertices int      `json:"vertices"`
	Edges    [][2]int `json:"edges"`
}

// main is the entry point for the CLI tool of the analog quantum compiler.
func main() {
	// Define command-line flags
	inputFile := flag.String("input", "", "Path to the JSON input file containing the graph (required)")
	h := flag.Float64("h", -1.0, "External field for the Ising model (typically negative)")
	j := flag.Float64("j", 2.0, "Interaction strength for the Ising model (typically positive)")
	jsonOutput := flag.Bool("json", false, "Output the Ising model in JSON format")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This tool compiles a graph into an Ising model for analog neutral atom quantum computers.\n")
		fmt.Fprintf(os.Stderr, "Example: %s -input=graph.json -h=-1.0 -j=2.0\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Expected JSON format for input file:\n")
		fmt.Fprintf(os.Stderr, `{"vertices": 4, "edges": [[0,1], [1,2], [2,3], [3,0]]}`+"\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Validate input file flag
	if *inputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: Input file must be specified")
		flag.Usage()
		os.Exit(1)
	}

	// Validate h and j parameters
	if *h >= 0 {
		fmt.Fprintln(os.Stderr, "Warning: h is typically negative for MIS problems")
	}
	if *j <= 0 {
		fmt.Fprintln(os.Stderr, "Error: j must be positive")
		flag.Usage()
		os.Exit(1)
	}

	// Read and parse the JSON input file
	graph, err := readGraphFromJSON(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading or parsing JSON file: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Validate graph
	if err := validateGraph(graph); err != nil {
		fmt.Fprintf(os.Stderr, "Error in graph data: %v\n", err)
		os.Exit(1)
	}

	// Compile the graph into an Ising model
	isingModel := analog.CompileMIS(graph, *h, *j)

	// Output the Ising model
	if *jsonOutput {
		output, err := json.MarshalIndent(isingModel, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating JSON output: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(string(output))
	} else {
		fmt.Printf("Ising Model: %+v\n", isingModel)
	}
}

// readGraphFromJSON reads and parses the graph from a JSON file.
func readGraphFromJSON(filePath string) (analog.Graph, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return analog.Graph{}, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	var graphInput GraphInput
	if err := json.Unmarshal(fileContent, &graphInput); err != nil {
		return analog.Graph{}, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return analog.Graph{
		Vertices: graphInput.Vertices,
		Edges:    graphInput.Edges,
	}, nil
}

// validateGraph checks the validity of the graph data.
func validateGraph(graph analog.Graph) error {
	if graph.Vertices <= 0 {
		return fmt.Errorf("number of vertices must be greater than 0")
	}
	for _, edge := range graph.Edges {
		if edge[0] < 0 || edge[0] >= graph.Vertices || edge[1] < 0 || edge[1] >= graph.Vertices {
			return fmt.Errorf("edge indices must be between 0 and %d: %v", graph.Vertices-1, edge)
		}
		if edge[0] == edge[1] {
			return fmt.Errorf("loops are not allowed: %v", edge)
		}
	}
	return nil
}
