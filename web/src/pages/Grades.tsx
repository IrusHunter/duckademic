import { useEffect, useMemo, useRef, useState } from "react";

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(" ");
}

type Student = {
  id: string;
  name: string;
  avatar: string;
  group: string; // наприклад SE-32
  course: string; // фільтр Course
  year: string; // фільтр Year
  score: number; // для рейтингу
  isMe?: boolean;
};

type CourseBreakdown = {
  id: string;
  title: string;
  markText: string; // 5+
  total: number; // 97
  blocks: Array<{ label: string; value: string }>;
};

const COURSE_OPTIONS = ["All courses", "Software Engineering", "Data Science", "Web Development"] as const;
const YEAR_OPTIONS = ["All years", "1 Bachelor", "2 Bachelor", "3 Bachelor (Current)", "4 Bachelor"] as const;

const studentsSeed: Student[] = [
  { id: "s1", name: "Emily Johnson", avatar: "/img/profile_pic.png", group: "SE-32", course: "Software Engineering", year: "3 Bachelor (Current)", score: 99.81 },
  { id: "s2", name: "Sarah Wilson", avatar: "/img/profile_pic.png", group: "SE-32", course: "Software Engineering", year: "3 Bachelor (Current)", score: 98.07 },
  { id: "s3", name: "Lily Thompson", avatar: "/img/profile_pic.png", group: "SE-32", course: "Web Development", year: "3 Bachelor (Current)", score: 97.67, isMe: true },
  { id: "s4", name: "Ethan Parker", avatar: "/img/profile_pic.png", group: "SE-31", course: "Data Science", year: "2 Bachelor", score: 96.58 },
  { id: "s5", name: "Mason Carter", avatar: "/img/profile_pic.png", group: "SE-34", course: "Software Engineering", year: "4 Bachelor", score: 95.64 },
  { id: "s6", name: "Harper Williams", avatar: "/img/profile_pic.png", group: "SE-33", course: "Web Development", year: "1 Bachelor", score: 95.31 },
  { id: "s7", name: "Chloe Anderson", avatar: "/img/profile_pic.png", group: "SE-31", course: "Data Science", year: "3 Bachelor (Current)", score: 94.98 },
  // ще більше учнів
  { id: "s8", name: "Noah Reed", avatar: "/img/profile_pic.png", group: "SE-32", course: "Software Engineering", year: "3 Bachelor (Current)", score: 94.51 },
  { id: "s9", name: "Ava Brooks", avatar: "/img/profile_pic.png", group: "SE-33", course: "Data Science", year: "2 Bachelor", score: 94.11 },
  { id: "s10", name: "Lucas King", avatar: "/img/profile_pic.png", group: "SE-34", course: "Web Development", year: "4 Bachelor", score: 93.78 },
  { id: "s11", name: "Mia Rivera", avatar: "/img/profile_pic.png", group: "SE-31", course: "Software Engineering", year: "1 Bachelor", score: 93.40 },
  { id: "s12", name: "Oliver Scott", avatar: "/img/profile_pic.png", group: "SE-32", course: "Data Science", year: "3 Bachelor (Current)", score: 92.88 },
  { id: "s13", name: "Amelia Ward", avatar: "/img/profile_pic.png", group: "SE-34", course: "Web Development", year: "2 Bachelor", score: 92.63 },
  { id: "s14", name: "James Turner", avatar: "/img/profile_pic.png", group: "SE-33", course: "Software Engineering", year: "4 Bachelor", score: 92.21 },
  { id: "s15", name: "Sofia Price", avatar: "/img/profile_pic.png", group: "SE-31", course: "Data Science", year: "1 Bachelor", score: 91.74 },
  { id: "s16", name: "Henry Cole", avatar: "/img/profile_pic.png", group: "SE-32", course: "Web Development", year: "3 Bachelor (Current)", score: 91.20 },
  { id: "s17", name: "Grace Diaz", avatar: "/img/profile_pic.png", group: "SE-33", course: "Software Engineering", year: "2 Bachelor", score: 90.96 },
  { id: "s18", name: "Daniel Young", avatar: "/img/profile_pic.png", group: "SE-34", course: "Data Science", year: "4 Bachelor", score: 90.44 },
  { id: "s19", name: "Ella Hart", avatar: "/img/profile_pic.png", group: "SE-31", course: "Web Development", year: "1 Bachelor", score: 90.01 },
  { id: "s20", name: "Jack Morgan", avatar: "/img/profile_pic.png", group: "SE-32", course: "Software Engineering", year: "3 Bachelor (Current)", score: 89.70 },
];

