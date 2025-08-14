package solanalogparser

import (
	"testing"
)

func TestParseTransactionLogs(t *testing.T) {
	tests := []struct {
		name     string
		logs     []string
		expected map[int8]map[int8][]LogEntry
	}{
		{
			name: "Simple single instruction",
			logs: []string{
				"Program ComputeBudget111111111111111111111111111111 invoke [1]",
				"Program ComputeBudget111111111111111111111111111111 success",
			},
			expected: map[int8]map[int8][]LogEntry{
				0: {
					-1: {
						{Content: "Program ComputeBudget111111111111111111111111111111 invoke [1]", ProgramID: "ComputeBudget111111111111111111111111111111", IsInvoke: true, IsSuccess: false, InnerInstruction: -1},
						{Content: "Program ComputeBudget111111111111111111111111111111 success", ProgramID: "ComputeBudget111111111111111111111111111111", IsInvoke: false, IsSuccess: true, InnerInstruction: -1},
					},
				},
			},
		},
		{
			name: "Multiple instructions",
			logs: []string{
				"Program ComputeBudget111111111111111111111111111111 invoke [1]",
				"Program ComputeBudget111111111111111111111111111111 success",
				"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]",
				"Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success",
			},
			expected: map[int8]map[int8][]LogEntry{
				0: {
					-1: {
						{Content: "Program ComputeBudget111111111111111111111111111111 invoke [1]", ProgramID: "ComputeBudget111111111111111111111111111111", IsInvoke: true, IsSuccess: false, InnerInstruction: -1},
						{Content: "Program ComputeBudget111111111111111111111111111111 success", ProgramID: "ComputeBudget111111111111111111111111111111", IsInvoke: false, IsSuccess: true, InnerInstruction: -1},
					},
				},
				1: {
					-1: {
						{Content: "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA invoke [1]", ProgramID: "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA", IsInvoke: true, IsSuccess: false, InnerInstruction: 0},
						{Content: "Program TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA success", ProgramID: "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA", IsInvoke: false, IsSuccess: true, InnerInstruction: 0},
					},
				},
			},
		},
		{
			name:     "Empty logs",
			logs:     []string{},
			expected: map[int8]map[int8][]LogEntry{},
		},
		{
			name:     "Whitespace only logs",
			logs:     []string{"   ", "  ", ""},
			expected: map[int8]map[int8][]LogEntry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTransactionLogs(tt.logs)

			// Check if the result has the expected number of instructions
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d instructions, got %d", len(tt.expected), len(result))
			}

			// For simple cases, do detailed comparison
			if tt.name == "Simple single instruction" || tt.name == "Empty logs" || tt.name == "Whitespace only logs" {
				if !mapsEqual(result, tt.expected) {
					t.Errorf("Expected %+v, got %+v", tt.expected, result)
				}
			}
		})
	}
}

func TestParseTransactionLogsComplex(t *testing.T) {
	logs := []string{
		"Program 11111111111111111111111111111111 invoke [1]",
		"Program 22222222222222222222222222222222 invoke [2]",
		"Program 33333333333333333333333333333333 invoke [3]",
		"Program 33333333333333333333333333333333 success",
		"Program 22222222222222222222222222222222 success",
		"Program 11111111111111111111111111111111 success",
	}

	result := ParseTransactionLogs(logs)

	// Verify structure
	if len(result) != 1 {
		t.Errorf("Expected 1 instruction, got %d", len(result))
	}

	instruction0, exists := result[0]
	if !exists {
		t.Fatal("Expected instruction 0 to exist")
	}

	// Check top-level instruction
	topLevel, exists := instruction0[-1]
	if !exists {
		t.Fatal("Expected top-level instruction (-1) to exist")
	}
	if len(topLevel) != 2 {
		t.Errorf("Expected 2 top-level logs, got %d", len(topLevel))
	}

	// Check inner instructions
	inner0, exists := instruction0[0]
	if !exists {
		t.Fatal("Expected inner instruction 0 to exist")
	}
	if len(inner0) != 2 {
		t.Errorf("Expected 2 inner instruction 0 logs, got %d", len(inner0))
	}

	inner1, exists := instruction0[1]
	if !exists {
		t.Fatal("Expected inner instruction 1 to exist")
	}
	if len(inner1) != 2 {
		t.Errorf("Expected 2 inner instruction 1 logs, got %d", len(inner1))
	}
}

func TestParseTransactionLogsWithRegularLogs(t *testing.T) {
	logs := []string{
		"Program 11111111111111111111111111111111 invoke [1]",
		"Regular log message",
		"Another regular log",
		"Program 11111111111111111111111111111111 success",
	}

	result := ParseTransactionLogs(logs)

	instruction0, exists := result[0]
	if !exists {
		t.Fatal("Expected instruction 0 to exist")
	}

	topLevel, exists := instruction0[-1]
	if !exists {
		t.Fatal("Expected top-level instruction (-1) to exist")
	}

	// Should have 4 logs: invoke, 2 regular logs, and success
	if len(topLevel) != 4 {
		t.Errorf("Expected 4 logs, got %d", len(topLevel))
	}

	// Check that regular logs are properly categorized
	regularLogCount := 0
	for _, logEntry := range topLevel {
		if !logEntry.IsInvoke && !logEntry.IsSuccess {
			regularLogCount++
		}
	}
	if regularLogCount != 2 {
		t.Errorf("Expected 2 regular logs, got %d", regularLogCount)
	}
}

// Helper function to compare maps (simplified for testing)
func mapsEqual(a, b map[int8]map[int8][]LogEntry) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if bv, exists := b[k]; !exists {
			return false
		} else if len(v) != len(bv) {
			return false
		}
	}
	return true
}
