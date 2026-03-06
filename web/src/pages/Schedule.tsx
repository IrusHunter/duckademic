import { useEffect, useMemo, useState } from "react";
import {
  DEFAULT_FILTERS,
  FACULTY_OPTIONS,
  GROUP_OPTIONS,
  SPECIALTIES_BY_FACULTY,
  SUBGROUP_OPTIONS,
  VIEW_OPTIONS,
  YEAR_OPTIONS,
  dayOrder,
  scheduleData,
  typeIconId,
  type DayName,
  type Lesson,
} from "../data/schedule";

type DropdownName = "view" | "faculty" | "specialty" | "year" | "group" | "subgroup";

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(" ");
}

function matchesFilters(lesson: Lesson, filters: typeof DEFAULT_FILTERS): boolean {
  const a = lesson.audience;
  if (!a) return true;

  if (a.faculty && a.faculty !== filters.faculty) return false;
  if (a.specialty && a.specialty !== filters.specialty) return false;
  if (a.year && a.year !== filters.year) return false;
  if (a.group && a.group !== filters.group) return false;
  if (a.subgroup && a.subgroup !== filters.subgroup) return false;

  return true;
}

export default function Schedule() {
  const [currentDayIndex, setCurrentDayIndex] = useState(0);
  const [openDD, setOpenDD] = useState<DropdownName | null>(null);
  const [state, setState] = useState(DEFAULT_FILTERS);

  const specialties = useMemo(() => SPECIALTIES_BY_FACULTY[state.faculty] ?? [], [state.faculty]);

  // якщо змінили faculty і specialty стала невалідна
  useEffect(() => {
    if (!specialties.includes(state.specialty)) {
      setState((s) => ({ ...s, specialty: specialties[0] ?? "" }));
    }
  }, [specialties, state.specialty]);

  const dayName: DayName = dayOrder[currentDayIndex];

  const lessons = useMemo(() => {
    const raw = scheduleData[dayName] ?? [];
    return raw.filter((l) => matchesFilters(l, state));
  }, [dayName, state]);

  useEffect(() => {
    function onDocClick(e: MouseEvent) {
      const el = e.target as HTMLElement | null;
      if (!el) return;
      const inside = el.closest(".filter-group, .schedule-view-dd");
      if (!inside) setOpenDD(null);
    }
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === "Escape") setOpenDD(null);
    }

    document.addEventListener("click", onDocClick);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("click", onDocClick);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, []);

  function prevDay() {
    setCurrentDayIndex((i) => (i - 1 + dayOrder.length) % dayOrder.length);
  }

  function nextDay() {
    setCurrentDayIndex((i) => (i + 1) % dayOrder.length);
  }

  function toggleDD(name: DropdownName) {
    setOpenDD((cur) => (cur === name ? null : name));
  }

  function selectValue(name: DropdownName, value: string) {
    setState((s) => {
      const next = { ...s, [name]: value } as typeof DEFAULT_FILTERS;

      if (name === "faculty") {
        const nextSpecs = SPECIALTIES_BY_FACULTY[value] ?? [];
        next.specialty = nextSpecs[0] ?? "";
        next.group = GROUP_OPTIONS[0];
        next.subgroup = SUBGROUP_OPTIONS[0];
      }

      if (name === "specialty") {
        next.group = GROUP_OPTIONS[0];
        next.subgroup = SUBGROUP_OPTIONS[0];
      }

      if (name === "group") {
        next.subgroup = SUBGROUP_OPTIONS[0];
      }

      return next;
    });

    setOpenDD(null);
  }

  return (
    <main className="sidebar-space schedule-info">
      <section className="schedule-layout">
        <h1 className="schedule-title">Schedule</h1>

        <div className="schedule-card">
          <section className="schedule-main">
            <header className="schedule-main-header">
              <div className="schedule-main-left">
                <svg width="20" height="20">
                  <use href="/img/icons.svg#icon-SVG-11"></use>
                </svg>
                <span className="schedule-main-heading">Today&apos;s Classes</span>
              </div>

              <div className="schedule-main-controls">
                <div>
                  <button className="schedule-nav-btn" type="button" onClick={prevDay}>
                    <svg width="14" height="20">
                      <use href="/img/icons.svg#icon-vector-left"></use>
                    </svg>
                  </button>

                  <button className="schedule-day-btn" type="button">
                    <span className="schedule-day-label">{dayName}</span>
                  </button>

                  <button className="schedule-nav-btn" type="button" onClick={nextDay}>
                    <svg width="14" height="20">
                      <use href="/img/icons.svg#icon-vector-right"></use>
                    </svg>
                  </button>
                </div>

                <div className={cx("schedule-view-dd", openDD === "view" && "is-open")}>
                  <button
                    className="schedule-view-select"
                    type="button"
                    aria-expanded={openDD === "view"}
                    onClick={(e) => {
                      e.stopPropagation();
                      toggleDD("view");
                    }}
                  >
                    <span className="schedule-view-label">{state.view}</span>

                    <span className="schedule-view-caret">
                      <svg width="18" height="13">
                        <use href="/img/icons.svg#icon-vector-down"></use>
                      </svg>
                    </span>
                  </button>

                  <div className="dropdown-menu">
                    {VIEW_OPTIONS.map((opt) => (
                      <button
                        key={opt}
                        className={cx("dropdown-item", opt === state.view && "is-active")}
                        type="button"
                        onClick={() => selectValue("view", opt)}
                      >
                        {opt}
                      </button>
                    ))}
                  </div>
                </div>
              </div>
            </header>

            <ul className="schedule-classes">
              {lessons.length === 0 ? (
                <li className="class-row">
                  <div className="class-time"></div>
                  <article className="class-card">
                    <div className="class-card-inner">
                      <p className="class-teacher">No classes for selected filters.</p>
                    </div>
                  </article>
                </li>
              ) : (
                lessons.map((lesson) => {
                  const timeRange = `${lesson.start} – ${lesson.end}`;
                  const colorClass = `class-row--${lesson.color}`;
                  const iconId = typeIconId[lesson.type] ?? "icon-navigation";
                  const tagClass = lesson.type === "Lab" ? "class-tag white-tag" : "class-tag";

                  return (
                    <li key={`${dayName}-${lesson.start}-${lesson.title}`} className={cx("class-row", colorClass)}>
                      <div className="class-time">{timeRange}</div>

                      <article className="class-card">
                        <div className="class-card-inner">
                          <header className="class-card-header">
                            <h2 className="class-title">{lesson.title}</h2>
                            <span className={tagClass}>{lesson.type}</span>
                          </header>

                          <p className="class-teacher">{lesson.teacher}</p>

                          <div className="class-meta">
                            <div className="class-meta-item">
                              <span className="class-meta-icon" aria-hidden="true">
                                <svg width="18" height="18">
                                  <use href={`/img/icons.svg#${iconId}`}></use>
                                </svg>
                              </span>
                              <span>{lesson.modeText}</span>
                            </div>

                            {/* View: General / Date */}
                            {state.view === "Date" && <div className="class-meta-dates">{lesson.dates}</div>}
                          </div>
                        </div>
                      </article>
                    </li>
                  );
                })
              )}
            </ul>
          </section>

          <aside className="schedule-filters">
            {/* Faculty */}
            <div className={cx("filter-group", openDD === "faculty" && "is-open")}>
              <label className="filter-label" htmlFor="faculty-btn">
                Faculty
              </label>

              <button
                className="filter-select"
                id="faculty-btn"
                type="button"
                aria-expanded={openDD === "faculty"}
                onClick={(e) => {
                  e.stopPropagation();
                  toggleDD("faculty");
                }}
              >
                <span className="filter-caret-value">{state.faculty}</span>
                <span className="schedule-view-caret filter-caret">
                  <svg width="18" height="13">
                    <use href="/img/icons.svg#icon-vector-down"></use>
                  </svg>
                </span>
              </button>

              <div className="dropdown-menu">
                {FACULTY_OPTIONS.map((opt) => (
                  <button
                    key={opt}
                    className={cx("dropdown-item", opt === state.faculty && "is-active")}
                    type="button"
                    onClick={() => selectValue("faculty", opt)}
                  >
                    {opt}
                  </button>
                ))}
              </div>
            </div>

            {/* Specialty */}
            <div className={cx("filter-group", openDD === "specialty" && "is-open")}>
              <label className="filter-label" htmlFor="specialty-btn">
                Specialty
              </label>

              <button
                className="filter-select"
                id="specialty-btn"
                type="button"
                aria-expanded={openDD === "specialty"}
                onClick={(e) => {
                  e.stopPropagation();
                  toggleDD("specialty");
                }}
              >
                <span className="filter-caret-value">{state.specialty}</span>
                <span className="schedule-view-caret filter-caret">
                  <svg width="18" height="13">
                    <use href="/img/icons.svg#icon-vector-down"></use>
                  </svg>
                </span>
              </button>

              <div className="dropdown-menu">
                {specialties.map((opt) => (
                  <button
                    key={opt}
                    className={cx("dropdown-item", opt === state.specialty && "is-active")}
                    type="button"
                    onClick={() => selectValue("specialty", opt)}
                  >
                    {opt}
                  </button>
                ))}
              </div>
            </div>

            {/* Year */}
            <div className={cx("filter-group", openDD === "year" && "is-open")}>
              <label className="filter-label" htmlFor="year-btn">
                Year
              </label>

              <button
                className="filter-select"
                id="year-btn"
                type="button"
                aria-expanded={openDD === "year"}
                onClick={(e) => {
                  e.stopPropagation();
                  toggleDD("year");
                }}
              >
                <span className="filter-caret-value">{state.year}</span>
                <span className="schedule-view-caret filter-caret">
                  <svg width="18" height="13">
                    <use href="/img/icons.svg#icon-vector-down"></use>
                  </svg>
                </span>
              </button>

              <div className="dropdown-menu">
                {YEAR_OPTIONS.map((opt) => (
                  <button
                    key={opt}
                    className={cx("dropdown-item", opt === state.year && "is-active")}
                    type="button"
                    onClick={() => selectValue("year", opt)}
                  >
                    {opt}
                  </button>
                ))}
              </div>
            </div>

            {/* Group */}
            <div className={cx("filter-group", openDD === "group" && "is-open")}>
              <label className="filter-label" htmlFor="group-btn">
                Group
              </label>

              <button
                className="filter-select"
                id="group-btn"
                type="button"
                aria-expanded={openDD === "group"}
                onClick={(e) => {
                  e.stopPropagation();
                  toggleDD("group");
                }}
              >
                <span className="filter-caret-value">{state.group}</span>
                <span className="schedule-view-caret filter-caret">
                  <svg width="18" height="13">
                    <use href="/img/icons.svg#icon-vector-down"></use>
                  </svg>
                </span>
              </button>

              <div className="dropdown-menu">
                {GROUP_OPTIONS.map((opt) => (
                  <button
                    key={opt}
                    className={cx("dropdown-item", opt === state.group && "is-active")}
                    type="button"
                    onClick={() => selectValue("group", opt)}
                  >
                    {opt}
                  </button>
                ))}
              </div>
            </div>

            {/* Subgroup */}
            <div className={cx("filter-group", openDD === "subgroup" && "is-open")}>
              <label className="filter-label" htmlFor="subgroup-btn">
                Subgroup
              </label>

              <button
                className="filter-select"
                id="subgroup-btn"
                type="button"
                aria-expanded={openDD === "subgroup"}
                onClick={(e) => {
                  e.stopPropagation();
                  toggleDD("subgroup");
                }}
              >
                <span className="filter-caret-value">{state.subgroup}</span>
                <span className="filter-caret">
                  <svg width="18" height="13">
                    <use href="/img/icons.svg#icon-vector-down"></use>
                  </svg>
                </span>
              </button>

              <div className="dropdown-menu">
                {SUBGROUP_OPTIONS.map((opt) => (
                  <button
                    key={opt}
                    className={cx("dropdown-item", opt === state.subgroup && "is-active")}
                    type="button"
                    onClick={() => selectValue("subgroup", opt)}
                  >
                    {opt}
                  </button>
                ))}
              </div>
            </div>
          </aside>
        </div>
      </section>
    </main>
  );
}
