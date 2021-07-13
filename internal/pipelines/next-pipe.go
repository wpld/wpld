package pipelines

import (
	"context"
)

type NextPipe func(ctx context.Context) error
