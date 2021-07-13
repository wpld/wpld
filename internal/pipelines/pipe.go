package pipelines

import (
	"context"
)

type Pipe func(ctx context.Context, next NextPipe) error
