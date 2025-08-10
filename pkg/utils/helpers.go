package utils

import (q
	"encoding/json"
)

// Waveform represents a waveform in the Pulser sequence.
type Waveform struct {
	Type     string  `json:"type"`
	Duration int     `json:"duration"`
	Start    float32 `json:"start,omitempty"`
	Stop     float32 `json:"stop,omitempty"`
	Area     float32 `json:"area,omitempty"` // For Blackman, area is like integral, but here we use amplitude as per example
}

// Pulse represents a pulse in the sequence.
type Pulse struct {
	Channel   string   `json:"channel"`
	Amplitude Waveform `json:"amplitude"`
	Detuning  Waveform `json:"detuning"`
	Phase     float64  `json:"phase"`
}

// Channel represents a declared channel.
type Channel struct {
	Type string `json:"type"`
}

// Register represents the qubit register.
type Register struct {
	Qubits struct {
		Positions map[string][]float64 `json:"positions"` // e.g., "q0": [0, 0]
	} `json:"qubits"`
}

// SequenceBuilder is the core structure for the Pulser sequence abstract representation.
type SequenceBuilder struct {
	Channels map[string]Channel `json:"channels"`
	Pulses   []Pulse            `json:"pulses"`
	Register Register           `json:"register"`
	// Additional fields can be added as needed, e.g., variables, measurement, etc.
}

// InputData wraps the sequence builder for submission (as per Azure/Pasqal format).
type InputData struct {
	SequenceBuilder SequenceBuilder `json:"sequence_builder"`
}

// CreateSimplePulserSequence generates a simple Pulser sequence in JSON format.
// This example creates a sequence with:
// - A single qubit register at position (0,0).
// - A Rydberg global channel "ch0".
// - A single pulse with Blackman amplitude and Ramp detuning.
func CreateSimplePulserSequence() (string, error) {
	// Define the register
	reg := Register{}
	reg.Qubits.Positions = map[string][]float64{
		"q0": {0.0, 0.0},
	}

	// Define channels
	channels := map[string]Channel{
		"ch0": {Type: "rydberg_global"},
	}

	// Define a simple pulse
	pulse := Pulse{
		Channel: "ch0",
		Amplitude: Waveform{
			Type:     "blackman",
			Duration: 1000,
			Area:     3.141592653589793, // Approximate pi for amplitude area
		},
		Detuning: Waveform{
			Type:     "ramp",
			Duration: 1000,
			Start:    -5.0,
			Stop:     5.0,
		},
		Phase: 0.0,
	}

	// Build the sequence
	builder := SequenceBuilder{
		Channels: channels,
		Pulses:   []Pulse{pulse},
		Register: reg,
	}

	inputData := InputData{SequenceBuilder: builder}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(inputData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
