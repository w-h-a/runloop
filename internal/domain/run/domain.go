package run

import (
	"encoding/json"
	"time"
)

type RunStatus string

const (
	RunStatusPending   RunStatus = "pending"
	RunStatusRunning   RunStatus = "running"
	RunStatusCompleted RunStatus = "completed"
	RunStatusFailed    RunStatus = "failed"
)

type StepType string

const (
	StepTypeLLMCall  StepType = "llm_call"
	StepTypeToolCall StepType = "tool_call"
)

type StepStatus string

const (
	StepStatusStarted   StepStatus = "started"
	StepStatusCompleted StepStatus = "completed"
	StepStatusFailed    StepStatus = "failed"
)

type Run struct {
	ID          string
	Goal        string
	Status      RunStatus
	ModelConfig ModelConfig
	Tools       []ToolDefinition
	Steps       []Step
	MaxSteps    int
	TokenBudget int
	CreatedAt   time.Time
	Error       string
}

type Step struct {
	ID          string
	Sequence    int
	Type        StepType
	Status      StepStatus
	ToolName    string
	Input       json.RawMessage
	Output      json.RawMessage
	StartedAt   time.Time
	CompletedAt time.Time
	Duration    time.Duration
	TokensUsed  int
	Error       string
}

type ModelConfig struct {
	Provider    string
	Model       string
	Temperature float64
}

type ToolDefinition struct {
	Name        string
	Description string
	InputSchema json.RawMessage
	Endpoint    string
}
