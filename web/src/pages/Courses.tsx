import { useMemo } from "react";

type Course = {
  title: string;
  teacher: string;
  description: string;
  assignments: number;
  students: number;
  grade: string;
  cardColorClass: string; // course-card-blue | course-card-mint | course-card-rose
  deadlineClass: string; // course-deadline-warning | info | danger
  deadlineText: string;
};

const courses: Course[] = [
  {
    title: "Introduction to Programming",
    teacher: "Dr. Smith",
    description: "Learn programming fundamentals with Python",
    assignments: 5,
    students: 32,
    grade: "5+",
    cardColorClass: "course-card-blue",
    deadlineClass: "course-deadline-warning",
    deadlineText: "Deadline: Tomorrow 9:00 AM",
  },
  {
    title: "Data Science",
    teacher: "Prof. Johnson",
    description: "Data analysis and machine learning concepts",
    assignments: 5,
    students: 32,
    grade: "5-",
    cardColorClass: "course-card-mint",
    deadlineClass: "course-deadline-info",
    deadlineText: "Deadline: Friday 2:00 PM",
  },
  {
    title: "Web Development",
    teacher: "Dr. Williams",
    description: "Modern web applications with React",
    assignments: 8,
    students: 32,
    grade: "4+",
    cardColorClass: "course-card-rose",
    deadlineClass: "course-deadline-danger",
    deadlineText: "Deadline: Yesterday 11:00 AM",
  },
  {
    title: "Database",
    teacher: "Dr. Smith",
    description: "Learn programming fundamentals with Python",
    assignments: 5,
    students: 32,
    grade: "5+",
    cardColorClass: "course-card-mint",
    deadlineClass: "course-deadline-warning",
    deadlineText: "Deadline: Wednesday 9:00 AM",
  },
];

function CourseCard({
  course,
  onOpen,
  onAssignments,
}: {
  course: Course;
  onOpen: (c: Course) => void;
  onAssignments: (c: Course) => void;
}) {
  return (
    <li className={`course-card ${course.cardColorClass}`}>
      <div className="course-card-header">
        <div className="course-card-text-head">
          <h2 className="course-title">{course.title}</h2>
          <p className="course-teacher">{course.teacher}</p>
        </div>

        <span className="course-grade">{course.grade}</span>
      </div>

      <p className="course-description">{course.description}</p>

      <ul className="course-meta">
        <li className="course-meta-item">
          <svg className="course-meta-icon" width="16" height="16" aria-hidden="true">
            <use href="/img/icons.svg#icon-SVG-5" />
          </svg>
          <span className="course-meta-text">{course.assignments} assignments</span>
        </li>

        <li className="course-meta-item">
          <svg className="course-meta-icon" width="16" height="16" aria-hidden="true">
            <use href="/img/icons.svg#icon-SVG-9" />
          </svg>
          <span className="course-meta-text">{course.students} students</span>
        </li>
      </ul>

      <div className={`course-deadline ${course.deadlineClass}`}>
        <svg className="course-deadline-icon" width="16" height="16" aria-hidden="true">
          <use href="/img/icons.svg#icon-SVG-11" />
        </svg>
        <p className="course-deadline-text">{course.deadlineText}</p>
      </div>

      <div className="course-actions">
        <button className="course-btn course-btn-primary" type="button" onClick={() => onOpen(course)}>
          Open Course
        </button>
        <button className="course-btn course-btn-secondary" type="button" onClick={() => onAssignments(course)}>
          Assignments
        </button>
      </div>
    </li>
  );
}

export default function Courses() {
  const list = useMemo(() => courses, []);

  function handleOpen(course: Course) {
    // TODO: навігація / відкриття курсу
    console.log("Open course:", course.title);
  }

  function handleAssignments(course: Course) {
    // TODO: навігація на assignments
    console.log("Assignments for:", course.title);
  }

  return (
    <main className="sidebar-space courses-info">
      <div className="courses-info-header">
        <h1 className="courses-info-title">My Courses</h1>
        <p className="courses-info-subtitle">Manage your enrolled courses</p>
      </div>

      <ul className="courses-list">
        {list.map((c) => (
          <CourseCard
            key={`${c.title}-${c.teacher}`}
            course={c}
            onOpen={handleOpen}
            onAssignments={handleAssignments}
          />
        ))}
      </ul>
    </main>
  );
}
