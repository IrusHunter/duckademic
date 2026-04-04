package entities

import (
	"fmt"
	"strings"

	"github.com/IrusHunter/duckademic/shared/db"
	"github.com/google/uuid"
)

type LessonType struct {
	ID            uuid.UUID
	Name          string
	HoursValue    int
	ReservedWeeks []int
}

type LessonTypeRequest struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	HoursValue    int       `json:"hours_value"`
	ReservedWeeks string    `json:"reserved_weeks"`
}

func (lt LessonTypeRequest) String() string {
	parts := make([]string, 0, 4)

	if lt.ID != uuid.Nil {
		parts = append(parts, fmt.Sprintf("id: %s", lt.ID))
	}

	parts = append(parts, fmt.Sprintf("name: %s", lt.Name))
	parts = append(parts, fmt.Sprintf("value: %d", lt.HoursValue))
	parts = append(parts, fmt.Sprintf("reserved_weeks: %s", lt.ReservedWeeks))

	return fmt.Sprintf("LessonType{%s}", strings.Join(parts, ", "))
}

func (lt *LessonTypeRequest) ValidateName() error {
	if len(lt.Name) == 0 {
		return fmt.Errorf("name required")
	}
	return nil
}

func (lt *LessonTypeRequest) ToLessonType() (LessonType, error) {
	rw, err := db.StringToIntSlice(lt.ReservedWeeks)
	if err != nil {
		return LessonType{}, fmt.Errorf("invalid reserved weeks: %w", err)
	}

	return LessonType{
		ID:            lt.ID,
		Name:          lt.Name,
		HoursValue:    lt.HoursValue,
		ReservedWeeks: rw,
	}, nil
}
