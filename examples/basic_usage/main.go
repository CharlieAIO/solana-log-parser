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
