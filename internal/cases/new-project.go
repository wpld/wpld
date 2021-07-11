package cases

import (
	"context"

	"wpld/internal/controllers/pipelines"
)

var NewProject = pipelines.NewPipeline(
	newProjectPromptPipe,
	newProjectMarshalPipe,
)

func newProjectPromptPipe(ctx context.Context, next pipelines.NextPipe) error {
	return next(ctx)
}

func newProjectMarshalPipe(ctx context.Context, next pipelines.NextPipe) error {
	return next(ctx)
}
