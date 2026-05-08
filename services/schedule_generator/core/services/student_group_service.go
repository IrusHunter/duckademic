package services

import (
	"fmt"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
	externalEntities "github.com/IrusHunter/duckademic/services/schedule_generator/entities"

	"github.com/google/uuid"
)

// StudentGroupService aggregates and manages student groups that the generator works with.
type StudentGroupService interface {
	// Returns a pointer to the student group with the given ID.
	Find(uuid.UUID) *entities.StudentGroup
	// Returns a slice with all student groups as pointers.
	GetAll() []*entities.StudentGroup
	// Returns windows (gaps between busy slots) and the sum of it.
	GetWindows() ([]responses.StudentGroupWindow, int)
	// Returns the overlapping lessons and the sum of it.
	GetLessonOverlapping() ([]responses.StudentGroupLessonOverlap, int)
	// Returns the overtime lessons (above the daily limit) and the sum of it.
	GetOvertimeLessons() ([]responses.StudentGroupOvertimeLesson, int)
	// Returns the lessons scheduled on days that are not allowed for their type, and the sum of them.
	GetInvalidLessonsByType() ([]responses.StudentGroupInvalidLesson, int)
	// Clears week binding of student groups.
	UnbindWeeks()
	GetSlotDeficitOfReservedSlotsForLT() int
}

// NewStudentGroupService creates a new StudentGroupService basic instance.
//
// It requires an array of database student groups (sg), day load limit (dll), and a busy grid for them (bg).
//
// Returns an error if any student group is an invalid model.
func NewStudentGroupService(sg []externalEntities.StudentGroup, dl int, bg [][]float32) (StudentGroupService, error) {
	sgs := studentGroupService{
		studentGroups: make([]*entities.StudentGroup, len(sg)),
	}

	for i := range sg {
		sgs.studentGroups[i] = entities.NewDefaultStudentGroup(
			sg[i].ID, sg[i].Name, dl, sg[i].StudentCount, entities.NewBusyGrid(bg),
		)
		// studentGroup := sgs.studentGroups[i]

		// set military day by marks slots on this day as blocked
		// md := sg[i].MilitaryDay
		// if md != -1 {
		// 	if err := studentGroup.CheckWeekDay(md); err != nil {
		// 		return nil, err
		// 	}
		// 	studentGroup.BlockWeekDay(md)
		// }
	}

	// create connection for student groups
	for i := range sg {
		mainGroup := sgs.Find(sg[i].ID)
		for _, gID := range sg[i].ConnectedGroups {
			otherG := sgs.Find(gID)
			if otherG == nil {
				return nil, fmt.Errorf("Can't find connected group %s for group %s (%s)", gID, sg[i].Name, sg[i].ID)
			}
			otherG.AddConnectedGroup(mainGroup)
		}
	}

	return &sgs, nil
}

// studentGroupService is the basic implementation of the StudentGroupService interface.
type studentGroupService struct {
	studentGroups []*entities.StudentGroup
}

func (sgs *studentGroupService) GetAll() []*entities.StudentGroup {
	return sgs.studentGroups
}
func (sgs *studentGroupService) Find(id uuid.UUID) *entities.StudentGroup {
	for i := range sgs.studentGroups {
		if sgs.studentGroups[i].ID == id {
			return sgs.studentGroups[i]
		}
	}

	return nil
}
func (sgs *studentGroupService) GetWindows() (res []responses.StudentGroupWindow, count int) {
	res = []responses.StudentGroupWindow{}
	for _, g := range sgs.studentGroups {
		windows := g.GetWindows()
		if len(windows) != 0 {
			count += len(windows)
			res = append(res, responses.StudentGroupWindow{
				CommonEntity: responses.CommonEntity{
					ID:   g.ID,
					Name: g.Name,
				},
				Windows: responses.FormTimeSlots(windows),
			})
		}
	}
	return
}
func (sgs *studentGroupService) GetLessonOverlapping() (res []responses.StudentGroupLessonOverlap, count int) {
	res = []responses.StudentGroupLessonOverlap{}

	for _, studentGroup := range sgs.studentGroups {
		lo := studentGroup.GetLessonOverlapping()
		if len(lo) != 0 {
			count += len(lo)
			res = append(res, responses.StudentGroupLessonOverlap{
				CommonEntity: responses.CommonEntity{
					ID:   studentGroup.ID,
					Name: studentGroup.Name,
				},
				OverlapLessons: responses.FormTimeSlots(lo),
			})
		}
	}

	return
}
func (sgs *studentGroupService) GetOvertimeLessons() (res []responses.StudentGroupOvertimeLesson, count int) {
	res = []responses.StudentGroupOvertimeLesson{}

	for _, sg := range sgs.studentGroups {
		days := sg.GetOvertimeLessons()
		if len(days) != 0 {
			count += len(days)
			res = append(res, responses.StudentGroupOvertimeLesson{
				CommonEntity: responses.FormCommonEntity(sg.ID, sg.Name),
				OvertimeDays: days,
			})
		}
	}

	return
}
func (sgs *studentGroupService) GetInvalidLessonsByType() (res []responses.StudentGroupInvalidLesson, count int) {
	res = []responses.StudentGroupInvalidLesson{}

	for _, sg := range sgs.studentGroups {
		il := sg.GetInvalidLessonsByType()
		if len(il) != 0 {
			count += len(il)
			res = append(res, responses.StudentGroupInvalidLesson{
				CommonEntity:   responses.FormCommonEntity(sg.ID, sg.Name),
				InvalidLessons: responses.FormInvalidLessonByType(il),
			})
		}
	}

	return
}
func (sgs *studentGroupService) UnbindWeeks() {
	for _, group := range sgs.studentGroups {
		group.UnbindWeeks()
	}
}
func (sgs *studentGroupService) GetSlotDeficitOfReservedSlotsForLT() int {
	count := 0

	for _, sg := range sgs.studentGroups {
		count += sg.GetSlotDeficitOfReservedSlotsForLT()
	}

	return count
}
