import { useEffect, useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { getPersonalSchedule, getAllPersonalSchedule } from '../../services/scheduleService'
import css from './App.module.css'

const DAY_NAMES = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
const VIEW_OPTIONS = ['General', 'Date'] as const
type ViewOption = (typeof VIEW_OPTIONS)[number]

const ROW_COLORS = [css.rowBlue, css.rowGreen, css.rowRose, css.rowAmber]

const LESSON_DURATION_MS = 80 * 60 * 1000
const DATE_RANGE_GAP_MS = 10 * 24 * 60 * 60 * 1000

function cx(...parts: Array<string | false | undefined | null>) {
  return parts.filter(Boolean).join(' ')
}

function todayIndex(): number {
  const d = new Date().getDay()
  return d === 0 ? 0 : d - 1
}

function utcTimeRange(iso: string): string {
  const start = new Date(iso)
  const end = new Date(start.getTime() + LESSON_DURATION_MS)
  const fmt = (d: Date) =>
    `${String(d.getUTCHours()).padStart(2, '0')}:${String(d.getUTCMinutes()).padStart(2, '0')}`
  return `${fmt(start)} – ${fmt(end)}`
}

function fmtDate(d: Date): string {
  return `${String(d.getUTCDate()).padStart(2, '0')}.${String(d.getUTCMonth() + 1).padStart(2, '0')}`
}

function computeDateRanges(isoStrings: string[]): string {
  if (isoStrings.length === 0) return ''
  const dates = isoStrings
    .map((s) => new Date(s))
    .sort((a, b) => a.getTime() - b.getTime())

  const ranges: [Date, Date][] = []
  let start = dates[0]
  let end = dates[0]

  for (let i = 1; i < dates.length; i++) {
    if (dates[i].getTime() - dates[i - 1].getTime() <= DATE_RANGE_GAP_MS) {
      end = dates[i]
    } else {
      ranges.push([start, end])
      start = dates[i]
      end = dates[i]
    }
  }
  ranges.push([start, end])

  return ranges
    .map(([s, e]) =>
      s.getTime() === e.getTime() ? fmtDate(s) : `${fmtDate(s)}–${fmtDate(e)}`,
    )
    .join(', ')
}

export default function App() {
  const [dayIndex, setDayIndex] = useState(todayIndex)
  const [view, setView] = useState<ViewOption>('General')
  const [viewOpen, setViewOpen] = useState(false)

  const schedule = useQuery({ queryKey: ['personal-schedule'], queryFn: getPersonalSchedule })
  const fullSchedule = useQuery({ queryKey: ['all-personal-schedule'], queryFn: getAllPersonalSchedule })

  // dayIndex 0=Mon … 5=Sat  →  getUTCDay() 1=Mon … 6=Sat
  const weekday = dayIndex + 1

  const dayLessons = useMemo(() => {
    return (schedule.data ?? [])
      .filter((o) => new Date(o.date).getUTCDay() === weekday)
      .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())
  }, [schedule.data, weekday])

  // study_load_id → formatted date ranges across full semester
  const dateRangesMap = useMemo(() => {
    const groups = new Map<string, string[]>()
    for (const occ of fullSchedule.data ?? []) {
      const arr = groups.get(occ.study_load_id) ?? []
      arr.push(occ.date)
      groups.set(occ.study_load_id, arr)
    }
    const result = new Map<string, string>()
    for (const [id, dates] of groups) {
      result.set(id, computeDateRanges(dates))
    }
    return result
  }, [fullSchedule.data])

  useEffect(() => {
    function onDocClick() { setViewOpen(false) }
    function onKeyDown(e: KeyboardEvent) { if (e.key === 'Escape') setViewOpen(false) }
    document.addEventListener('click', onDocClick)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('click', onDocClick)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [])

  const isDate = view === 'Date'

  return (
    <main className={css.page}>
      <section className={css.layout}>
        <h1 className={css.title}>Schedule</h1>

        <div className={css.card}>
          <header className={css.mainHeader}>
            <div className={css.mainLeft}>
              <svg width="20" height="20">
                <use href="/img/icons.svg#icon-SVG-11" />
              </svg>
              <span className={css.heading}>Today&apos;s Classes</span>
            </div>

            <div className={css.controls}>
              <div className={css.dayNavRow}>
                <button
                  className={css.navBtn}
                  type="button"
                  onClick={() => setDayIndex((i) => (i - 1 + DAY_NAMES.length) % DAY_NAMES.length)}
                >
                  <svg width="14" height="20">
                    <use href="/img/icons.svg#icon-vector-left" />
                  </svg>
                </button>

                <button className={css.dayBtn} type="button">
                  <span className={css.dayLabel}>{DAY_NAMES[dayIndex]}</span>
                </button>

                <button
                  className={css.navBtn}
                  type="button"
                  onClick={() => setDayIndex((i) => (i + 1) % DAY_NAMES.length)}
                >
                  <svg width="14" height="20">
                    <use href="/img/icons.svg#icon-vector-right" />
                  </svg>
                </button>
              </div>

              <div className={cx(css.viewDd, viewOpen && css.open)}>
                <button
                  className={css.viewSelect}
                  type="button"
                  aria-expanded={viewOpen}
                  onClick={(e) => {
                    e.stopPropagation()
                    setViewOpen((v) => !v)
                  }}
                >
                  <span className={css.viewLabel}>{view}</span>
                  <span className={css.viewCaret}>
                    <svg width="18" height="13">
                      <use href="/img/icons.svg#icon-vector-down" />
                    </svg>
                  </span>
                </button>

                <div className={css.dropdown}>
                  {VIEW_OPTIONS.map((opt) => (
                    <button
                      key={opt}
                      className={cx(css.dropdownItem, opt === view && css.dropdownItemActive)}
                      type="button"
                      onClick={() => {
                        setView(opt)
                        setViewOpen(false)
                      }}
                    >
                      {opt}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          </header>

          <ul className={css.classes}>
            {schedule.isLoading && (
              <li className={css.classRow}>
                <div className={css.classTime} />
                <article className={css.classCard}>
                  <p className={css.classTeacher}>Loading…</p>
                </article>
              </li>
            )}

            {schedule.isError && (
              <li className={css.classRow}>
                <div className={css.classTime} />
                <article className={css.classCard}>
                  <p className={css.classTeacher}>Failed to load schedule.</p>
                </article>
              </li>
            )}

            {!schedule.isLoading && !schedule.isError && dayLessons.length === 0 && (
              <li className={css.classRow}>
                <div className={css.classTime} />
                <article className={css.classCard}>
                  <p className={css.classTeacher}>No classes on this day.</p>
                </article>
              </li>
            )}

            {dayLessons.map((occ, idx) => {
              const colorClass = ROW_COLORS[idx % ROW_COLORS.length]
              const dateRanges = dateRangesMap.get(occ.study_load_id)

              return (
                <li key={occ.id} className={cx(css.classRow, colorClass)}>
                  <div className={css.classTime}>{utcTimeRange(occ.date)}</div>

                  <article className={css.classCard}>
                    <header className={css.classCardHeader}>
                      <h2 className={css.classTitle}>
                        {occ.study_load?.discipline_name ?? '—'}
                      </h2>
                      {occ.study_load?.lesson_type_name && (
                        <span className={css.classTag}>{occ.study_load.lesson_type_name}</span>
                      )}
                    </header>

                    {occ.study_load?.teacher_name && (
                      <p className={css.classTeacher}>{occ.study_load.teacher_name}</p>
                    )}

                    <div className={css.classMeta}>
                      {occ.classroom && (
                        <div className={css.classMetaItem}>
                          Room {occ.classroom.number}
                        </div>
                      )}
                      {occ.study_load?.student_group_name && (
                        <div className={css.classMetaItem}>
                          {occ.study_load.student_group_name}
                        </div>
                      )}
                      {isDate && dateRanges && (
                        <div className={css.classMetaDates}>{dateRanges}</div>
                      )}
                    </div>
                  </article>
                </li>
              )
            })}
          </ul>
        </div>
      </section>
    </main>
  )
}
