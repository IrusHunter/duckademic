import { useEffect, useState } from 'react'
import {
  loadDataIntoGenerator,
  loadClassroomsIntoGenerator,
  submitAndGo,
  processStep,
  extractDataFromGenerator,
  fetchSemesters,
  fetchClassrooms,
  initGenerator,
  type Semester,
  type Classroom,
} from '../../api/generatorApi'
import css from './ScheduleGenerator.module.css'

interface PipelineStep {
  id: string
  label: string
  methods: string[]
  optional?: boolean
  fixed?: boolean
  resultType: 'weekday' | 'slot' | 'classroom' | 'none'
}

const PIPELINE: PipelineStep[] = [
  { id: 'weekday',      label: 'Weekday Allocation',           methods: ['even_weekday_allocator'],                                      resultType: 'weekday'   },
  { id: 'weekly-slot',  label: 'Weekly Time Slot Assignment',  methods: ['one_per_week_time_slot_assigner', 'brute_time_slot_assigner'],  resultType: 'slot'      },
  { id: 'weekly-class', label: 'Weekly Classroom Assignment',  methods: ['munkres_classroom_assigner'],  optional: true,                  resultType: 'classroom' },
  { id: 'expansion',    label: 'Weekly Schedule Expansion',    methods: [],                              fixed: true,                     resultType: 'none'      },
  { id: 'full-slot',    label: 'Full Time Slot Assignment',    methods: ['one_per_week_time_slot_assigner', 'brute_time_slot_assigner'],  resultType: 'slot'      },
  { id: 'full-class',   label: 'Full Classroom Assignment',    methods: ['munkres_classroom_assigner'],  optional: true,                  resultType: 'classroom' },
]

type Phase = 'setup' | 'pipeline' | 'export' | 'done'
type StepSub = 'idle' | 'processing' | 'processed' | 'advancing'
type ExportSub = 'idle' | 'busy' | 'done'

interface StepConfig { enabled: boolean; method: string }
interface LogLine { text: string; kind: 'run' | 'ok' | 'warn' | 'error' | 'info' }

function initConfigs(): Record<string, StepConfig> {
  return Object.fromEntries(PIPELINE.map(s => [s.id, { enabled: true, method: s.methods[0] ?? '' }]))
}


function isEnabled(step: PipelineStep, configs: Record<string, StepConfig>): boolean {
  return !step.optional || configs[step.id].enabled
}

function nextEnabledIdx(from: number, configs: Record<string, StepConfig>): number | null {
  for (let i = from + 1; i < PIPELINE.length; i++) {
    if (isEnabled(PIPELINE[i], configs)) return i
  }
  return null
}

// ── Result display ────────────────────────────────────────────────────────────

type ResultData = Record<string, unknown>

function WeekdayResult({ d }: { d: ResultData }) {
  const groups = (d.groups as unknown[]) ?? []
  const errors = (d.errors as unknown[]) ?? []
  return (
    <div className={css.result}>
      <ResultRow label="Groups assigned" value={groups.length} />
      {errors.length > 0 && <ResultRow label="Unaccounted slots" value={`${errors.length} issue(s)`} warn />}
    </div>
  )
}

function SlotResult({ d }: { d: ResultData }) {
  const assigned = (d.assigned as unknown[]) ?? []
  const unplaced = (d.unplaced as unknown[]) ?? []
  return (
    <div className={css.result}>
      <ResultRow label="Slots assigned" value={assigned.length} />
      {unplaced.length > 0 && <ResultRow label="Unplaced lessons" value={unplaced.length} warn />}
    </div>
  )
}

function ClassroomResult({ d }: { d: ResultData }) {
  const assigned = (d.assigned as unknown[]) ?? []
  const unassigned = (d.unassigned as unknown[]) ?? []
  return (
    <div className={css.result}>
      <ResultRow label="Classrooms assigned" value={assigned.length} />
      {unassigned.length > 0 && <ResultRow label="Without classroom" value={unassigned.length} warn />}
    </div>
  )
}

function ResultRow({ label, value, warn }: { label: string; value: number | string; warn?: boolean }) {
  return (
    <div className={css.resultRow}>
      <span className={css.resultLabel}>{label}</span>
      <span className={`${css.resultValue} ${warn ? css.resultWarn : ''}`}>{value}</span>
    </div>
  )
}

