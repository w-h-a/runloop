package run

import "errors"

var (
	ErrRunNotFound         = errors.New("run not found")
	ErrRunAlreadyCompleted = errors.New("run already completed")
	ErrToolNotAllowed      = errors.New("tool not allowed")
	ErrMaxStepsExceeded    = errors.New("max steps exceeded")
	ErrTokenBudgetExceeded = errors.New("token budget exceeded")
	ErrStepNotFound        = errors.New("step not found")
	ErrUnknownEvent        = errors.New("unknown event")
)
