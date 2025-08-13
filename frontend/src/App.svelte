<script lang="ts">
  import { onMount } from 'svelte';

  // Type definitions
  interface Edge {
    from: number;
    to: number;
  }

  // State variables
  let numVertices: number = 3; // Default number of vertices
  let numColors: number = 2; // Default number of colors (k)
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

  // Function to generate Rydberg Hamiltonian JSON for k-Coloring
  function generateJSON() {
    // For k-Coloring, each vertex-color pair (v, c) is a qubit. We need to:
    // 1. Ensure each vertex is assigned exactly one color (local fields).
    // 2. Penalize adjacent vertices having the same color (interactions).
    const h: { [key: string]: number } = {};
    const J: { [key: string]: number } = {};

    // For each vertex, encourage exactly one color (via penalty for deviation)
    for (let v = 0; v < numVertices; v++) {
      for (let c = 0; c < numColors; c++) {
        const qubit = `${v}_${c}`; // Qubit for vertex v, color c
        h[qubit] = -1; // Encourage assigning a color
      }
    }

    // Penalize adjacent vertices having the same color
    edges.forEach(edge => {
      for (let c = 0; c < numColors; c++) {
        const key = `${edge.from}_${c},${edge.to}_${c}`;
        J[key] = 1; // Positive coupling to penalize same color on adjacent vertices
      }
    });

    // Optionally, add constraints to ensure each vertex has exactly one color
    // This can be approximated in the Hamiltonian with penalty terms
    for (let v = 0; v < numVertices; v++) {
      for (let c1 = 0; c1 < numColors; c1++) {
        for (let c2 = c1 + 1; c2 < numColors; c2++) {
          const key = `${v}_${c1},${v}_${c2}`;
          J[key] = 1; // Penalize assigning multiple colors to the same vertex
        }
      }
    }

    const rydbergHamiltonian = { h, J };
    jsonOutput = JSON.stringify(rydbergHamiltonian, null, 2);
  }

  // Simulate solving process
  async function calculateColoring() {
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

    // Simulated k-Coloring solution (greedy heuristic for demo purposes)
    const colors: number[] = new Array(numVertices).fill(-1); // -1 means uncolored
    const availableColors = Array.from({ length: numColors }, (_, i) => i);
    const vertices = Array.from({ length: numVertices }, (_, i) => i).sort(() => Math.random() - 0.5);

    for (const v of vertices) {
      const neighborColors = edges
        .filter(e => e.from === v || e.to === v)
        .map(e => e.from === v ? colors[e.to] : colors[e.from])
        .filter(c => c !== -1);
      const validColor = availableColors.find(c => !neighborColors.includes(c));
      if (validColor !== undefined) {
        colors[v] = validColor;
      }
    }

    const isValid = colors.every(c => c !== -1) && edges.every(e => colors[e.from] !== colors[e.to]);
    if (isValid) {
      result = `Valid ${numColors}-Coloring: ${colors.map((c, v) => `Vertex ${v}: Color ${c}`).join(', ')}`;
    } else {
      result = `No valid ${numColors}-coloring found for the given graph.`;
    }

    solving = false;
  }

  // Generate initial JSON
  onMount(() => {
    generateJSON();
  });
</script>

<main>
  <h1>k-Coloring Problem Input</h1>
  
  <section>
    <h2>Graph Parameters</h2>
    <label>
      Number of Vertices:
      <input type="number" bind:value={numVertices} min="1" on:change={generateJSON} />
    </label>
    <label>
      Number of Colors (k):
      <input type="number" bind:value={numColors} min="1" on:change={generateJSON} />
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
    <button on:click={calculateColoring} disabled={solving}>Calculate</button>
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
    background-color: #666; /* Grayed-out button */
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