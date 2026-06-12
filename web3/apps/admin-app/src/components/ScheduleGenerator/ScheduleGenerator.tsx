import { useState } from 'react'
import {
  loadDataIntoGenerator,
  loadClassroomsIntoGenerator,
  submitAndGo,
  processStep,
  extractDataFromGenerator,
} from '../../api/generatorApi'
import css from './ScheduleGenerator.module.css'

// Редагований пайплайн (за docs/schedule_generator_pipeline.md).
// Кожен крок: дозволені методи (алгоритми) + чи можна вимкнути.
interface PipelineStep {
  id: string
  label: string
  methods: string[]   // дозволені алгоритми; порожньо = фіксований метод ("")
  optional?: boolean  // крок можна вимкнути
  fixed?: boolean      // метод не редагується (Expansion)
}

const PIPELINE: PipelineStep[] = [
  { id: 'weekday', label: '2 · Weekday Allocation', methods: ['even_weekday_allocator'] },
  { id: 'weekly-slot', label: '3 · Weekly Time Slot Assignment', methods: ['one_per_week_time_slot_assigner', 'brute_time_slot_assigner'] },
  { id: 'weekly-class', label: '4 · Weekly Classroom Assignment', methods: ['munkres_classroom_assigner'], optional: true },
  { id: 'expansion', label: '5 · Weekly Schedule Expansion', methods: [], fixed: true },
  { id: 'full-slot', label: '6 · Full Time Slot Assignment', methods: ['one_per_week_time_slot_assigner', 'brute_time_slot_assigner'] },
  { id: 'full-class', label: '7 · Full Classroom Assignment', methods: ['munkres_classroom_assigner'], optional: true },
]

interface StepState {
  enabled: boolean
  method: string
}

interface LogLine {
  text: string
  status: 'run' | 'ok' | 'fail'
}

function initialSteps(): Record<string, StepState> {
  const s: Record<string, StepState> = {}
  for (const step of PIPELINE) {
    s[step.id] = { enabled: true, method: step.methods[0] ?? '' }
  }
  return s
}

function parseIds(raw: string): string[] {
  return raw.split(/[\s,]+/).map(x => x.trim()).filter(Boolean)
}

