package cases

import (
	"github.com/spf13/afero"

	"wpld/internal/controllers/pipelines"
)

func StartProjectPipeline(fs afero.Fs) pipelines.Pipeline {
	return pipelines.NewPipeline()
}
