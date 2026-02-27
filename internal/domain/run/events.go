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

func NewCreatedEvent(
	runID string,
	goal string,
	modelConfig ModelConfig,
	tools []ToolDefinition,
	maxSteps int,
	tokenBudget int,
	createdAt time.Time,
) Event {
	return created{
		runID:       runID,
		goal:        goal,
		modelConfig: modelConfig,
		tools:       tools,
		maxSteps:    maxSteps,
		tokenBudget: tokenBudget,
		createdAt:   createdAt,
	}
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

func NewLLMCallReqEvent(
	runID string,
	stepID string,
	sequence int,
	input json.RawMessage,
	startedAt time.Time,
) Event {
	return llmCallReq{
		runID:     runID,
		stepID:    stepID,
		sequence:  sequence,
		input:     input,
		startedAt: startedAt,
	}
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

func NewLLMCallRspEvent(
	runID string,
	stepID string,
	output json.RawMessage,
	tokensUsed int,
	duration time.Duration,
	err string,
) Event {
	return llmCallRsp{
		runID:      runID,
		stepID:     stepID,
		output:     output,
		tokensUsed: tokensUsed,
		duration:   duration,
		err:        err,
	}
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

func NewToolCallReqEvent(
	runID string,
	stepID string,
	sequence int,
	toolName string,
	arguments json.RawMessage,
	startedAt time.Time,
) Event {
	return toolCallReq{
		runID:     runID,
		stepID:    stepID,
		sequence:  sequence,
		toolName:  toolName,
		arguments: arguments,
		startedAt: startedAt,
	}
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

func NewToolCallRspEvent(
	runID string,
	stepID string,
	toolName string,
	result json.RawMessage,
	duration time.Duration,
	err string,
) Event {
	return toolCallRsp{
		runID:    runID,
		stepID:   stepID,
		toolName: toolName,
		result:   result,
		duration: duration,
		err:      err,
	}
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

func NewCompletedEvent(
	runID string,
) Event {
	return completed{
		runID: runID,
	}
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

func NewFailedEvent(
	runID string,
	err string,
	failedAtStepID string,
) Event {
	return failed{
		runID:          runID,
		err:            err,
		failedAtStepID: failedAtStepID,
	}
}
