package components

import (
	"fmt"
	"slices"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

// DayBlocker selects days for student groups
type DayBlocker interface {
	// Basic interface for generator component
	GeneratorComponent[responses.LessonTypeDayDebt, *SetDayTypeError]
	// Add a SetDayTypeError to ErrorService if at not enough days per group
	SetDayTypes()
}

// NewDayBlocker creates a DayBlocker instance.
// It requires an ErrorService and a list of student groups.
func NewDayBlocker(
	sg []*entities.StudentGroup,
	es ErrorService[responses.LessonTypeDayDebt, *SetDayTypeError],
	w int, lfr float64,
) DayBlocker {
	db := dayBlocker{
		errorService:   es,
		weekCount:      w,
		lessonFillRate: lfr,
	}
	db.setGroupExtensions(sg)

	return &db
}

// Extension of group (store data to not calculate every time)
type groupExtension struct {
	group         *entities.StudentGroup // Original StudentGroup
	dayPriorities []float32              // Bigger number - better day for lessons (<0.99 if day is uncomfortable) (length - 7)
	freeDayCount  int                    // count of free (comfortable) days
}

func (ge *groupExtension) IsFreeDay(day int) bool {
	return ge.dayPriorities[day] > 0.99
}

func newGroupExtension(group *entities.StudentGroup) *groupExtension {
	ge := groupExtension{
		group:         group,
		dayPriorities: group.GetWeekDaysPriority(),
	}

	for day := range ge.dayPriorities {
		if ge.IsFreeDay(day) {
			ge.freeDayCount++
		}
	}

	return &ge
}

type dayBlocker struct {
	groupExtensions []groupExtension
	errorService    ErrorService[responses.LessonTypeDayDebt, *SetDayTypeError]
	weekCount       int
	lessonFillRate  float64
}

func (db *dayBlocker) SetDayTypes() {
	mainDayBlocked := make([]int, 7)

	for len(db.groupExtensions) != 0 {
		daysBlocked := make([]int, 7) // contains num of groups that chose this day
		copy(daysBlocked, mainDayBlocked)
		mainGroup := db.groupExtensions[0]
		for i := 0; i < len(db.groupExtensions); {
			group := db.groupExtensions[i]
			if !mainGroup.group.ConnectedTo(db.groupExtensions[i].group) && mainGroup.group != group.group {
				i++
				continue
			}
			db.groupExtensions = append(db.groupExtensions[:i], db.groupExtensions[i+1:]...)

			availableDays := []int{0, 1, 2, 3, 4, 5, 6}

			for _, lt := range group.group.GetLessonTypes() {
				requiredSlots := float64(group.group.GetSlotCountForLType(lt))

				for requiredSlots >= 0 { // break after error assigned
					// select day that free and blocked the fewest times
					min := 1000000000
					mIndex := -1
					for _, day := range availableDays {
						if group.IsFreeDay(day) && min > daysBlocked[day] {
							min = daysBlocked[day]
							mIndex = day
						}
					}

					// day not found
					if mIndex == -1 {
						db.errorService.AddError(&SetDayTypeError{
							LessonType:    lt,
							StudentGroup:  group.group,
							DayPriorities: group.dayPriorities,
							AvailableDays: availableDays,
							SlotsDept:     requiredSlots,
						})
						break // continue with next lesson type for group
					}

					// if an error occurs, ignore this day, delete it from available days, continue the search
					err := group.group.BindWeekday(lt, mIndex)
					if err != nil {
						dayIndex := slices.Index(availableDays, mIndex)
						availableDays = append(availableDays[:dayIndex], availableDays[dayIndex+1:]...)
						continue
					}

					// all good, add to blocked day
					daysBlocked[mIndex]++
					requiredSlots -= float64(db.weekCount*group.group.GetAverageSlotCountOnWeekday(mIndex)) * db.lessonFillRate
				}
			}

			if i == 0 {
				copy(mainDayBlocked, daysBlocked)
			}
		}
	}
}

func (db *dayBlocker) GetErrorService() ErrorService[responses.LessonTypeDayDebt, *SetDayTypeError] {
	return db.errorService
}
func (db *dayBlocker) Run() {
	db.SetDayTypes()
}

func (db *dayBlocker) setGroupExtensions(studentGroups []*entities.StudentGroup) {
	db.groupExtensions = make([]groupExtension, len(studentGroups))
	for i := range studentGroups {
		db.groupExtensions[i] = *newGroupExtension(studentGroups[i])
	}
	// sorts by connected groups count in decreasing order, and then by free day count in increasing order
	slices.SortFunc(db.groupExtensions, func(a, b groupExtension) int {
		if a.group.CountConnectedGroupsNumber() == b.group.CountConnectedGroupsNumber() {
			if a.freeDayCount == b.freeDayCount {
				return 0
			} else if a.freeDayCount > b.freeDayCount {
				return 1
			}
			return -1
		} else if a.group.CountConnectedGroupsNumber() < b.group.CountConnectedGroupsNumber() {
			return 1
		}
		return -1
	})
}

type SetDayTypeError struct {
	LessonType    *entities.LessonType
	StudentGroup  *entities.StudentGroup
	DayPriorities []float32
	AvailableDays []int
	SlotsDept     float64
}

func (e *SetDayTypeError) Error() string {
	return fmt.Sprintf("can't add a day of type %s to group %s", e.LessonType.Name, e.StudentGroup.Name)
}
func (e *SetDayTypeError) GeneratorResponseError() responses.LessonTypeDayDebt {
	return responses.LessonTypeDayDebt{
		StudentGroup: responses.CommonEntity{
			ID:   e.StudentGroup.ID,
			Name: e.StudentGroup.Name,
		},
		LessonType: responses.CommonEntity{
			ID:   e.LessonType.ID,
			Name: e.LessonType.Name,
		},
		SlotsDept: e.SlotsDept,
	}
}