function StepResult({ data, type }: { data: unknown; type: PipelineStep['resultType'] }) {
  if (!data || type === 'none') return null
  const d = data as ResultData
  if (type === 'weekday')   return <WeekdayResult d={d} />
  if (type === 'slot')      return <SlotResult d={d} />
  if (type === 'classroom') return <ClassroomResult d={d} />
  return null
}

const WEEKDAYS = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday']
const DEFAULT_SLOT_PREFERENCE: number[][] = [
  [],
  [1, 1, 1, 1, 1, 1, 1],
  [1, 1, 1, 1, 1, 1, 1],
  [1, 1, 1, 1, 1, 1, 1],
  [1, 1, 1, 1, 1, 1, 1],
  [1, 1, 1, 1, 1, 1, 1],
  [],
]
const SLOTS_PER_DAY = 7

// ── Main component ────────────────────────────────────────────────────────────

export function ScheduleGenerator({ onClose }: { onClose: () => void }) {
  const [semesterList,         setSemesterList]         = useState<Semester[]>([])
  const [classroomList,        setClassroomList]        = useState<Classroom[]>([])
  const [loadingOptions,       setLoadingOptions]       = useState(true)
  const [selectedSemesterIds,  setSelectedSemesterIds]  = useState<string[]>([])
  const [selectedClassroomIds, setSelectedClassroomIds] = useState<string[]>([])
  const [genStartDate,         setGenStartDate]         = useState('')
  const [genEndDate,           setGenEndDate]           = useState('')
  const [slotPreference,       setSlotPreference]       = useState<number[][]>(DEFAULT_SLOT_PREFERENCE)
  const [maxDailyLoad,         setMaxDailyLoad]         = useState(4)
  const [lessonFillRate,       setLessonFillRate]       = useState(0.8)
  const [classroomOccupancy,   setClassroomOccupancy]   = useState(0.8)
  const [startTime,            setStartTime]            = useState('')
  const [configs,              setConfigs]              = useState<Record<string, StepConfig>>(initConfigs)

  const [phase,       setPhase]       = useState<Phase>('setup')
  const [stepIdx,     setStepIdx]     = useState(0)
  const [sub,         setSub]         = useState<StepSub>('idle')
  const [stepResult,  setStepResult]  = useState<unknown>(null)
  const [exportSub,   setExportSub]   = useState<ExportSub>('idle')
  const [setupError,  setSetupError]  = useState<string | null>(null)
  const [log,         setLog]         = useState<LogLine[]>([])

  useEffect(() => {
    Promise.all([fetchSemesters(), fetchClassrooms()])
      .then(([semRes, clsRes]) => {
        setSemesterList(semRes.data ?? [])
        setClassroomList(clsRes.data ?? [])
      })
      .catch(() => {})
      .finally(() => setLoadingOptions(false))
  }, [])

  const addLog = (text: string, kind: LogLine['kind']) =>
    setLog(prev => [...prev, { text, kind }])

  const patchConfig = (id: string, patch: Partial<StepConfig>) =>
    setConfigs(prev => ({ ...prev, [id]: { ...prev[id], ...patch } }))

  const busy = sub === 'processing' || sub === 'advancing' || exportSub === 'busy'

  // ── Setup → Pipeline ───────────────────────────────────────────────────────

  async function handleStart() {
    if (selectedSemesterIds.length === 0)         { setSetupError('Select at least one curriculum semester'); return }
    if (!genStartDate)                             { setSetupError('Provide a generator start date'); return }
    if (!genEndDate)                               { setSetupError('Provide a generator end date'); return }
    if (new Date(genEndDate) <= new Date(genStartDate)) { setSetupError('End date must be after start date'); return }
    if (slotPreference.every(d => d.length === 0)) { setSetupError('Select at least one active day'); return }
    if (!startTime)                                { setSetupError('Provide a schedule start time'); return }
    setSetupError(null)
    setSub('advancing')

    try {
      addLog('Initializing generator…', 'run')
      await initGenerator({
        start_date: new Date(genStartDate).toISOString(),
        end_date: new Date(genEndDate).toISOString(),
        slot_preference: slotPreference,
        max_daily_student_load: maxDailyLoad,
        lesson_fill_rate: lessonFillRate,
        classroom_occupancy: classroomOccupancy,
      })
      addLog('Generator initialized ✓', 'ok')

      addLog(`Loading data (${selectedSemesterIds.length} semester(s))…`, 'run')
      await loadDataIntoGenerator(selectedSemesterIds)
      addLog('Data loaded ✓', 'ok')

      const clsIds = selectedClassroomIds
      if (clsIds.length > 0) {
        addLog(`Loading classrooms (${clsIds.length})…`, 'run')
        await loadClassroomsIntoGenerator(clsIds)
        addLog('Classrooms loaded ✓', 'ok')
      }

      addLog('Entering pipeline…', 'run')
      await submitAndGo(true)
      addLog('Pipeline started ✓', 'ok')

      setStepIdx(0)
      setSub('idle')
      setStepResult(null)
      setPhase('pipeline')
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error: ${msg}`, 'error')
      setSetupError(msg)
      setSub('idle')
    }
  }

  // ── Process current step ───────────────────────────────────────────────────

  async function handleProcess() {
    const step = PIPELINE[stepIdx]
    const method = step.fixed ? '' : configs[step.id].method
    setSub('processing')
    setStepResult(null)

    try {
      addLog(`Processing: ${step.label}${method ? ` (${method})` : ''}…`, 'run')
      const res = await processStep(method)
      const data = (res as { data: unknown }).data
      setStepResult(data)
      addLog(`${step.label} ✓`, 'ok')
      setSub('processed')
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error: ${msg}`, 'error')
      setSub('idle')
    }
  }

  // ── Advance to next step ───────────────────────────────────────────────────

  async function handleNext() {
    setSub('advancing')

    try {
      addLog('Submit & go…', 'run')
      await submitAndGo(true)
      addLog('Step submitted ✓', 'ok')
      await advanceStepIdx()
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error: ${msg}`, 'error')
      setSub('processed')
    }
  }

  async function advanceStepIdx() {
    const next = nextEnabledIdx(stepIdx, configs)

    // For each disabled optional step between current and next, advance backend state
    if (next !== null) {
      for (let i = stepIdx + 1; i < next; i++) {
        addLog(`Skip: ${PIPELINE[i].label}`, 'info')
        await submitAndGo(true)
      }
      setStepIdx(next)
      setSub('idle')
      setStepResult(null)
    } else {
      setPhase('export')
      setExportSub('idle')
      addLog('All steps done. Ready to export.', 'ok')
    }
  }

  // ── Skip optional step ─────────────────────────────────────────────────────

  async function handleSkip() {
    setSub('advancing')
    try {
      addLog(`Skip: ${PIPELINE[stepIdx].label}`, 'info')
      await submitAndGo(true)
      await advanceStepIdx()
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error skipping: ${msg}`, 'error')
      setSub('idle')
    }
  }

  // ── Export ─────────────────────────────────────────────────────────────────

  async function handleExport() {
    setExportSub('busy')
    try {
      addLog('Extracting data from generator…', 'run')
      const iso = new Date(startTime).toISOString()
      await extractDataFromGenerator(iso)
      addLog('Data extracted ✓', 'ok')
      addLog('Schedule generated successfully ✓', 'ok')
      setExportSub('done')
      setPhase('done')
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error: ${msg}`, 'error')
      setExportSub('idle')
    }
  }

  // ── Render ─────────────────────────────────────────────────────────────────

  const currentStep = phase === 'pipeline' ? PIPELINE[stepIdx] : null

  return (
    <div className={css.overlay} onClick={onClose}>
      <div className={css.modal} onClick={e => e.stopPropagation()}>

        {/* Header */}
        <div className={css.header}>
          <div>
            <h2 className={css.title}>Schedule Generation</h2>
            {phase === 'pipeline' && (
              <p className={css.stepProgress}>Step {stepIdx + 1} of {PIPELINE.length}</p>
            )}
          </div>
          <button className={css.close} onClick={onClose} aria-label="Close" disabled={busy}>×</button>
        </div>

        {/* Body */}
        <div className={css.body}>

          {/* ── Setup ── */}
          {phase === 'setup' && (
            <>
              <div className={css.field}>
                <span className={css.label}>Curriculum semester IDs</span>
                {loadingOptions
                  ? <span className={css.optHint}>Loading…</span>
                  : semesterList.length === 0
                    ? <span className={css.optHint}>No semesters found</span>
                    : (
                      <div className={css.multiSelect}>
                        {semesterList.map(s => (
                          <label key={s.id} className={css.optionRow}>
                            <input
                              type="checkbox"
                              checked={selectedSemesterIds.includes(s.id)}
                              onChange={e =>
                                setSelectedSemesterIds(prev =>
                                  e.target.checked ? [...prev, s.id] : prev.filter(id => id !== s.id)
                                )
                              }
                            />
                            <span>Semester {s.number} — {s.slug}</span>
                          </label>
                        ))}
                      </div>
                    )
                }
              </div>

              <div className={css.field}>
                <div className={css.labelRow}>
                  <span className={css.label}>Classroom IDs (optional)</span>
                  {!loadingOptions && classroomList.length > 0 && (
                    <button type="button" className={css.selectAllBtn}
                      onClick={() => setSelectedClassroomIds(
                        selectedClassroomIds.length === classroomList.length
                          ? []
                          : classroomList.map(c => c.id)
                      )}>
                      {selectedClassroomIds.length === classroomList.length ? 'Deselect all' : 'Select all'}
                    </button>
                  )}
                </div>
                {loadingOptions
                  ? <span className={css.optHint}>Loading…</span>
                  : classroomList.length === 0
                    ? <span className={css.optHint}>No classrooms found</span>
                    : (
                      <div className={css.multiSelect}>
                        {classroomList.map(c => (
                          <label key={c.id} className={css.optionRow}>
                            <input
                              type="checkbox"
                              checked={selectedClassroomIds.includes(c.id)}
                              onChange={e =>
                                setSelectedClassroomIds(prev =>
                                  e.target.checked ? [...prev, c.id] : prev.filter(id => id !== c.id)
                                )
                              }
                            />
                            <span>Room {c.number} (cap. {c.capacity})</span>
                          </label>
                        ))}
                      </div>
                    )
                }
              </div>

              <div className={css.configSection}>
                <p className={css.label}>Generator configuration</p>

                <div className={css.dateRow}>
                  <label className={css.field}>
                    <span className={css.label}>Start date</span>
                    <input className={css.input} type="datetime-local"
                      value={genStartDate} onChange={e => setGenStartDate(e.target.value)} />
                  </label>
                  <label className={css.field}>
                    <span className={css.label}>End date</span>
                    <input className={css.input} type="datetime-local"
                      value={genEndDate} onChange={e => setGenEndDate(e.target.value)} />
                  </label>
                </div>

                <div className={css.field}>
                  <span className={css.label}>Slot preferences by day</span>
                  <div className={css.slotTable}>
                    <div className={css.slotHeader}>
                      <span className={css.slotDayCol} />
                      {Array.from({ length: SLOTS_PER_DAY }, (_, i) => (
                        <span key={i} className={css.slotNumCol}>S{i + 1}</span>
                      ))}
                    </div>
                    {WEEKDAYS.map((day, dayIdx) => {
                      const active = slotPreference[dayIdx].length > 0
                      return (
                        <div key={day} className={`${css.slotRow} ${active ? css.slotRowActive : ''}`}>
                          <label className={css.slotDayCol}>
                            <input
                              type="checkbox"
                              checked={active}
                              onChange={e => setSlotPreference(prev => {
                                const next = prev.map(d => [...d])
                                next[dayIdx] = e.target.checked ? Array<number>(SLOTS_PER_DAY).fill(1) : []
                                return next
                              })}
                            />
                            <span>{day.slice(0, 3)}</span>
                          </label>
                          {Array.from({ length: SLOTS_PER_DAY }, (_, slotIdx) => (
                            <input
                              key={slotIdx}
                              className={css.slotInput}
                              type="number"
                              min={0.1}
                              max={9.9}
                              step={0.1}
                              disabled={!active}
                              value={active ? (slotPreference[dayIdx][slotIdx] ?? 1) : 1}
                              onChange={e => setSlotPreference(prev => {
                                const next = prev.map(d => [...d])
                                next[dayIdx][slotIdx] = Number(e.target.value)
                                return next
                              })}
                            />
                          ))}
                        </div>
                      )
                    })}
                  </div>
                </div>

                <div className={css.numRow}>
                  <label className={css.field}>
                    <span className={css.label}>Max daily student load</span>
                    <input className={css.input} type="number" min={1} max={20}
                      value={maxDailyLoad} onChange={e => setMaxDailyLoad(Number(e.target.value))} />
                  </label>
                  <label className={css.field}>
                    <span className={css.label}>Lesson fill rate</span>
                    <input className={css.input} type="number" min={0.01} max={1} step={0.01}
                      value={lessonFillRate} onChange={e => setLessonFillRate(Number(e.target.value))} />
                  </label>
                  <label className={css.field}>
                    <span className={css.label}>Classroom occupancy</span>
                    <input className={css.input} type="number" min={0.01} max={1} step={0.01}
                      value={classroomOccupancy} onChange={e => setClassroomOccupancy(Number(e.target.value))} />
                  </label>
                </div>
              </div>

              <label className={css.field}>
                <span className={css.label}>Schedule start time</span>
                <input className={css.input} type="datetime-local"
                  value={startTime} onChange={e => setStartTime(e.target.value)} />
              </label>

              <div className={css.pipeline}>
                <p className={css.label}>Pipeline (algorithm)</p>
                {PIPELINE.map(step => {
                  const cfg = configs[step.id]
                  return (
                    <div key={step.id} className={css.step}>
                      <label className={css.stepToggle}>
                        <input type="checkbox" checked={cfg.enabled}
                          disabled={!step.optional}
                          onChange={e => patchConfig(step.id, { enabled: e.target.checked })} />
                        <span>{step.label}</span>
                      </label>
                      {step.fixed
                        ? <span className={css.fixedTag}>fixed</span>
                        : (
                          <select className={css.select} value={cfg.method}
                            disabled={!cfg.enabled}
                            onChange={e => patchConfig(step.id, { method: e.target.value })}>
                            {step.methods.map(m => <option key={m} value={m}>{m}</option>)}
                          </select>
                        )
                      }
                    </div>
                  )
                })}
              </div>

              {setupError && <p className={css.errorMsg}>{setupError}</p>}
            </>
          )}

          {/* ── Pipeline ── */}
          {phase === 'pipeline' && currentStep && (
            <div className={css.pipelinePhase}>
              <div className={css.stepCard}>
                <div className={css.stepCardHeader}>
                  <span className={css.stepBadge}>{stepIdx + 1}</span>
                  <div>
                    <p className={css.stepName}>{currentStep.label}</p>
                    {!currentStep.fixed && configs[currentStep.id].method && (
                      <p className={css.stepMethod}>{configs[currentStep.id].method}</p>
                    )}
                  </div>
                </div>
              </div>

              {stepResult != null && (
                <StepResult data={stepResult} type={currentStep.resultType} />
              )}

            </div>
          )}

          {/* ── Export ── */}
          {phase === 'export' && (
            <div className={css.exportPhase}>
              <p className={css.exportTitle}>Export Results</p>
              <p className={css.exportHint}>Extract the generated schedule data into the database.</p>
              <div className={css.exportSteps}>
                <div className={`${css.exportStep} ${exportSub === 'done' ? css.exportStepDone : ''}`}>
                  <span className={css.exportStepNum}>1</span>
                  <span>Extract data from generator</span>
                  {exportSub === 'done' && <span className={css.exportCheck}>✓</span>}
                </div>
              </div>
            </div>
          )}

          {/* ── Done ── */}
          {phase === 'done' && (
            <div className={css.donePhase}>
              <div className={css.doneIcon}>✓</div>
              <p className={css.doneText}>Schedule generated successfully</p>
            </div>
          )}

          {/* Log */}
          {log.length > 0 && (
            <div className={css.log}>
              {log.map((line, i) => (
                <div key={i} className={`${css.logLine} ${css[line.kind]}`}>
                  {line.kind === 'ok' ? '✓' : line.kind === 'error' ? '✕' : line.kind === 'warn' ? '⚠' : line.kind === 'run' ? '…' : '·'} {line.text}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className={css.footer}>
          <button className={css.cancel} onClick={onClose} disabled={busy}>
            {phase === 'done' ? 'Close' : 'Cancel'}
          </button>

          {phase === 'setup' && (
            <button className={css.generate} onClick={handleStart} disabled={busy}>
              {sub === 'advancing' ? 'Loading…' : 'Start'}
            </button>
          )}

          {phase === 'pipeline' && (
            <>
              {currentStep?.optional && (sub === 'idle' || sub === 'processed') && (
                <button className={css.secondary} onClick={handleSkip} disabled={busy}>Skip</button>
              )}

              {sub === 'idle' && (
                <button className={css.generate} onClick={handleProcess}>Process</button>
              )}

              {sub === 'processing' && (
                <button className={css.generate} disabled>Processing…</button>
              )}

              {sub === 'processed' && (
                <button className={css.generate} onClick={handleNext} disabled={busy}>
                  Next Step
                </button>
              )}

              {sub === 'advancing' && (
                <button className={css.generate} disabled>Advancing…</button>
              )}
            </>
          )}

          {phase === 'export' && (
            <>
              {exportSub === 'idle' && (
                <button className={css.generate} onClick={handleExport}>Extract Data</button>
              )}
              {exportSub === 'busy' && (
                <button className={css.generate} disabled>Extracting…</button>
              )}
            </>
          )}
        </div>

      </div>
    </div>
  )
}
