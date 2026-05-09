package components

import (
	"fmt"
	"slices"

	"github.com/IrusHunter/duckademic/services/schedule_generator/core/entities"
	"github.com/IrusHunter/duckademic/services/schedule_generator/core/responses"
)

// NewClassroomAssigner creates a new ClassroomAssigner instance that uses the Munkres
// assignment algorithm (Hungarian algorithm) to assign classrooms to lessons.
func NewClassroomAssigner() ClassroomAssigner {
	return &classroomAssigner{
		errorService: NewErrorService[responses.LessonWithoutClassroom, *ClassroomAssignError](),
		fault:        float32(0.000_000_1),
		maxValue:     1_000_000_000_000_000,
	}
}

// studentGroupService is the basic implementation of the StudentGroupService interface.
// It uses the Munkres assignment algorithm (Hungarian algorithm) to assign classrooms to lessons
type classroomAssigner struct {
	classrooms   []*entities.Classroom
	lessons      map[entities.LessonSlot][]*entities.Lesson
	slotsOrder   []entities.LessonSlot
	errorService ErrorService[responses.LessonWithoutClassroom, *ClassroomAssignError]
	// busynessOfClassrooms []int
	fault    float32 // Used for zero comparison.
	maxValue float32
}

func (ca *classroomAssigner) Run(input ClassroomAssignerInput) []responses.LessonWithoutClassroom {
	ca.setLessons(input.Lessons)
	ca.classrooms = input.Classrooms

	for _, slot := range ca.slotsOrder {
		classrooms := ca.getAvailableClassrooms(slot)
		n := len(classrooms)
		lessons := ca.lessons[slot]

		delta := n - len(lessons)
		if delta < 0 {
			NewUnexpectedError("number of classrooms is less than number of simultaneous lessons",
				"ClassroomAssigner", "AssignClassrooms",
				fmt.Errorf("classes: %d, lessons: %d, slot: %s", n, len(lessons), slot.String()),
			)
			continue
		}

		matrix := make([][]float32, 0, n) // matrix indices: [lesson][classroom].
		// fill matrix with values (lower values indicate better assignments)
		for _, lesson := range lessons {
			values := make([]float32, 0, n)

			for _, classroom := range classrooms {
				value := float32(1.0)
				if err := classroom.CheckLesson(lesson); err != nil {
					value = ca.maxValue
				}
				value *= lesson.Teacher.CalculateClassValueFor(lesson, classroom)

				values = append(values, value)
			}

			matrix = append(matrix, values)
		}

		// add unavailable lessons to ensure correct algorithm execution
		for range delta {
			values := make([]float32, n)
			for j := range values {
				values[j] = ca.maxValue
			}
			matrix = append(matrix, values)
		}

		// step 1 subtract the row minimum from each element in the row
		for i, row := range matrix {
			min := slices.Min(row)
			for j := range row {
				matrix[i][j] -= min
			}
		}

		// step 2 subtract the column minimum from each element in the column
		for j := range n {
			min := float32(1_000_000_000_000_000)
			for i := range matrix {
				if min > matrix[i][j] {
					min = matrix[i][j]
				}
			}

			for i := range matrix {
				matrix[i][j] -= min
			}
		}

		for {
			// step 3 cover zeros with lines
			// step 3.1 the maximum bipartite matching using the Kuhn algorithm
			busynessOfClassrooms := ca.ConstructBusynessOfClassrooms(n)
			for i := range matrix {
				ca.dfc(matrix, i, busynessOfClassrooms, []int{})
			}
			// step 3.2 build the minimum vertex cover from the maximum matching
			// step 3.2.1 find free lessons
			freeLessons := []int{}
			for i := range n {
				ind := slices.Index(busynessOfClassrooms, i)
				if ind == -1 {
					freeLessons = append(freeLessons, i)
				}
			}
			// step 3.2.2 find "visited" lessons and classrooms
			visitedLessons := make([]bool, n)
			visitedClassrooms := make([]bool, n)
			for _, fl := range freeLessons {
				ca.dfcVisited(matrix, visitedLessons, visitedClassrooms, fl, busynessOfClassrooms, []int{})
			}
			// step 3.2.3 end of the covering
			// cover all unvisited rows
			checkSum := 0
			rowLines := make([]bool, n)
			for i, visited := range visitedLessons {
				if !visited {
					rowLines[i] = true
					checkSum++
				}
			}
			// cover all visited columns
			columnLines := make([]bool, n)
			for j, visited := range visitedClassrooms {
				if visited {
					columnLines[j] = true
					checkSum++
				}
			}

			if checkSum < n {
				// step 4 shift zeros
				// step 4.1 find minimum uncovered value
				min := ca.maxValue
				for i := range n {
					if rowLines[i] {
						continue
					}
					for j := range n {
						if columnLines[j] {
							continue
						}
						if matrix[i][j] < min {
							min = matrix[i][j]
						}
					}
				}
				// step 4.2 subtract the found minimum from the uncovered values
				for i := range n {
					if rowLines[i] {
						continue
					}
					for j := range n {
						if columnLines[j] {
							continue
						}
						matrix[i][j] -= min
					}
				}
				// step 4.3 add the found minimum to interception of two lines
				for i := range n {
					if !rowLines[i] {
						continue
					}
					for j := range n {
						if !columnLines[j] {
							continue
						}
						matrix[i][j] += min
					}
				}

				continue // jump back to step 3
			}

			break // continue with step 5
		}

		// step 5 making final assignment
		// step 5.1 select independent zeros using the Kuhn algorithm
		busynessOfClassrooms := ca.ConstructBusynessOfClassrooms(n)
		for i := range n {
			if !ca.dfc(matrix, i, busynessOfClassrooms, []int{}) {
				if i < len(lessons) {
					ca.errorService.AddError(&ClassroomAssignError{Lesson: lessons[i]})
				}
			}
		}
		// step 5.2 assign classrooms to lessons
		for classroom, lesson := range busynessOfClassrooms {
			if lesson == -1 {
				continue
			}
			if lesson >= len(lessons) {
				continue
			}
			err := lessons[lesson].SetClassroom(classrooms[classroom])
			if err != nil {
				NewUnexpectedError("could not assign a classroom to the lesson",
					"ClassroomAssigner", "AssignClassrooms", newUnavailableClassroomForLessonError(
						lessons[lesson], ca.classrooms[classroom], err,
					))
			}
		}
	}

	return ca.errorService.GetGeneratorResponseErrors()
}
func (ca *classroomAssigner) GetComponentIdentifier() ComponentIdentifier {
	return MunkresClassroomAssignerID
}
func (ca *classroomAssigner) CheckAvailability(input ClassroomAssignerInput) error {
	ca.setLessons(input.Lessons)
	ca.classrooms = input.Classrooms

	for slot, lessons := range ca.lessons {
		n := len(ca.getAvailableClassrooms(slot))
		if len(lessons) > n {
			return fmt.Errorf("not enough classrooms for slot %s (%d < %d)", slot.String(), n, len(lessons))
		}
	}

	return nil
}

