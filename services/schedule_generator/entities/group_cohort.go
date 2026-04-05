package entities

import (
	"fmt"

	"github.com/google/uuid"
)

type GroupCohort struct {
	ID     uuid.UUID      `json:"id"`
	Name   string         `json:"name"`
	Groups []StudentGroup `json:"groups"`
}

func (gc *GroupCohort) ValidateGroups() error {
	if len(gc.Groups) == 0 {
		return fmt.Errorf("group cohort should contain at least one group")
	}

	return nil
}

type StudentGroup struct {
	ID              uuid.UUID   `json:"id"`
	Name            string      `json:"name"`
	ConnectedGroups []uuid.UUID `json:"connected_groups"`
	StudentCount    int         `json:"student_count"`
}

func (sg *StudentGroup) ValidateName() error {
	if len(sg.Name) == 0 {
		return fmt.Errorf("student group name must not be empty")
	}
	return nil
}
func (sg *StudentGroup) ValidateConnectedGroups() error {
	seen := make(map[uuid.UUID]struct{})
	for _, id := range sg.ConnectedGroups {
		if _, exists := seen[id]; exists {
			return fmt.Errorf("connected_groups contains duplicate id: %s", id)
		}
		seen[id] = struct{}{}
	}
	return nil
}
func (sg *StudentGroup) ValidateStudentCount() error {
	if sg.StudentCount <= 0 {
		return fmt.Errorf("student_count must be positive")
	}
	return nil
}
