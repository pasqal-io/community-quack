<script lang="ts">
  import { onMount } from 'svelte';

  // Type definitions for the graph
  interface Edge {
    from: number;
    to: number;
    weight: number;
  }

  // State variables
  let numNodes: number = 3; // Default number of nodes
  let edges: Edge[] = []; // List of edges
  let newFrom: number = 0;
  let newTo: number = 1;
  let newWeight: number = 1.0;
  let jsonOutput: string = ''; // Generated JSON output

  // Function to add an edge
  function addEdge() {
    if (newFrom >= numNodes || newTo >= numNodes || newFrom === newTo || newFrom < 0 || newTo < 0) {
      alert('Invalid nodes! Nodes must be between 0 and ' + (numNodes - 1) + ' and different.');
      return;
    }
    // Check if edge already exists
    const exists = edges.some(e => 
      (e.from === newFrom && e.to === newTo) || (e.from === newTo && e.to === newFrom)
    );
    if (exists) {
      alert('Edge already exists!');
      return;
    }
    edges = [...edges, { from: newFrom, to: newTo, weight: newWeight }];
  }

  // Function to remove an edge
  function removeEdge(index: number) {
    edges = edges.filter((_, i) => i !== index);
  }

  // Function to generate the Ising Hamiltonian JSON
  function generateJSON() {
    const h: { [key: string]: number } = {};
    for (let i = 0; i < numNodes; i++) {
      h[i.toString()] = 0; // No external fields for pure Max-Cut
    }

    const J: { [key: string]: number } = {};
    edges.forEach(edge => {
      // For undirected graphs: smaller ID first for consistency
      const key = edge.from < edge.to ? `${edge.from},${edge.to}` : `${edge.to},${edge.from}`;
      J[key] = -edge.weight; // Negative sign for Max-Cut (since Ising minimizes H = -sum J sigma_i sigma_j)
    });

    const isingHamiltonian = { h, J };
    jsonOutput = JSON.stringify(isingHamiltonian, null, 2);
  }

  // Generate initial JSON
  onMount(() => {
    generateJSON();
  });
</script>

<main>
  <h1>Max-Cut Problem Input</h1>
  
  <section>
    <h2>Graph Parameters</h2>
    <label>
      Number of Nodes:
      <input type="number" bind:value={numNodes} min="2" on:change={generateJSON} />
    </label>
    
    <h3>Add Edges</h3>
    <label>
      From Node (0 to {numNodes - 1}):
      <input type="number" bind:value={newFrom} min="0" max={numNodes - 1} />
    </label>
    <label>
      To Node (0 to {numNodes - 1}):
      <input type="number" bind:value={newTo} min="0" max={numNodes - 1} />
    </label>
    <label>
      Weight:
      <input type="number" bind:value={newWeight} step="0.1" />
    </label>
    <button on:click={addEdge}>Add Edge</button>
    
    <h3>Current Edges</h3>
    <ul>
      {#each edges as edge, index}
        <li>
          {edge.from} -- {edge.to} (Weight: {edge.weight})
          <button on:click={() => removeEdge(index)}>Remove</button>
        </li>
      {/each}
    </ul>
    
    <button on:click={generateJSON}>Generate Ising Hamiltonian</button>
  </section>
  
  <section>
    <h2>Generated Ising Hamiltonian JSON (for Analog Neutral Atom Compiler)</h2>
    <pre>{jsonOutput}</pre>
  </section>
</main>

<style>
  main {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
  }
  label {
    display: block;
    margin: 10px 0;
  }
  button {
    margin: 10px 0;
  }
  pre {
    background: #f4f4f4;
    padding: 10px;
    overflow: auto;
  }
</style>