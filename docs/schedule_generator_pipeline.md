## General Diagram Description

The state diagram represents a step-by-step pipeline for schedule generation, where each state corresponds to a distinct phase of the algorithm. Transitions between states are controlled by the selected processing methods (algorithms), allowing flexibility in how each step is executed.

Each state in the diagram explicitly defines:

- **Allowed methods** — algorithms that can be used to process the step
- **Input data** — required context or intermediate results
- **Response type** — the output produced after the step is executed

This design enables modularity, extensibility, and the ability to experiment with different scheduling strategies at each stage.

## States Description

### 1. Setup State

This initial state is responsible for preparing all input data required for schedule generation. It includes study loads, student groups, teachers, classrooms, lesson types, and available time slots. At this stage, data is validated, normalized, and transformed into a unified format suitable for further processing. They all separate endpoints.

**Allowed methods**

- [/set-teachers](#schedule-generator-set-teachers)
- [/set-disciplines](#schedule-generator-set-disciplines)
- [/set-lesson-types](#schedule-generator-set-lesson-types)
- [/set-lesson-type-assignments](#schedule-generator-set-lesson-type-assignments)
- [/set-student-groups](#schedule-generator-set-student-groups)
- [/set-teacher-loads](#schedule-generator-set-teacher-loads)
- [/set-classrooms](#schedule-generator-set-classrooms)

### 2. Weekday Allocation State

In this state, weekdays are preliminarily distributed among different lesson types for each student group.
A simplified rule is applied: each day is assigned only one type of lesson. The allocation is based on the study load, meaning that lesson types with fewer hours may receive fewer days. This step structures the weekly schedule and simplifies further time slot assignment.

**Allowed methods**

- even_weekday_allocator

**Response** [=> DaysForLessonTypes](schemas.md#schedule-generator-days-for-lesson-types)

### 3. Weekly Time Slot Assignment State

This state assigns specific time slots within a single week. Irregular or one-time constraints (e.g., temporary teacher unavailability) are ignored to simplify the problem. The result is a **skeleton schedule**, representing the core structure of the timetable. Graph-based algorithms (e.g., bipartite matching) are typically used here, optionally improved with heuristics to reduce soft constraint violations.

**Allowed methods**

- one_per_week_time_slot_assigner
- brute_time_slot_assigner

**Response** [=> BoneLessons](schemas.md#schedule-generator-bone-lessons)

### 4. Weekly Classroom Assignment State

This optional state assigns classrooms to lessons in the skeleton schedule. Since different groups may have classes on different days, early global classroom assignment may be inefficient. Therefore, this step can be skipped or only partially executed depending on the scenario.

**Allowed methods**

- munkres_classroom_assigner

**Response** [=> BoneLessonsWithC](schemas.md#schedule-generator-bone-lessons-with-c)

### 5. Weekly Schedule Expansion State

At this stage, the skeleton schedule is extended to all academic weeks.Teacher-specific constraints and availability preferences are taken into account, which may introduce gaps but results in a more realistic and adaptable timetable.

**Allowed methods**

- any (fixed - single predefined method, not configurable)

**Response** [=> GeneratedLessons](schemas.md#schedule-generator-generated-lessons)

### 6. Full Time Slot Assignment State

This state handles the assignment of time slots for previously unassigned (“floating”) lessons. Due to increased problem size (multiple weeks), computational complexity is higher, so simpler heuristic approaches (e.g., greedy algorithms) are preferred over exact methods.

**Allowed methods**

- one_per_week_time_slot_assigner
- brute_time_slot_assigner

**Response** [=> GeneratedLessons](schemas.md#schedule-generator-generated-lessons)

### 7. Full Classroom Assignment State

This state assigns classrooms to all remaining lessons that do not yet have one. Previously assigned classrooms are not modified. The assignment is performed per time slot using the **Hungarian algorithm** on a bipartite graph to achieve an optimal distribution of lessons across available classrooms.

**Allowed methods**

- munkres_classroom_assigner

**Response** [=> GeneratedLessonsWithC](schemas.md#schedule-generator-generated-lessons-with-c)

### 8. Extraction State

This final state exports the generated schedule.
It produces two related datasets:

- Study load models (group–teacher–subject–lesson type–ID)
- Scheduled lessons with assigned time slots and classrooms

At this stage, the final schedule can also be evaluated based on predefined quality criteria.
They all separate endpoints.

**Allowed methods**

- [/get-study-loads](#schedule-generator-get-study-loads)
- [/get-lessons](#schedule-generator-get-lessons)
- [/get-fault](#schedule-generator-get-fault)
