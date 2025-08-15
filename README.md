## Analog Neutral Atom Compiler

A cross-platform desktop application for compiling Ising Hamiltonians into Pasqal pulser sequences. This tool provides a user-friendly interface for researchers and developers working with neutral atom quantum computers, simplifying the process of translating a theoretical model into a hardware-compatible pulse sequence.


### Ising Hamiltonian Input

Define your quantum system using the familiar Ising model format.

#### Pasqal Pulse Generation 

Automatically compiles the model into a valid sequence for Pasqal's pulser library.

#### High-Performance Backend

Core compilation logic is written in Go for speed and efficiency.

#### Modern GUI
A sleek and responsive user interface built with Wails and Svelte.

#### Cross-Platform
Build and run the application on Windows, macOS and Linux.

### Architecture

This application leverages the power of Wails to create a bridge between a Go backend and a Svelte web-based frontend.

### Frontend (Svelte): 

The frontend/ directory contains the Svelte application. It is responsible for rendering the UI components, managing user input, and displaying the results. It makes calls to the Go backend for the core compilation logic.

### Backend (Go): 

The main.go file and other .go source files contain the application's core logic. This includes:

Parsing and validating the input Ising Hamiltonian.

Implementing the compilation algorithm to map Hamiltonian coefficients to pulse parameters (amplitude, detuning, duration).

Exposing methods that can be called directly from the Svelte frontend.

Wails Bridge: Wails handles the communication between the frontend and backend, allowing you to call Go functions from JavaScript and vice-versa. It also manages window creation, menus, and bundling the application into a single executable.

### Getting Started

Follow these instructions to get a local copy up and running.

### Prerequisites

You must have the following tools installed on your system:

Go (version 1.18 or newer)

Node.js (LTS version)

Wails CLI: Follow the official installation guide at wails.io/docs/gettingstarted/installation.

### Installation

Clone the repository:

git clone https://github.com/pasqal-io/community-quack
cd community-quack

Build the application:
The Wails CLI will handle all dependencies and compile the application into a single executable for your platform.

wails build

The executable will be located in the build/bin/ directory.

### How to Use

#### Launch the application.

Enter your Ising Hamiltonian parameters into the input text area on the left. See the required format below.

Click the "Compile" button.

The generated Pasqal pulser sequence will appear in the output panel on the right. You can use the "Copy" button to copy it to your clipboard.


## For Developers

If you want to contribute to the development, you can run the application in live development mode. This provides hot-reloading for both the Go backend and the Svelte frontend.

Navigate to the project directory:

cd analog-neutral-atom-compiler

Run the development server:

wails dev

The application will launch, and any changes you make to the Go or Svelte source files will be automatically reloaded.

## Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

Fork the Project

Create your Feature Branch (git checkout -b feature/AmazingFeature)

Commit your Changes (git commit -m 'Add some AmazingFeature')

Push to the Branch (git push origin feature/AmazingFeature)

Open a Pull Request

## License
Distributed under the MIT License. See LICENSE for more information.
