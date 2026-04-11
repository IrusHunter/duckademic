package entities

import (
	"fmt"

	"github.com/google/uuid"
)

// Discipline represents a university subject in the scheduling context.
type Discipline struct {
	ID   uuid.UUID           // Unique identifier of the Discipline.
	Name string              // Human-readable identifier of the Discipline.
	Load map[*LessonType]int // Required hours per lesson type.
}

// NewDiscipline creates a new Discipline instance.
//
// It requires discipline's id and name.
func NewDiscipline(id uuid.UUID, name string) *Discipline {
	return &Discipline{
		ID:   id,
		Name: name,
		Load: map[*LessonType]int{},
	}
}

// AddLoad sets the required number of hours for the lesson type.
//
// Returns an error if the lesson type is already registered.
func (d *Discipline) AddLoad(lt *LessonType, hours int) error {
	_, ok := d.Load[lt]
	if ok {
		return fmt.Errorf("load for %q already exists", lt.Name)
	}

	d.Load[lt] = hours
	return nil
}

// GetRequiredHours returns the required number of hours for the lesson type.
//
// Returns 0 if no hours were set for the lesson type.
func (d *Discipline) GetRequiredHours(lt *LessonType) int {
	// h, ok := d.Load[lt]
	// if !ok {
	// 	return 0
	// }
	return d.Load[lt]
}

// CutLoadTo limits the number of hours for each lesson type.
//
// Any lesson type with hours greater than the specified limit
// will be reduced to the given value. Other values remain unchanged.
func (d *Discipline) CutLoadTo(hours int) {
	for key := range d.Load {
		if d.Load[key] > hours {
			d.Load[key] = hours
		}
	}
}
