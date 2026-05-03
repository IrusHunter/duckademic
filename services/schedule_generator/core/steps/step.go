package steps

import (
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/components"
)

type PipelineStep interface {
	GetNextStep(*GeneratorContext) PipelineStep
	CanGoToTheNextStep() error
	InsertData(data any) error
	Process(components.ComponentIdentifier) (any, error)
}
