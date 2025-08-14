package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/CharlieAIO/solana-log-parser/internal/structs"
)

func ParseTransactionLogs(logs []string) map[int8]map[int8][]structs.LogEntry {
	result := make(map[int8]map[int8][]structs.LogEntry)

	invokeRegex := regexp.MustCompile(`^Program (.+) invoke \[(\d+)]`)
	successRegex := regexp.MustCompile(`^Program (.+) success`)

	type CallFrame struct {
		ProgramID             string
		Depth                 int
		InstructionIndex      int8
		InnerInstructionIndex int8
	}

	var callStack []CallFrame
	instructionIndex := int8(-1)
	innerInstructionIndex := int8(-1)

	for _, logX := range logs {
		logX = strings.TrimSpace(logX)
		if logX == "" {
			continue
		}

		if matches := invokeRegex.FindStringSubmatch(logX); matches != nil {
			programID := matches[1]
			depth, _ := strconv.Atoi(matches[2])

			// If this is an invoke with depth 1, it's a new instruction
			if depth == 1 {
				instructionIndex++
				innerInstructionIndex = -1 // Reset inner instruction counter
			} else {
				innerInstructionIndex++
			}
			if result[instructionIndex] == nil {
				result[instructionIndex] = make(map[int8][]structs.LogEntry)
			}

			logEntry := structs.LogEntry{
				InnerInstruction: innerInstructionIndex, // -1 for top-level, 0+ for inner
				Content:          logX,
				ProgramID:        programID,
				IsInvoke:         true,
				IsSuccess:        false,
			}

			// Add to appropriate logs
			if instructionIndex >= 0 {
				result[instructionIndex][innerInstructionIndex] = append(
					result[instructionIndex][innerInstructionIndex],
					logEntry,
				)
			}

			callStack = append(callStack, CallFrame{
				ProgramID:             programID,
				Depth:                 depth,
				InstructionIndex:      instructionIndex,
				InnerInstructionIndex: innerInstructionIndex,
			})

			continue
		}

		if matches := successRegex.FindStringSubmatch(logX); matches != nil {
			programID := matches[1]

			// Find matching invoke on call stack
			foundIndex := -1
			for i := len(callStack) - 1; i >= 0; i-- {
				if callStack[i].ProgramID == programID {
					foundIndex = i
					break
				}
			}

			if foundIndex != -1 {
				// Get the instruction index from the call stack frame
				frame := callStack[foundIndex]

				logEntry := structs.LogEntry{
					InnerInstruction: frame.InnerInstructionIndex,
					Content:          logX,
					ProgramID:        programID,
					IsInvoke:         false,
					IsSuccess:        true,
				}

				result[frame.InstructionIndex][frame.InnerInstructionIndex] = append(
					result[frame.InstructionIndex][frame.InnerInstructionIndex],
					logEntry,
				)
				callStack = append(callStack[:foundIndex], callStack[foundIndex+1:]...)
			}

			continue
		}

		currentContext := CallFrame{
			InstructionIndex:      instructionIndex,
			InnerInstructionIndex: -1, // Default to top-level
			Depth:                 1,
		}

		// If there are active calls on the stack, use the most recent one for context
		if len(callStack) > 0 {
			currentContext = callStack[len(callStack)-1]
		}

		logEntry := structs.LogEntry{
			InnerInstruction: currentContext.InnerInstructionIndex,
			Content:          logX,
			ProgramID:        currentContext.ProgramID,
			IsInvoke:         false,
			IsSuccess:        false,
		}
		result[currentContext.InstructionIndex][currentContext.InnerInstructionIndex] = append(
			result[currentContext.InstructionIndex][currentContext.InnerInstructionIndex],
			logEntry,
		)

	}

	return result
}
