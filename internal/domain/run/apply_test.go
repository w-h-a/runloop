package run

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApply_LLMCallHappyPath(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID: "run-1",
		goal:  "solve a task",
		modelConfig: ModelConfig{
			Provider:    "openai",
			Model:       "gpt-4",
			Temperature: 0.7,
		},
		tools: []ToolDefinition{
			{Name: "bash", Description: "run commands"},
		},
		maxSteps:    10,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, RunStatusPending, r.Status)

	// Act
	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		input:     json.RawMessage(`{"prompt":"hello"}`),
		startedAt: now,
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, RunStatusRunning, r.Status)
	assert.Len(t, r.Steps, 1)

	// Act
	err = r.Apply(llmCallRsp{
		runID:      "run-1",
		stepID:     "step-1",
		output:     json.RawMessage(`{"response":"world"}`),
		tokensUsed: 100,
		duration:   500 * time.Millisecond,
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, StepStatusCompleted, r.Steps[0].Status)
	assert.Equal(t, 100, r.Steps[0].TokensUsed)
	assert.Equal(t, 100, r.TotalTokensUsed)

	// Act
	err = r.Apply(completed{
		runID: "run-1",
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, RunStatusCompleted, r.Status)
}

func TestApply_ToolCallHappyPath(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, RunStatusPending, r.Status)

	// Act
	err = r.Apply(toolCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		toolName:  "bash",
		arguments: json.RawMessage(`{"cmd":"ls"}`),
		startedAt: now,
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, RunStatusRunning, r.Status)
	assert.Len(t, r.Steps, 1)
	assert.Equal(t, StepTypeToolCall, r.Steps[0].Type)
	assert.Equal(t, "bash", r.Steps[0].ToolName)

	// Act
	err = r.Apply(toolCallRsp{
		runID:    "run-1",
		stepID:   "step-1",
		toolName: "bash",
		result:   json.RawMessage(`{"output":"file.txt"}`),
		duration: 200 * time.Millisecond,
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, StepStatusCompleted, r.Steps[0].Status)

	// Act
	err = r.Apply(completed{
		runID: "run-1",
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, RunStatusCompleted, r.Status)
}

func TestApply_LLMCallFailure(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		input:     json.RawMessage(`{"prompt":"hello"}`),
		startedAt: now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallRsp{
		runID:    "run-1",
		stepID:   "step-1",
		duration: 100 * time.Millisecond,
		err:      "permission denied",
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, StepStatusFailed, r.Steps[0].Status)
	assert.Equal(t, "permission denied", r.Steps[0].Error)
}

func TestApply_ToolCallFailure(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(toolCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		toolName:  "bash",
		arguments: json.RawMessage(`{"cmd":"ls"}`),
		startedAt: now,
	})
	require.NoError(t, err)

	err = r.Apply(toolCallRsp{
		runID:    "run-1",
		stepID:   "step-1",
		toolName: "bash",
		duration: 100 * time.Millisecond,
		err:      "permission denied",
	})
	require.NoError(t, err)

	// Assert
	assert.Equal(t, StepStatusFailed, r.Steps[0].Status)
	assert.Equal(t, "permission denied", r.Steps[0].Error)
}

func TestApply_RejectsStepAfterRunCompleted(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(completed{
		runID: "run-1",
	})
	require.NoError(t, err)

	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		startedAt: now,
	})

	// Assert
	assert.ErrorIs(t, err, ErrRunAlreadyCompleted)
}

func TestApply_RejectsStepAfterRunFailed(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(failed{
		runID:          "run-1",
		err:            "failed",
		failedAtStepID: "step-1",
	})
	require.NoError(t, err)

	err = r.Apply(toolCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		toolName:  "bash",
		arguments: json.RawMessage(`{"cmd":"ls"}`),
		startedAt: now,
	})

	// Assert
	assert.ErrorIs(t, err, ErrRunAlreadyCompleted)
	assert.Equal(t, "failed", r.Error)
	assert.Equal(t, "step-1", r.FailedAtStepID)
}

func TestApply_TokenBudgetExceeded(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 100,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		startedAt: now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallRsp{
		runID:      "run-1",
		stepID:     "step-1",
		output:     json.RawMessage(`{"response":"world"}`),
		tokensUsed: 101,
		duration:   100 * time.Millisecond,
	})

	// Assert
	assert.ErrorIs(t, err, ErrTokenBudgetExceeded)
}

func TestApply_TokenBudgetExceededAcrossSteps(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    10,
		tokenBudget: 150,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		startedAt: now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallRsp{
		runID:      "run-1",
		stepID:     "step-1",
		output:     json.RawMessage(`{"response":"world"}`),
		tokensUsed: 80,
		duration:   100 * time.Millisecond,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-2",
		sequence:  2,
		startedAt: now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallRsp{
		runID:      "run-1",
		stepID:     "step-2",
		output:     json.RawMessage(`{}`),
		tokensUsed: 80,
		duration:   time.Second,
	})

	// Assert
	assert.ErrorIs(t, err, ErrTokenBudgetExceeded)
}

func TestApply_MaxStepsExceeded(t *testing.T) {
	if len(os.Getenv("INTEGRATION")) > 0 {
		t.Skip()
	}

	// Arrange
	r := &Run{}
	now := time.Now()

	// Act
	err := r.Apply(created{
		runID:       "run-1",
		goal:        "test",
		maxSteps:    1,
		tokenBudget: 1000,
		createdAt:   now,
	})
	require.NoError(t, err)

	err = r.Apply(llmCallReq{
		runID:     "run-1",
		stepID:    "step-1",
		sequence:  1,
		startedAt: now,
	})
	require.NoError(t, err)

	err = r.Apply(toolCallReq{
		runID:     "run-1",
		stepID:    "step-2",
		sequence:  2,
		toolName:  "bash",
		arguments: json.RawMessage(`{}`),
		startedAt: now,
	})

	// Assert
	assert.ErrorIs(t, err, ErrMaxStepsExceeded)
}