func (ca *classroomAssigner) getAvailableClassrooms(slot entities.LessonSlot) []*entities.Classroom {
	res := make([]*entities.Classroom, 0, len(ca.classrooms))
	for _, classroom := range ca.classrooms {
		if classroom.IsFree(slot) {
			res = append(res, classroom)
		}
	}
	return res
}
func (ca *classroomAssigner) setLessons(lessons []*entities.Lesson) {
	lm := make(map[entities.LessonSlot][]*entities.Lesson, 0)
	for _, lesson := range lessons {
		lm[lesson.LessonSlot] = append(lm[lesson.LessonSlot], lesson)
	}

	slotsOrder := make([]entities.LessonSlot, 0, len(lessons))
	for slot := range lm {
		slotsOrder = append(slotsOrder, slot)
	}

	ca.lessons = lm
	ca.slotsOrder = slotsOrder
}
func (ca *classroomAssigner) ConstructBusynessOfClassrooms(n int) []int {
	res := make([]int, n)
	for i := range n {
		res[i] = -1
	}
	return res
}

// depth-first-search
func (ca *classroomAssigner) dfc(matrix [][]float32, lessonIndex int, classroomBusyness, usedLessons []int) bool {
	ind := slices.Index(usedLessons, lessonIndex)
	if ind != -1 {
		return false
	}

	n := len(classroomBusyness)
	for j := range n {
		if matrix[lessonIndex][j] <= ca.fault {
			if classroomBusyness[j] == -1 {
				classroomBusyness[j] = lessonIndex
				return true
			} else {
				if ca.dfc(matrix, classroomBusyness[j], classroomBusyness, append(usedLessons, lessonIndex)) {
					classroomBusyness[j] = lessonIndex
					return true
				}
			}
		}
	}
	return false
}

func (ca *classroomAssigner) dfcVisited(matrix [][]float32, visitedRows, visitedCols []bool, row int, classroomB, used []int) {
	ind := slices.Index(used, row)
	if ind != -1 {
		return
	}

	n := len(classroomB)
	visitedRows[row] = true
	for j := range n {
		if matrix[row][j] < ca.fault {
			visitedCols[j] = true
			if classroomB[j] != -1 {
				ca.dfcVisited(matrix, visitedRows, visitedCols, classroomB[j], classroomB, append(used, row))
			}
		}
	}
}

// ==========================================================================================================
// ================================================= ERRORS =================================================
// ==========================================================================================================

// UnavailableClassroomForLessonError indicates that the ClassroomAssigner found a classroom for a lesson,
// but it cannot be assigned.
//
// Error: "failed to assign classroom "%classroom%" to lesson "%lesson%" after algorithm selection: %BasicError%"
type UnavailableClassroomForLessonError struct {
	Lesson     *entities.Lesson
	Classroom  *entities.Classroom
	BasicError error
}

func newUnavailableClassroomForLessonError(
	l *entities.Lesson, c *entities.Classroom, err error,
) *UnavailableClassroomForLessonError {
	return &UnavailableClassroomForLessonError{Lesson: l, Classroom: c, BasicError: err}
}

func (e *UnavailableClassroomForLessonError) Error() string {
	return fmt.Sprintf("failed to assign classroom %q to lesson %q after algorithm selection: %s",
		e.Classroom.RoomNumber, e.Lesson.String(), e.BasicError.Error(),
	)
}
func (e *UnavailableClassroomForLessonError) Unwrap() error {
	return e.BasicError
}

type ClassroomAssignError struct {
	*entities.Lesson
}

func (e *ClassroomAssignError) Error() string {
	return fmt.Sprintf("failed to assign classroom to %s", e.Lesson.String())
}
func (e *ClassroomAssignError) GeneratorResponseError() responses.LessonWithoutClassroom {
	return responses.LessonWithoutClassroom{
		CommonLesson: responses.CommonLesson{
			Teacher: responses.CommonEntity{
				ID:   e.Teacher.ID,
				Name: e.Teacher.UserName,
			},
			StudentGroup: responses.CommonEntity{
				ID:   e.StudentGroup.ID,
				Name: e.StudentGroup.Name,
			},
			Discipline: responses.CommonEntity{
				ID:   e.Discipline.ID,
				Name: e.Discipline.Name,
			},
			LessonType: responses.CommonEntity{
				ID:   e.Type.ID,
				Name: e.Type.Name,
			},
		},
		Day:  e.Day,
		Slot: e.Slot,
	}
}
