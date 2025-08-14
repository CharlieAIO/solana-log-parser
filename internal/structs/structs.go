package structs

type LogEntry struct {
	Content          string `json:"content"`
	ProgramID        string `json:"program_id,omitempty"`
	IsInvoke         bool   `json:"is_invoke"`
	IsSuccess        bool   `json:"is_success"`
	InnerInstruction int8   `json:"inner_instruction"` // -1 if not an inner instruction
}

type InstructionLogs struct {
	InstructionIndex int8       `json:"instruction_index"`
	Logs             []LogEntry `json:"logs"`
}

type CallFrame struct {
	ProgramID             string
	Depth                 int
	InstructionIndex      int8
	InnerInstructionIndex int8
}
