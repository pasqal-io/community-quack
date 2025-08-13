<script lang="ts">
  import { onMount } from 'svelte';

  // Type definitions
  interface Edge {
    from: number;
    to: number;
  }

  // State variables
  let numVertices: number = 3; // Default number of vertices
  let edges: Edge[] = []; // List of edges
  let newFrom: number = 0;
  let newTo: number = 1;
  let jsonOutput: string = ''; // Generated JSON output
  let solving: boolean = false;
  let progress: number = 0;
  let progressText: string = '';
  let result: string = '';

  // Function to add an edge
  function addEdge() {
    if (newFrom >= numVertices || newTo >= numVertices || newFrom === newTo || newFrom < 0 || newTo < 0) {
      alert(`Invalid vertices! Vertices must be between 0 and ${numVertices - 1} and different.`);
      return;
    }
    const exists = edges.some(e => 
      (e.from === newFrom && e.to === newTo) || (e.from === newTo && e.to === newFrom)
    );
    if (exists) {
      alert('Edge already exists!');
      return;
    }
    edges = [...edges, { from: newFrom, to: newTo }];
    generateJSON();
  }

  // Function to remove an edge
  function removeEdge(index: number) {
    edges = edges.filter((_, i) => i !== index);
    generateJSON();
  }

  // Function to generate Rydberg Hamiltonian JSON for MIS
  function generateJSON() {
    // For MIS, we map to a Rydberg Hamiltonian where selected vertices are in the Rydberg state (|1⟩).
    // The blockade effect ensures adjacent vertices are not both selected.
    const h: { [key: string]: number } = {};
    for (let i = 0; i < numVertices; i++) {
      h[i.toString()] = -1; // Local detuning to favor selecting vertices (maximize independent set)
    }

    const J: { [key: string]: number } = {};
    edges.forEach(edge => {
      // For undirected graphs: smaller ID first for consistency
      const key = edge.from < edge.to ? `${edge.from},${edge.to}` : `${edge.to},${edge.from}`;
      J[key] = 1; // Positive coupling to penalize adjacent vertices being selected
    });

    const rydbergHamiltonian = { h, J };
    jsonOutput = JSON.stringify(rydbergHamiltonian, null, 2);
  }

  // Simulate solving process
  async function calculateMIS() {
    solving = true;
    progress = 0;
    result = '';

    const stages = [
      { text: 'Compiling', duration: 1000 },
      { text: 'Calling Pascal API', duration: 1500 },
      { text: 'Calculating', duration: 2000 },
      { text: 'Result', duration: 500 }
    ];

    for (const stage of stages) {
      progressText = stage.text;
      progress = (stages.indexOf(stage) + 1) / stages.length * 100;
      await new Promise(resolve => setTimeout(resolve, stage.duration));
    }

    // Simulated MIS solution (heuristic for demo purposes)
    const independentSet: number[] = [];
    const vertices = Array.from({ length: numVertices }, (_, i) => i);
    const shuffledVertices = vertices.sort(() => Math.random() - 0.5); // Randomize for simplicity
    for (const v of shuffledVertices) {
      if (!independentSet.some(u => edges.some(e => 
        (e.from === u && e.to === v) || (e.from === v && e.to === u)
      ))) {
        independentSet.push(v);
      }
    }
    result = `Maximum Independent Set: {${independentSet.join(', ')}} (Size: ${independentSet.length})`;

    solving = false;
  }

  // Generate initial JSON
  onMount(() => {
    generateJSON();
  });
</script>

<main>
  <h1>Maximum Independent Set Problem Input</h1>
  
  <section>
    <h2>Graph Parameters</h2>
    <label>
      Number of Vertices:
      <input type="number" bind:value={numVertices} min="1" on:change={generateJSON} />
    </label>
    
    <h3>Add Edges</h3>
    <label>
      From Vertex (0 to {numVertices - 1}):
      <input type="number" bind:value={newFrom} min="0" max={numVertices - 1} />
    </label>
    <label>
      To Vertex (0 to {numVertices - 1}):
      <input type="number" bind:value={newTo} min="0" max={numVertices - 1} />
    </label>
    <button on:click={addEdge}>Add Edge</button>
    
    <h3>Current Edges</h3>
    <ul>
      {#each edges as edge, index}
        <li>
          {edge.from} -- {edge.to}
          <button on:click={() => removeEdge(index)}>Remove</button>
        </li>
      {/each}
    </ul>
  </section>

  <section>
    <button on:click={calculateMIS} disabled={solving}>Calculate</button>
    {#if solving}
      <div class="progress-container">
        <div class="progress-bar" style="width: {progress}%"></div>
      </div>
      <p>{progressText}</p>
    {/if}
    {#if result}
      <h3>Result</h3>
      <p>{result}</p>
    {/if}
  </section>
</main>

<style>
  main {
    max-width: 100%;
    margin: 0 auto;
    padding: 20px;
    font-family: Arial, sans-serif;
    background-color: #1a1a1a; /* Dark background */
    color: #f0f0f0; /* Light text for contrast */
  }
  h1, h2, h3 {
    color: #ffffff; /* White headings for visibility */
  }
  label {
    display: block;
    margin: 10px 0;
    font-weight: bold;
    color: #e0e0e0; /* Light gray for labels */
  }
  input {
    padding: 8px;
    margin: 5px 0;
    border: 1px solid #555; /* Darker border for contrast */
    border-radius: 4px;
    width: 100px;
    background-color: #333; /* Dark input background */
    color: #f0f0f0; /* Light text in inputs */
  }
  button {
    margin: 10px 0;
    padding: 10px 20px;
    background-color: #28a745; /* Green button */
    color: #ffffff; /* White text */
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 16px;
  }
  button:disabled {
    background-color: #111; /* Grayed-out button */
    cursor: not-allowed;
  }
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    margin: 5px 0;
    padding: 5px;
    background: #333; /* Dark list item background */
    border-radius: 4px;
    color: #f0f0f0; /* Light text */
  }
  li button {
    margin-left: 10px;
    background-color: #dc3545; /* Red remove button */
    padding: 5px 10px;
    font-size: 14px;
  }
  pre {
    background: #2a2a2a; /* Darker code background */
    color: #f0f0f0; /* Light code text */
    padding: 15px;
    border-radius: 4px;
    overflow: auto;
    max-height: 300px;
  }
  .progress-container {
    width: 100%;
    background-color: #333; /* Dark progress container */
    border-radius: 4px;
    margin: 10px 0;
  }
  .progress-bar {
    height: 20px;
    background-color: #28a745; /* Green progress bar */
    border-radius: 4px;
    transition: width 0.3s ease-in-out;
  }
  p {
    color: #e0e0e0; /* Light gray for progress text and result */
  }
</style>