const coursesBreakdownSeed: CourseBreakdown[] = [
  {
    id: "c1",
    title: "Introduction to Programming",
    markText: "5+",
    total: 97,
    blocks: [
      { label: "Assignments", value: "40/40" },
      { label: "Module 1", value: "8/10" },
      { label: "Module 2", value: "10/10" },
      { label: "Exam", value: "39/40" },
    ],
  },
  {
    id: "c2",
    title: "Data Science",
    markText: "5+",
    total: 98,
    blocks: [
      { label: "Assignments", value: "40/40" },
      { label: "Module", value: "18/20" },
      { label: "Exam", value: "40/40" },
    ],
  },
  {
    id: "c3",
    title: "Web Development",
    markText: "5+",
    total: 98,
    blocks: [
      { label: "Assignments", value: "40/40" },
      { label: "Module 1", value: "8/10" },
      { label: "Module 2", value: "10/10" },
      { label: "Exam", value: "40/40" },
    ],
  },
];

function Dropdown({
  label,
  value,
  items,
  open,
  onToggle,
  onSelect,
}: {
  label: string;
  value: string;
  items: readonly string[];
  open: boolean;
  onToggle: () => void;
  onSelect: (v: string) => void;
}) {
  const rootRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    function onDocClick(e: MouseEvent) {
      if (!open) return;
      const el = rootRef.current;
      if (!el) return;
      if (e.target instanceof Node && !el.contains(e.target)) onToggle(); // закрити
    }
    document.addEventListener("mousedown", onDocClick);
    return () => document.removeEventListener("mousedown", onDocClick);
  }, [open, onToggle]);

  return (
    <div className="grades-filter schedule-view-dd" ref={rootRef}>
      <label className="grades-filter-label">{label}</label>

      <div className={cx(open && "is-open")}>
        <button type="button" className="grades-filter-select schedule-view-select" onClick={onToggle}>
          <span className="grades-filter-select-span schedule-view-label">{value}</span>
          <span className="grades-filter-caret schedule-view-caret">
            <svg width="18" height="13" aria-hidden="true">
              <use href="/img/icons.svg#icon-vector-down" />
            </svg>
          </span>
        </button>

        <div className="dropdown-menu">
          {items.map((it) => (
            <button
              key={it}
              type="button"
              className={cx("dropdown-item", it === value && "is-active")}
              onClick={() => {
                onSelect(it);
                onToggle(); // закрити після вибору
              }}
            >
              {it}
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

export default function Grades() {
  const [courseFilter, setCourseFilter] = useState<(typeof COURSE_OPTIONS)[number]>("All courses");
  const [yearFilter, setYearFilter] = useState<(typeof YEAR_OPTIONS)[number]>("3 Bachelor (Current)");
  const [courseOpen, setCourseOpen] = useState(false);
  const [yearOpen, setYearOpen] = useState(false);

  const viewportRef = useRef<HTMLDivElement | null>(null);
  const [canPrev, setCanPrev] = useState(false);
  const [canNext, setCanNext] = useState(true);

  const filteredStudents = useMemo(() => {
    return studentsSeed.filter((s) => {
      const okCourse = courseFilter === "All courses" ? true : s.course === courseFilter;
      const okYear = yearFilter === "All years" ? true : s.year === yearFilter;
      return okCourse && okYear;
    });
  }, [courseFilter, yearFilter]);

  const ranked = useMemo(() => {
    // сортуємо по score DESC, перераховуємо місця
    const sorted = [...filteredStudents].sort((a, b) => b.score - a.score);
    return sorted.map((s, idx) => ({ ...s, place: idx + 1 }));
  }, [filteredStudents]);

  // коли фільтри змінюються — повертаємо скрол на початок і оновлюємо кнопки
  useEffect(() => {
    const el = viewportRef.current;
    if (!el) return;
    el.scrollLeft = 0;
    setCanPrev(false);
    setCanNext(el.scrollWidth > el.clientWidth + 2);
  }, [courseFilter, yearFilter]);

  function syncArrows() {
    const el = viewportRef.current;
    if (!el) return;
    const max = el.scrollWidth - el.clientWidth;
    setCanPrev(el.scrollLeft > 0);
    setCanNext(el.scrollLeft < max - 1);
  }

  function scrollByCards(dir: "prev" | "next") {
    const el = viewportRef.current;
    if (!el) return;

    // скрол приблизно на ширину видимої області - трохи
    const step = Math.max(200, Math.floor(el.clientWidth * 0.85));
    el.scrollBy({ left: dir === "next" ? step : -step, behavior: "smooth" });
  }

  // прості "Details" як було (можеш підв’язати під реальні дані)
  const overallGpa = ranked.find((s) => s.isMe)?.score ?? 97.67;
  const myPlace = ranked.find((s) => s.isMe)?.place ?? 3;

  return (
    <main className="sidebar-space grades-info">
      <h1 className="grades-title">Grades</h1>

      <div className="grades-content">
        {/* ===== Rating ===== */}
        <section className="grades-rating">
          <header className="grades-rating-header">
            <h2 className="grades-section-title">Rating</h2>

            <div className="grades-rating-filters">
              <Dropdown
                label="Course"
                value={courseFilter}
                items={COURSE_OPTIONS}
                open={courseOpen}
                onToggle={() => {
                  setCourseOpen((v) => !v);
                  setYearOpen(false);
                }}
                onSelect={(v) => setCourseFilter(v as any)}
              />

              <Dropdown
                label="Year"
                value={yearFilter}
                items={YEAR_OPTIONS}
                open={yearOpen}
                onToggle={() => {
                  setYearOpen((v) => !v);
                  setCourseOpen(false);
                }}
                onSelect={(v) => setYearFilter(v as any)}
              />
            </div>
          </header>

          {/* ===== Slider ===== */}
          <div className="grades-rating-slider grades-rating-slider--fixed">
            <button
              className="grades-swiper-btn grades-swiper-prev"
              type="button"
              aria-label="Previous"
              onClick={() => scrollByCards("prev")}
              disabled={!canPrev}
              style={{ opacity: canPrev ? 1 : 0.35, pointerEvents: canPrev ? "auto" : "none" }}
            >
              <svg width="20" height="36" aria-hidden="true">
                <use href="/img/icons.svg#icon-left2" />
              </svg>
            </button>

            <div className="grades-swiper swiper">
              <div
                className="grades-swiper-viewport"
                ref={viewportRef}
                onScroll={syncArrows}
              >
                <div className="swiper-wrapper">
                  {ranked.map((s) => {
                    const medalClass =
                      s.place === 1 ? "" : s.place === 2 ? "second-place" : s.place === 3 ? "third-place" : "";
                    const showMedal = s.place <= 3;
                    const cardClass = cx("rating-card", s.isMe ? "rating-card--current" : "");
                    const scoreClass = cx("rating-card-score", s.isMe && "rating-card-score--accent");

                    return (
                      <div className="swiper-slide" key={s.id}>
                        <article className={cardClass}>
                          <img src={s.avatar} alt={`${s.name} avatar`} className="rating-card-avatar" />

                          {showMedal && (
                            <svg width="20" height="26" className={cx("rating-card-medal", medalClass)} aria-hidden="true">
                              <use href="/img/icons.svg#icon-medal" />
                            </svg>
                          )}

                          <h3 className="rating-card-name">{s.name}</h3>
                          <p className="rating-card-group">{s.group}</p>
                          <p className="rating-card-position">{s.place}</p>
                        </article>

                        <p className={scoreClass}>{s.score.toFixed(2)}</p>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>

            <button
              className="grades-swiper-btn grades-swiper-next"
              type="button"
              aria-label="Next"
              onClick={() => scrollByCards("next")}
              disabled={!canNext}
              style={{ opacity: canNext ? 1 : 0.35, pointerEvents: canNext ? "auto" : "none" }}
            >
              <svg width="20" height="36" aria-hidden="true">
                <use href="/img/icons.svg#icon-right2" />
              </svg>
            </button>
          </div>
        </section>

        {/* ===== Details ===== */}
        <section className="grades-details">
          <h2 className="grades-section-title">Details</h2>

          <div className="grades-stats-grid">
            <article className="grades-stat-card">
              <p className="grades-stat-value grades-stat-value--blue">{overallGpa.toFixed(2)}</p>
              <p className="grades-stat-label">Overall GPA</p>
            </article>

            <article className="grades-stat-card">
              <p className="grades-stat-value grades-stat-value--green">5+</p>
              <p className="grades-stat-label">Performance</p>
            </article>

            <article className="grades-stat-card">
              <p className="grades-stat-value grades-stat-value--purple">{coursesBreakdownSeed.length}</p>
              <p className="grades-stat-label">Subjects</p>
            </article>

            <article className="grades-stat-card">
              <p className="grades-stat-value grades-stat-value--orange">{myPlace}</p>
              <p className="grades-stat-label">Rating</p>
            </article>
          </div>
        </section>

        {/* ===== Courses breakdown ===== */}
        <section className="grades-courses">
          {coursesBreakdownSeed.map((c) => (
            <article className="course-grade-card" key={c.id}>
              <header className="course-grade-header">
                <div className="course-grade-title-wrap">
                  <svg width="20" height="20" className="icon-medal-course" aria-hidden="true">
                    <use href="/img/icons.svg#icon-medal2" />
                  </svg>

                  <h3 className="course-grade-title">{c.title}</h3>
                  <span className="course-grade-mark course-grade-mark--green">{c.markText}</span>
                </div>
              </header>

              <div className="course-grade-body">
                <p className="course-grade-total">{c.total}</p>

                {c.blocks.map((b) => (
                  <div className="course-grade-block" key={b.label}>
                    <p className="course-grade-block-value">{b.value}</p>
                    <p className="course-grade-block-label">{b.label}</p>
                  </div>
                ))}
              </div>
            </article>
          ))}
        </section>
      </div>
    </main>
  );
}
