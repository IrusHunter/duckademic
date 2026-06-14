import { useEffect, useMemo, useRef, useState } from 'react'
import css from './App.module.css'

const PALETTE = [
  '#DB4437', '#E67E22', '#F4B400', '#0F9D58', '#16A085',
  '#4285F4', '#2980B9', '#8E44AD', '#9B59B6', '#E91E63',
  '#FF5722', '#795548', '#607D8B', '#009688',
]

function hashString(str: string): number {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash)
    hash |= 0
  }
  return Math.abs(hash)
}

function colorFor(name: string): string {
  return PALETTE[hashString(name || '?') % PALETTE.length]
}

function initialsFor(name: string): string {
  const parts = (name || '').trim().split(/\s+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return (name.slice(0, 1) || '?').toUpperCase()
}

function AvatarCircle({ name, size, className }: { name: string; size: number; className: string }) {
  return (
    <div
      className={className}
      style={{ backgroundColor: colorFor(name), fontSize: Math.round(size * 0.4) }}
    >
      {initialsFor(name)}
    </div>
  )
}

function getAuthRole(): string | null {
  try {
    const raw = localStorage.getItem('auth_user')
    return raw ? (JSON.parse(raw) as { role: string }).role : null
  } catch { return null }
}

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(' ')
}

type Student = {
  id: string
  name: string
  avatar: string
  group: string
  course: string
  year: string
  score: number
  isMe?: boolean
}

type CourseBreakdown = {
  id: string
  title: string
  markText: string
  total: number
  blocks: Array<{ label: string; value: string }>
}

const COURSE_OPTIONS = ['All courses', 'Software Engineering', 'Data Science', 'Web Development'] as const
const YEAR_OPTIONS = ['All years', '1 Bachelor', '2 Bachelor', '3 Bachelor (Current)', '4 Bachelor'] as const

