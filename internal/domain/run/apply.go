package run

func (r *Run) Apply(e Event) error {
	switch ev := e.(type) {
	case created:
		r.ID = ev.runID
		r.Goal = ev.goal
		r.Status = RunStatusPending
		r.ModelConfig = ev.modelConfig
		r.Tools = ev.tools
		r.MaxSteps = ev.maxSteps
		r.TokenBudget = ev.tokenBudget
		r.CreatedAt = ev.createdAt
		return nil
	case llmCallReq:
		if r.isTerminal() {
			return ErrRunAlreadyCompleted
		}
		if r.MaxSteps > 0 && len(r.Steps) >= r.MaxSteps {
			return ErrMaxStepsExceeded
		}
		r.Status = RunStatusRunning
		r.Steps = append(r.Steps, Step{
			ID:        ev.stepID,
			Sequence:  ev.sequence,
			Type:      StepTypeLLMCall,
			Status:    StepStatusStarted,
			Input:     ev.input,
			StartedAt: ev.startedAt,
		})
		return nil
	case llmCallRsp:
		if r.isTerminal() {
			return ErrRunAlreadyCompleted
		}
		step := r.findStep(ev.stepID)
		if step == nil {
			return ErrStepNotFound
		}
		if r.TokenBudget > 0 && r.TotalTokensUsed+ev.tokensUsed > r.TokenBudget {
			return ErrTokenBudgetExceeded
		}
		step.Output = ev.output
		step.TokensUsed = ev.tokensUsed
		r.TotalTokensUsed += ev.tokensUsed
		step.Duration = ev.duration
		if len(ev.err) > 0 {
			step.Error = ev.err
			step.Status = StepStatusFailed
		} else {
			step.Status = StepStatusCompleted
		}
		return nil
	case toolCallReq:
		if r.isTerminal() {
			return ErrRunAlreadyCompleted
		}
		if r.MaxSteps > 0 && len(r.Steps) >= r.MaxSteps {
			return ErrMaxStepsExceeded
		}
		r.Status = RunStatusRunning
		r.Steps = append(r.Steps, Step{
			ID:        ev.stepID,
			Sequence:  ev.sequence,
			Type:      StepTypeToolCall,
			Status:    StepStatusStarted,
			ToolName:  ev.toolName,
			Input:     ev.arguments,
			StartedAt: ev.startedAt,
		})
		return nil
	case toolCallRsp:
		if r.isTerminal() {
			return ErrRunAlreadyCompleted
		}
		step := r.findStep(ev.stepID)
		if step == nil {
			return ErrStepNotFound
		}
		step.Output = ev.result
		step.Duration = ev.duration
		if len(ev.err) > 0 {
			step.Error = ev.err
			step.Status = StepStatusFailed
		} else {
			step.Status = StepStatusCompleted
		}
		return nil
	case completed:
		if r.isTerminal() {
			return ErrRunAlreadyCompleted
		}
		r.Status = RunStatusCompleted
		return nil
	case failed:
		if r.isTerminal() {
			return ErrRunAlreadyCompleted
		}
		r.Status = RunStatusFailed
		r.Error = ev.err
		r.FailedAtStepID = ev.failedAtStepID
		return nil
	default:
		return ErrUnknownEvent
	}
}

func (r *Run) isTerminal() bool {
	return r.Status == RunStatusCompleted || r.Status == RunStatusFailed
}

func (r *Run) findStep(id string) *Step {
	for i := range r.Steps {
		if r.Steps[i].ID == id {
			return &r.Steps[i]
		}
	}
	return nil
}
