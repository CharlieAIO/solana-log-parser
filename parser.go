// Package solana-log-parser provides utilities for parsing Solana transaction logs
// and organizing them by instruction and inner instruction indices.
package solanalogparser

import (
	"github.com/CharlieAIO/solana-log-parser/internal/parser"
	"github.com/CharlieAIO/solana-log-parser/internal/structs"
)

// ParseTransactionLogs parses a slice of Solana transaction log strings and returns
// a structured representation organized by instruction index and inner instruction index.
//
// The returned map structure is:
//   - First level key: instruction index (int8)
//   - Second level key: inner instruction index (int8, -1 for top-level instructions)
//   - Value: slice of LogEntry structs for that instruction/inner instruction
//
// Example usage:
//
//	logs := []string{
//		"Program 11111111111111111111111111111111 invoke [1]",
//		"Program 22222222222222222222222222222222 invoke [2]",
//		"Program 22222222222222222222222222222222 success",
//		"Program 11111111111111111111111111111111 success",
//	}
//	result := ParseTransactionLogs(logs)
func ParseTransactionLogs(logs []string) map[int8]map[int8][]structs.LogEntry {
	return parser.ParseTransactionLogs(logs)
}

// LogEntry represents a single log entry from a Solana transaction
type LogEntry = structs.LogEntry

// InstructionLogs represents logs for a specific instruction
type InstructionLogs = structs.InstructionLogs
