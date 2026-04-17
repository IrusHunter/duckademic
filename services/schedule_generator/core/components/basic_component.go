package components

// GeneratorComponent represents any component participating in the generation process.
type GeneratorComponent[TR any, T GeneratorComponentError[TR]] interface {
	Run()                                 // The main improvement of schedule for generator
	GetErrorService() ErrorService[TR, T] // Each component must embed an ErrorService to report generation errors.
}
