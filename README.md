# Solana Log Parser

A Go package for parsing and organizing Solana transaction logs by instruction and inner instruction indices.

## Features

- **Structured Parsing**: Organizes logs by instruction index and inner instruction index
- **Program Tracking**: Tracks program invocations and their success/failure status
- **Call Stack Management**: Maintains proper call stack context for nested program calls
- **Type Safety**: Full Go type safety with exported structs and functions
- **Zero Dependencies**: Pure Go implementation with no external dependencies

## Installation

```bash
go get github.com/CharlieAIO/solana-log-parser
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/CharlieAIO/solana-log-parser"
)

func main() {
	// Example Solana transaction logs
	rawLogs := []string{
		"Program ComputeBudget111111111111111111111111111111 invoke [1]",
		"Program ComputeBudget111111111111111111111111111111 success",
		"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]",
		"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
	}

	// Parse the logs
	result := solanalogparser.ParseTransactionLogs(rawLogs)

	instructionIndex := int8(0)
	innerInstructionIndex := int8(-1) // -1 indicates not an inner instruction

	for _, logs := range result[instructionIndex][innerInstructionIndex] {
		fmt.Println(logs)
	}

}
```

## API Reference

### Functions

#### `ParseTransactionLogs(logs []string) map[int8]map[int8][]LogEntry`

Parses a slice of Solana transaction log strings and returns a structured representation.

**Parameters:**
- `logs`: Slice of log strings from a Solana transaction

**Returns:**
- A nested map where:
  - First level key: instruction index (int8)
  - Second level key: inner instruction index (int8, -1 for top-level instructions)
  - Value: slice of LogEntry structs for that instruction/inner instruction

### Types

#### `LogEntry`

```go
type LogEntry struct {
    Content          string `json:"content"`           // The log content
    ProgramID        string `json:"program_id,omitempty"` // Program ID if applicable
    IsInvoke         bool   `json:"is_invoke"`         // True if this is a program invoke
    IsSuccess        bool   `json:"is_success"`        // True if this is a program success
    InnerInstruction int8   `json:"inner_instruction"` // -1 if not an inner instruction
}
```

#### `InstructionLogs`

```go
type InstructionLogs struct {
    InstructionIndex int8       `json:"instruction_index"` // The instruction index
    Logs             []LogEntry `json:"logs"`              // All logs for this instruction
}
```

## Examples

See the `examples/` directory for complete working examples:

- **Basic Usage** (`examples/basic_usage/`): Simple parsing and JSON output

## How It Works

The parser processes Solana transaction logs sequentially, maintaining a call stack to track program invocations. It identifies:

1. **Program Invokes**: Lines matching `Program <ID> invoke [<depth>]`
2. **Program Successes**: Lines matching `Program <ID> success`
3. **Regular Logs**: All other lines, categorized by current call context

The parser handles nested program calls (inner instructions) by tracking the call depth and organizing logs accordingly.

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built for the Solana blockchain ecosystem
- Designed to work with standard Solana transaction log formats
- Inspired by the need for better log analysis tools in DeFi development