const studentsSeed: Student[] = [
  { id: 's1',  name: 'Emily Johnson',   avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Software Engineering', year: '3 Bachelor (Current)', score: 99.81 },
  { id: 's2',  name: 'Sarah Wilson',    avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Software Engineering', year: '3 Bachelor (Current)', score: 98.07 },
  { id: 's3',  name: 'Lily Thompson',   avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Web Development',      year: '3 Bachelor (Current)', score: 97.67, isMe: true },
  { id: 's4',  name: 'Ethan Parker',    avatar: '/img/profile_pic.png', group: 'SE-31', course: 'Data Science',         year: '2 Bachelor',           score: 96.58 },
  { id: 's5',  name: 'Mason Carter',    avatar: '/img/profile_pic.png', group: 'SE-34', course: 'Software Engineering', year: '4 Bachelor',           score: 95.64 },
  { id: 's6',  name: 'Harper Williams', avatar: '/img/profile_pic.png', group: 'SE-33', course: 'Web Development',      year: '1 Bachelor',           score: 95.31 },
  { id: 's7',  name: 'Chloe Anderson',  avatar: '/img/profile_pic.png', group: 'SE-31', course: 'Data Science',         year: '3 Bachelor (Current)', score: 94.98 },
  { id: 's8',  name: 'Noah Reed',       avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Software Engineering', year: '3 Bachelor (Current)', score: 94.51 },
  { id: 's9',  name: 'Ava Brooks',      avatar: '/img/profile_pic.png', group: 'SE-33', course: 'Data Science',         year: '2 Bachelor',           score: 94.11 },
  { id: 's10', name: 'Lucas King',      avatar: '/img/profile_pic.png', group: 'SE-34', course: 'Web Development',      year: '4 Bachelor',           score: 93.78 },
  { id: 's11', name: 'Mia Rivera',      avatar: '/img/profile_pic.png', group: 'SE-31', course: 'Software Engineering', year: '1 Bachelor',           score: 93.40 },
  { id: 's12', name: 'Oliver Scott',    avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Data Science',         year: '3 Bachelor (Current)', score: 92.88 },
  { id: 's13', name: 'Amelia Ward',     avatar: '/img/profile_pic.png', group: 'SE-34', course: 'Web Development',      year: '2 Bachelor',           score: 92.63 },
  { id: 's14', name: 'James Turner',    avatar: '/img/profile_pic.png', group: 'SE-33', course: 'Software Engineering', year: '4 Bachelor',           score: 92.21 },
  { id: 's15', name: 'Sofia Price',     avatar: '/img/profile_pic.png', group: 'SE-31', course: 'Data Science',         year: '1 Bachelor',           score: 91.74 },
  { id: 's16', name: 'Henry Cole',      avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Web Development',      year: '3 Bachelor (Current)', score: 91.20 },
  { id: 's17', name: 'Grace Diaz',      avatar: '/img/profile_pic.png', group: 'SE-33', course: 'Software Engineering', year: '2 Bachelor',           score: 90.96 },
  { id: 's18', name: 'Daniel Young',    avatar: '/img/profile_pic.png', group: 'SE-34', course: 'Data Science',         year: '4 Bachelor',           score: 90.44 },
  { id: 's19', name: 'Ella Hart',       avatar: '/img/profile_pic.png', group: 'SE-31', course: 'Web Development',      year: '1 Bachelor',           score: 90.01 },
  { id: 's20', name: 'Jack Morgan',     avatar: '/img/profile_pic.png', group: 'SE-32', course: 'Software Engineering', year: '3 Bachelor (Current)', score: 89.70 },
]

const coursesBreakdownSeed: CourseBreakdown[] = [
  {
    id: 'c1',
    title: 'Introduction to Programming',
    markText: '5+',
    total: 97,
    blocks: [
      { label: 'Assignments', value: '40/40' },
      { label: 'Module 1',    value: '8/10'  },
      { label: 'Module 2',    value: '10/10' },
      { label: 'Exam',        value: '39/40' },
    ],
  },
  {
    id: 'c2',
    title: 'Data Science',
    markText: '5+',
    total: 98,
    blocks: [
      { label: 'Assignments', value: '40/40' },
      { label: 'Module',      value: '18/20' },
      { label: 'Exam',        value: '40/40' },
    ],
  },
  {
    id: 'c3',
    title: 'Web Development',
    markText: '5+',
    total: 98,
    blocks: [
      { label: 'Assignments', value: '40/40' },
      { label: 'Module 1',    value: '8/10'  },
      { label: 'Module 2',    value: '10/10' },
      { label: 'Exam',        value: '40/40' },
    ],
  },
]

function Dropdown({
  label,
  value,
  items,
  open,
  onToggle,
  onSelect,
}: {
  label: string
  value: string
  items: readonly string[]
  open: boolean
  onToggle: () => void
  onSelect: (v: string) => void
}) {
  const rootRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    function onDocClick(e: MouseEvent) {
      if (!open) return
      const el = rootRef.current
      if (!el) return
      if (e.target instanceof Node && !el.contains(e.target)) onToggle()
    }
    document.addEventListener('mousedown', onDocClick)
    return () => document.removeEventListener('mousedown', onDocClick)
  }, [open, onToggle])

  return (
    <div className={css.filter} ref={rootRef}>
      <label className={css.filterLabel}>{label}</label>

      <div className={cx(open && css.open)}>
        <button type="button" className={css.filterSelect} onClick={onToggle}>
          <span className={css.filterSelectSpan}>{value}</span>
          <span className={css.filterCaret}>
            <svg width="18" height="13" aria-hidden="true">
              <use href="/img/icons.svg#icon-vector-down" />
            </svg>
          </span>
        </button>

        <div className={css.dropdown}>
          {items.map((it) => (
            <button
              key={it}
              type="button"
              className={cx(css.dropdownItem, it === value && css.dropdownItemActive)}
              onClick={() => { onSelect(it); onToggle() }}
            >
              {it}
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

export default function App() {
  const isTeacher = getAuthRole() === 'teacher'

  const [courseFilter, setCourseFilter] = useState<(typeof COURSE_OPTIONS)[number]>('All courses')
  const [yearFilter, setYearFilter] = useState<(typeof YEAR_OPTIONS)[number]>('3 Bachelor (Current)')
  const [courseOpen, setCourseOpen] = useState(false)
  const [yearOpen, setYearOpen] = useState(false)

  const viewportRef = useRef<HTMLDivElement | null>(null)
  const [canPrev, setCanPrev] = useState(false)
  const [canNext, setCanNext] = useState(true)

  const filteredStudents = useMemo(() => {
    return studentsSeed.filter((s) => {
      const okCourse = courseFilter === 'All courses' || s.course === courseFilter
      const okYear   = yearFilter   === 'All years'   || s.year   === yearFilter
      return okCourse && okYear
    })
  }, [courseFilter, yearFilter])

  const ranked = useMemo(() => {
    return [...filteredStudents]
      .sort((a, b) => b.score - a.score)
      .map((s, idx) => ({ ...s, place: idx + 1 }))
  }, [filteredStudents])

  useEffect(() => {
    const el = viewportRef.current
    if (!el) return
    el.scrollLeft = 0
    setCanPrev(false)
    setCanNext(el.scrollWidth > el.clientWidth + 2)
  }, [courseFilter, yearFilter])

  function syncArrows() {
    const el = viewportRef.current
    if (!el) return
    const max = el.scrollWidth - el.clientWidth
    setCanPrev(el.scrollLeft > 0)
    setCanNext(el.scrollLeft < max - 1)
  }

  function scrollByCards(dir: 'prev' | 'next') {
    const el = viewportRef.current
    if (!el) return
    const step = Math.max(200, Math.floor(el.clientWidth * 0.85))
    el.scrollBy({ left: dir === 'next' ? step : -step, behavior: 'smooth' })
  }

  const overallGpa = isTeacher ? null : (ranked.find((s) => s.isMe)?.score ?? 97.67)
  const myPlace    = isTeacher ? null : (ranked.find((s) => s.isMe)?.place  ?? 3)

  return (
    <main className={css.page}>
      <h1 className={css.title}>Grades</h1>

      <div className={css.content}>

        {/* ===== Rating ===== */}
        <section className={css.rating}>
          <header className={css.ratingHeader}>
            <h2 className={css.sectionTitle}>Rating</h2>

            <div className={css.ratingFilters}>
              <Dropdown
                label="Course"
                value={courseFilter}
                items={COURSE_OPTIONS}
                open={courseOpen}
                onToggle={() => { setCourseOpen((v) => !v); setYearOpen(false) }}
                onSelect={(v) => setCourseFilter(v as typeof courseFilter)}
              />
              <Dropdown
                label="Year"
                value={yearFilter}
                items={YEAR_OPTIONS}
                open={yearOpen}
                onToggle={() => { setYearOpen((v) => !v); setCourseOpen(false) }}
                onSelect={(v) => setYearFilter(v as typeof yearFilter)}
              />
            </div>
          </header>

          {/* Slider */}
          <div className={css.ratingSlider}>
            <button
              className={css.swiperBtn}
              type="button"
              aria-label="Previous"
              onClick={() => scrollByCards('prev')}
              disabled={!canPrev}
              style={{ opacity: canPrev ? 1 : 0.35, pointerEvents: canPrev ? 'auto' : 'none' }}
            >
              <svg width="20" height="36" aria-hidden="true">
                <use href="/img/icons.svg#icon-left2" />
              </svg>
            </button>

            <div className={css.swiper}>
              <div
                className={css.swiperViewport}
                ref={viewportRef}
                onScroll={syncArrows}
              >
                <div className={css.swiperWrapper}>
                  {ranked.map((s) => {
                    const medalClass =
                      s.place === 1 ? '' :
                      s.place === 2 ? css.secondPlace :
                      s.place === 3 ? css.thirdPlace  : ''

                    return (
                      <div className={css.swiperSlide} key={s.id}>
                        <article className={cx(css.ratingCard, !isTeacher && s.isMe && css.ratingCardCurrent)}>
                          <AvatarCircle
                            name={s.name}
                            size={!isTeacher && s.isMe ? 80 : 72}
                            className={css.ratingCardAvatar}
                          />

                          {s.place <= 3 && (
                            <svg
                              width="20" height="26"
                              className={cx(css.ratingCardMedal, medalClass)}
                              aria-hidden="true"
                            >
                              <use href="/img/icons.svg#icon-medal" />
                            </svg>
                          )}

                          <h3 className={css.ratingCardName}>{s.name}</h3>
                          <p className={css.ratingCardGroup}>{s.group}</p>
                          <p className={css.ratingCardPosition}>{s.place}</p>
                        </article>

                        <p className={cx(css.ratingCardScore, !isTeacher && s.isMe && css.ratingCardScoreAccent)}>
                          {s.score.toFixed(2)}
                        </p>
                      </div>
                    )
                  })}
                </div>
              </div>
            </div>

            <button
              className={css.swiperBtn}
              type="button"
              aria-label="Next"
              onClick={() => scrollByCards('next')}
              disabled={!canNext}
              style={{ opacity: canNext ? 1 : 0.35, pointerEvents: canNext ? 'auto' : 'none' }}
            >
              <svg width="20" height="36" aria-hidden="true">
                <use href="/img/icons.svg#icon-right2" />
              </svg>
            </button>
          </div>
        </section>

        {/* ===== Details ===== */}
        <section className={css.details}>
          <h2 className={css.sectionTitle}>Details</h2>

          <div className={css.statsGrid}>
            <article className={css.statCard}>
              <p className={cx(css.statValue, css.statValueBlue)}>{overallGpa != null ? overallGpa.toFixed(2) : '—'}</p>
              <p className={css.statLabel}>Overall GPA</p>
            </article>
            <article className={css.statCard}>
              <p className={cx(css.statValue, css.statValueGreen)}>5+</p>
              <p className={css.statLabel}>Performance</p>
            </article>
            <article className={css.statCard}>
              <p className={cx(css.statValue, css.statValuePurple)}>{coursesBreakdownSeed.length}</p>
              <p className={css.statLabel}>Subjects</p>
            </article>
            <article className={css.statCard}>
              <p className={cx(css.statValue, css.statValueOrange)}>{myPlace ?? '—'}</p>
              <p className={css.statLabel}>Rating</p>
            </article>
          </div>
        </section>

        {/* ===== Courses breakdown ===== */}
        <section className={css.courses}>
          {coursesBreakdownSeed.map((c) => (
            <article className={css.courseGradeCard} key={c.id}>
              <header className={css.courseGradeHeader}>
                <div className={css.courseGradeTitleWrap}>
                  <svg width="20" height="20" aria-hidden="true">
                    <use href="/img/icons.svg#icon-medal2" />
                  </svg>
                  <h3 className={css.courseGradeTitle}>{c.title}</h3>
                  <span className={cx(css.courseGradeMark, css.courseGradeMarkGreen)}>{c.markText}</span>
                </div>
              </header>

              <div className={css.courseGradeBody}>
                <p className={css.courseGradeTotal}>{c.total}</p>

                {c.blocks.map((b) => (
                  <div className={css.courseGradeBlock} key={b.label}>
                    <p className={css.courseGradeBlockValue}>{b.value}</p>
                    <p className={css.courseGradeBlockLabel}>{b.label}</p>
                  </div>
                ))}
              </div>
            </article>
          ))}
        </section>

      </div>
    </main>
  )
}
