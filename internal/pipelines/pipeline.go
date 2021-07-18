package pipelines

import (
	"context"
)

type Pipeline struct {
	pipes []Pipe
}

func NewPipeline(pipes ...Pipe) Pipeline {
	return Pipeline{
		pipes: pipes,
	}
}

func (p Pipeline) Run(ctx context.Context) error {
	if len(p.pipes) == 0 {
		return nil
	}

	select {
	case <-ctx.Done():
		return nil
	default:
		return p.pipes[0](ctx, func(nextCtx context.Context) error {
			return NewPipeline(p.pipes[1:]...).Run(nextCtx)
		})
	}
}
