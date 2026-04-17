package components

import (
	"fmt"
	"log"
)

// GeneratorComponentErrorTypes defines categories of generator errors.
type GeneratorComponentErrorTypes string

const (
	SetDayTypeErrorType          GeneratorComponentErrorTypes = "SetDayTypeErrorType"
	BoneWeekErrorType            GeneratorComponentErrorTypes = "BoneWeekErrorType"
	MissingLessonsAdderErrorType GeneratorComponentErrorTypes = "MissingLessonsAdderErrorType"
	ClassroomAssignerErrorType   GeneratorComponentErrorTypes = "ClassroomAssignerErrorType"

	unexpectedErrorType GeneratorComponentErrorTypes = "unexpectedErrorType"
)

// GeneratorComponentError represents a typed generator component error.
type GeneratorComponentError[TR any] interface {
	error // Basic interface for errors
	GeneratorResponseError() TR
}

// NewUnexpectedError create new unexpectedError instance.
// unexpectedError indicates an internal state that should be unreachable.
//
// It requires error description (d), method where error accrued (m), class name that contains this method (c),
// error that best describes internal state.
func NewUnexpectedError(d, c, m string, err error) *unexpectedError {
	ue := &unexpectedError{description: d, className: c, methodName: m, err: err}
	log.Println(ue.Error())
	return ue
}

type unexpectedError struct {
	description string
	className   string
	methodName  string
	err         error
}

func (e *unexpectedError) Error() string {
	return fmt.Sprintf("%s %s ==> %s. \n└-- basic error: %s", e.className, e.methodName, e.description, e.err.Error())
}
func (e *unexpectedError) GetTypeOfError() GeneratorComponentErrorTypes {
	return unexpectedErrorType
}

// ErrorService aggregates and manages errors produced by generator components.
type ErrorService[TR any, T GeneratorComponentError[TR]] interface {
	AddError(T) // Add error to collection. The service automatically handles ordering or deduplication.
	GetGeneratorResponseErrors() []TR
}

// NewErrorService creates new ErrorService instance
func NewErrorService[TR any, T GeneratorComponentError[TR]]() ErrorService[TR, T] {
	return &errorService[TR, T]{errors: []T{}}
}

type errorService[TR any, T GeneratorComponentError[TR]] struct {
	errors []T
}

func (ec *errorService[TR, T]) AddError(err T) {
	ec.errors = append(ec.errors, err)
}
func (ec *errorService[TR, T]) GetGeneratorResponseErrors() []TR {
	result := make([]TR, 0, len(ec.errors))

	for _, err := range ec.errors {
		result = append(result, err.GeneratorResponseError())
	}

	return result
}
