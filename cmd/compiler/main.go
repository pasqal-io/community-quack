package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"quack/pkg/analog"
	"sync"
)

// GraphInput represents the JSON structure for graph input.
type GraphInput struct {
	Vertices int      `json:"vertices"`
	Edges    [][2]int `json:"edges"`
}

// ProcessingResult represents the result of processing a file
type ProcessingResult struct {
	FilePath   string
	IsingModel analog.IsingModel
	Error      error
}

// main is the entry point for the CLI tool of the analog quantum compiler.
func main() {
	// Define command-line flags
	var inputFiles []string
	flag.Func("input", "Path to the JSON input file containing the graph (can be specified multiple times)", func(s string) error {
		inputFiles = append(inputFiles, s)
		return nil
	})
	h := flag.Float64("h", -1.0, "External field for the Ising model (typically negative)")
	j := flag.Float64("j", 2.0, "Interaction strength for the Ising model (typically positive)")
	jsonOutput := flag.Bool("json", false, "Output the Ising model in JSON format")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "This tool compiles graphs into Ising models for analog neutral atom quantum computers.\n")
		fmt.Fprintf(os.Stderr, "Example: %s -input=graph1.json -input=graph2.json -h=-1.0 -j=2.0\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Expected JSON format for input file:\n")
		fmt.Fprintf(os.Stderr, `{"vertices": 4, "edges": [[0,1], [1,2], [2,3], [3,0]]}`+"\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Validate input file flag
	if len(inputFiles) == 0 {
		fmt.Fprintln(os.Stderr, "Error: At least one input file must be specified")
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

	// Channels for pipeline
	fileChan := make(chan string, len(inputFiles))
	graphChan := make(chan analog.Graph, len(inputFiles))
	resultChan := make(chan ProcessingResult, len(inputFiles))

	// WaitGroup for synchronization
	var wg sync.WaitGroup
	wg.Add(3) // One for each stage: read, validate, compile

	// Start goroutines
	go readFiles(inputFiles, fileChan, &wg)
	go processGraphs(fileChan, graphChan, resultChan, &wg)
	go compileGraphs(graphChan, resultChan, *h, *j, &wg)

	// Collect results in a separate goroutine
	go func() {
		for result := range resultChan {
			if result.Error != nil {
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", result.FilePath, result.Error)
				continue
			}
			if *jsonOutput {
				output, err := json.MarshalIndent(result.IsingModel, "", "  ")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error generating JSON output for %s: %v\n", result.FilePath, err)
					continue
				}
				fmt.Printf("Result for %s:\n%s\n", result.FilePath, string(output))
			} else {
				fmt.Printf("Result for %s:\nIsing Model: %+v\n", result.FilePath, result.IsingModel)
			}
		}
	}()

	// Wait for all goroutines to complete
	wg.Wait()
	close(resultChan) // Close result channel after all processing is done

	fmt.Println("Compilation completed.")
}

// readFiles sends input file paths to the file channel
func readFiles(inputFiles []string, fileChan chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, filePath := range inputFiles {
		fileChan <- filePath
	}
	close(fileChan)
}

// processGraphs reads and validates graphs
func processGraphs(fileChan <-chan string, graphChan chan<- analog.Graph, resultChan chan<- ProcessingResult, wg *sync.WaitGroup) {
	defer wg.Done()
	for filePath := range fileChan {
		graph, err := readGraphFromJSON(filePath)
		if err != nil {
			resultChan <- ProcessingResult{FilePath: filePath, Error: err}
			continue
		}
		if err := validateGraph(graph); err != nil {
			resultChan <- ProcessingResult{FilePath: filePath, Error: err}
			continue
		}
		graphChan <- graph
	}
	close(graphChan)
}

// compileGraphs compiles graphs into Ising models
func compileGraphs(graphChan <-chan analog.Graph, resultChan chan<- ProcessingResult, h, j float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for graph := range graphChan {
		isingModel := analog.CompileMIS(graph, h, j)
		resultChan <- ProcessingResult{FilePath: fmt.Sprintf("graph_%d", graph.Vertices), IsingModel: isingModel}
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
