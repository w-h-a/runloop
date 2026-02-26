package run

import (
	"encoding/json"
	"time"
)

type Event interface {
	EventType() string
	RunID() string
}

type created struct {
	runID       string
	goal        string
	modelConfig ModelConfig
	tools       []ToolDefinition
	maxSteps    int
	tokenBudget int
	createdAt   time.Time
}

func (e created) EventType() string {
	return "created"
}

func (e created) RunID() string {
	return e.runID
}

type llmCallReq struct {
	runID     string
	stepID    string
	sequence  int
	input     json.RawMessage
	startedAt time.Time
}

func (e llmCallReq) EventType() string {
	return "llm_call_req"
}

func (e llmCallReq) RunID() string {
	return e.runID
}

type llmCallRsp struct {
	runID      string
	stepID     string
	output     json.RawMessage
	tokensUsed int
	duration   time.Duration
	err        string
}

func (e llmCallRsp) EventType() string {
	return "llm_call_rsp"
}

func (e llmCallRsp) RunID() string {
	return e.runID
}

type toolCallReq struct {
	runID     string
	stepID    string
	sequence  int
	toolName  string
	arguments json.RawMessage
	startedAt time.Time
}

func (e toolCallReq) EventType() string {
	return "tool_call_req"
}

func (e toolCallReq) RunID() string {
	return e.runID
}

type toolCallRsp struct {
	runID    string
	stepID   string
	toolName string
	result   json.RawMessage
	duration time.Duration
	err      string
}

func (e toolCallRsp) EventType() string {
	return "tool_call_rsp"
}

func (e toolCallRsp) RunID() string {
	return e.runID
}

type completed struct {
	runID string
}

func (e completed) EventType() string {
	return "completed"
}

func (e completed) RunID() string {
	return e.runID
}

type failed struct {
	runID          string
	err            string
	failedAtStepID string
}

func (e failed) EventType() string {
	return "failed"
}

func (e failed) RunID() string {
	return e.runID
}
