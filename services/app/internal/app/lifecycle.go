package app

import (
	"context"

	"github.com/highonsemicolon/aura/pkg/logging"
)

type Lifecycle struct {
	steps []shutdownStep
	log   logging.Logger
}

type shutdownStep struct {
	name string
	fn   func(ctx context.Context) error
}

func NewLifecycle(log logging.Logger) *Lifecycle {
	return &Lifecycle{log: log}
}

func (l *Lifecycle) Add(name string, fn func(ctx context.Context) error) {
	l.steps = append(l.steps, shutdownStep{name, fn})
}

func (l *Lifecycle) Shutdown(ctx context.Context) error {
	for i := len(l.steps) - 1; i >= 0; i-- {
		step := l.steps[i]
		l.log.InfoF("%s: shutting down", step.name)
		if err := step.fn(ctx); err != nil {
			l.log.ErrorF("error while shutting down %s: %v", step.name, err)
		} else {
			l.log.InfoF("%s: shut down successfully", step.name)
		}
	}
	return nil
}