export function ScheduleGenerator({ onClose }: { onClose: () => void }) {
  const [semesters, setSemesters] = useState('')
  const [classrooms, setClassrooms] = useState('')
  const [startTime, setStartTime] = useState('')
  const [steps, setSteps] = useState<Record<string, StepState>>(initialSteps)
  const [log, setLog] = useState<LogLine[]>([])
  const [running, setRunning] = useState(false)

  const append = (text: string, status: LogLine['status']) =>
    setLog(prev => [...prev, { text, status }])

  const setStep = (id: string, patch: Partial<StepState>) =>
    setSteps(prev => ({ ...prev, [id]: { ...prev[id], ...patch } }))

  async function run<T>(label: string, fn: () => Promise<T>): Promise<boolean> {
    append(label, 'run')
    try {
      await fn()
      setLog(prev => {
        const next = [...prev]
        next[next.length - 1] = { text: label, status: 'ok' }
        return next
      })
      return true
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : String(e)
      setLog(prev => {
        const next = [...prev]
        next[next.length - 1] = { text: `${label} — ${msg}`, status: 'fail' }
        return next
      })
      return false
    }
  }

  async function handleGenerate() {
    setLog([])
    setRunning(true)

    const semesterIds = parseIds(semesters)
    const classroomIds = parseIds(classrooms)

    try {
      // 1. Setup — завантажуємо дані з БД у генератор
      if (semesterIds.length === 0) {
        append('Provide at least one curriculum semester ID', 'fail')
        return
      }
      if (!(await run(`Load data (${semesterIds.length} semester(s))`, () => loadDataIntoGenerator(semesterIds)))) return

      if (classroomIds.length > 0) {
        if (!(await run(`Load classrooms (${classroomIds.length})`, () => loadClassroomsIntoGenerator(classroomIds)))) return
      }

      // 2. Вхід у перший крок
      if (!(await run('Submit & go (enter pipeline)', () => submitAndGo()))) return

      // 3. Кроки пайплайну: process-step(method) → submit-and-go
      for (const step of PIPELINE) {
        const st = steps[step.id]
        // Опціональний крок можна вимкнути; обов'язкові завжди виконуються
        if (step.optional && !st.enabled) {
          append(`Skip ${step.label}`, 'ok')
          continue
        }
        const method = step.fixed ? '' : st.method
        const label = step.fixed ? `${step.label}` : `${step.label} · ${method}`
        if (!(await run(label, () => processStep(method)))) return
        if (!(await run('Submit & go', () => submitAndGo()))) return
      }

      // 4. Extraction — записуємо розклад у БД
      if (!startTime) {
        append('Provide a start time to extract the schedule', 'fail')
        return
      }
      const iso = new Date(startTime).toISOString()
      if (!(await run('Extract schedule into database', () => extractDataFromGenerator(iso)))) return

      append('Done ✓ Schedule generated', 'ok')
    } finally {
      setRunning(false)
    }
  }

  return (
    <div className={css.overlay} onClick={onClose}>
      <div className={css.modal} onClick={e => e.stopPropagation()}>
        <div className={css.header}>
          <h2 className={css.title}>Schedule Generation</h2>
          <button className={css.close} onClick={onClose} aria-label="Close">×</button>
        </div>

        <div className={css.body}>
          {/* Inputs */}
          <label className={css.field}>
            <span className={css.label}>Curriculum semester IDs</span>
            <textarea
              className={css.input}
              rows={2}
              placeholder="uuid, uuid, …"
              value={semesters}
              onChange={e => setSemesters(e.target.value)}
            />
          </label>

          <label className={css.field}>
            <span className={css.label}>Classroom IDs (optional)</span>
            <textarea
              className={css.input}
              rows={2}
              placeholder="uuid, uuid, …"
              value={classrooms}
              onChange={e => setClassrooms(e.target.value)}
            />
          </label>

          <label className={css.field}>
            <span className={css.label}>Schedule start time</span>
            <input
              className={css.input}
              type="datetime-local"
              value={startTime}
              onChange={e => setStartTime(e.target.value)}
            />
          </label>

          {/* Editable pipeline (the algorithm) */}
          <div className={css.pipeline}>
            <p className={css.label}>Pipeline (algorithm)</p>
            {PIPELINE.map(step => {
              const st = steps[step.id]
              return (
                <div key={step.id} className={css.step}>
                  <label className={css.stepToggle}>
                    <input
                      type="checkbox"
                      checked={st.enabled}
                      disabled={!step.optional}
                      onChange={e => setStep(step.id, { enabled: e.target.checked })}
                    />
                    <span>{step.label}</span>
                  </label>

                  {step.fixed ? (
                    <span className={css.fixedTag}>fixed</span>
                  ) : (
                    <select
                      className={css.select}
                      value={st.method}
                      disabled={!st.enabled}
                      onChange={e => setStep(step.id, { method: e.target.value })}
                    >
                      {step.methods.map(m => <option key={m} value={m}>{m}</option>)}
                    </select>
                  )}
                </div>
              )
            })}
          </div>

          {/* Live log */}
          {log.length > 0 && (
            <div className={css.log}>
              {log.map((line, i) => (
                <div key={i} className={`${css.logLine} ${css[line.status]}`}>
                  {line.status === 'ok' ? '✓' : line.status === 'fail' ? '✕' : '…'} {line.text}
                </div>
              ))}
            </div>
          )}
        </div>

        <div className={css.footer}>
          <button className={css.cancel} onClick={onClose} disabled={running}>Close</button>
          <button className={css.generate} onClick={handleGenerate} disabled={running}>
            {running ? 'Generating…' : 'Generate'}
          </button>
        </div>
      </div>
    </div>
  )
}
