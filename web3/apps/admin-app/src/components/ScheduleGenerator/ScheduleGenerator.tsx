import { useState } from 'react'
import {
  loadDataIntoGenerator,
  loadClassroomsIntoGenerator,
  submitAndGo,
  processStep,
  extractWorkloadsFromGenerator,
  extractDataFromGenerator,
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
type StepSub = 'idle' | 'processing' | 'processed' | 'advancing' | 'warning'
type ExportSub = 'idle' | 'busy' | 'workloads-done' | 'done'

interface StepConfig { enabled: boolean; method: string }
interface LogLine { text: string; kind: 'run' | 'ok' | 'warn' | 'error' | 'info' }

function initConfigs(): Record<string, StepConfig> {
  return Object.fromEntries(PIPELINE.map(s => [s.id, { enabled: true, method: s.methods[0] ?? '' }]))
}

function parseIds(raw: string): string[] {
  return raw.split(/[\s,]+/).map(x => x.trim()).filter(Boolean)
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

// ── Main component ────────────────────────────────────────────────────────────

export function ScheduleGenerator({ onClose }: { onClose: () => void }) {
  const [semesters,   setSemesters]   = useState('')
  const [classrooms,  setClassrooms]  = useState('')
  const [startTime,   setStartTime]   = useState('')
  const [configs,     setConfigs]     = useState<Record<string, StepConfig>>(initConfigs)

  const [phase,       setPhase]       = useState<Phase>('setup')
  const [stepIdx,     setStepIdx]     = useState(0)
  const [sub,         setSub]         = useState<StepSub>('idle')
  const [stepResult,  setStepResult]  = useState<unknown>(null)
  const [warnings,    setWarnings]    = useState<string[]>([])
  const [exportSub,   setExportSub]   = useState<ExportSub>('idle')
  const [setupError,  setSetupError]  = useState<string | null>(null)
  const [log,         setLog]         = useState<LogLine[]>([])

  const addLog = (text: string, kind: LogLine['kind']) =>
    setLog(prev => [...prev, { text, kind }])

  const patchConfig = (id: string, patch: Partial<StepConfig>) =>
    setConfigs(prev => ({ ...prev, [id]: { ...prev[id], ...patch } }))

  const busy = sub === 'processing' || sub === 'advancing' || exportSub === 'busy'

  // ── Setup → Pipeline ───────────────────────────────────────────────────────

  async function handleStart() {
    const semIds = parseIds(semesters)
    if (semIds.length === 0) { setSetupError('Provide at least one curriculum semester ID'); return }
    if (!startTime)          { setSetupError('Provide a schedule start time'); return }
    setSetupError(null)
    setSub('advancing')

    try {
      addLog(`Loading data (${semIds.length} semester(s))…`, 'run')
      await loadDataIntoGenerator(semIds)
      addLog('Data loaded ✓', 'ok')

      const clsIds = parseIds(classrooms)
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

  async function handleNext(force = false) {
    setSub('advancing')

    try {
      addLog(`Submit & go${force ? ' (forced)' : ''}…`, 'run')
      const res = await submitAndGo(force)
      const warns = res.data?.warnings ?? []

      if (warns.length > 0 && !force) {
        setWarnings(warns)
        addLog(`Warnings (${warns.length}) — review before continuing`, 'warn')
        setSub('warning')
        return
      }

      setWarnings([])
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

  async function handleExportWorkloads() {
    setExportSub('busy')
    try {
      addLog('Exporting workloads…', 'run')
      await extractWorkloadsFromGenerator()
      addLog('Workloads exported ✓', 'ok')
      setExportSub('workloads-done')
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error: ${msg}`, 'error')
      setExportSub('idle')
    }
  }

  async function handleExportLessons() {
    setExportSub('busy')
    try {
      addLog('Exporting lessons…', 'run')
      const iso = new Date(startTime).toISOString()
      await extractDataFromGenerator(iso)
      addLog('Lessons exported ✓', 'ok')
      addLog('Schedule generated successfully ✓', 'ok')
      setExportSub('done')
      setPhase('done')
    } catch (e) {
      const msg = e instanceof Error ? e.message : String(e)
      addLog(`Error: ${msg}`, 'error')
      setExportSub('workloads-done')
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
              <label className={css.field}>
                <span className={css.label}>Curriculum semester IDs</span>
                <textarea className={css.input} rows={2} placeholder="uuid, uuid, …"
                  value={semesters} onChange={e => setSemesters(e.target.value)} />
              </label>

              <label className={css.field}>
                <span className={css.label}>Classroom IDs (optional)</span>
                <textarea className={css.input} rows={2} placeholder="uuid, uuid, …"
                  value={classrooms} onChange={e => setClassrooms(e.target.value)} />
              </label>

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

              {sub === 'warning' && warnings.length > 0 && (
                <div className={css.warningBox}>
                  <p className={css.warningTitle}>Warnings from submit</p>
                  <ul className={css.warningList}>
                    {warnings.map((w, i) => <li key={i}>{w}</li>)}
                  </ul>
                </div>
              )}
            </div>
          )}

          {/* ── Export ── */}
          {phase === 'export' && (
            <div className={css.exportPhase}>
              <p className={css.exportTitle}>Export Results</p>
              <p className={css.exportHint}>Export workloads first, then the lesson schedule.</p>
              <div className={css.exportSteps}>
                <div className={`${css.exportStep} ${exportSub !== 'idle' ? css.exportStepDone : ''}`}>
                  <span className={css.exportStepNum}>1</span>
                  <span>Workloads</span>
                  {(exportSub === 'workloads-done' || exportSub === 'done') && (
                    <span className={css.exportCheck}>✓</span>
                  )}
                </div>
                <div className={`${css.exportStep} ${exportSub === 'done' ? css.exportStepDone : ''}`}>
                  <span className={css.exportStepNum}>2</span>
                  <span>Lessons</span>
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
                <button className={css.generate} onClick={() => handleNext(false)} disabled={busy}>
                  Next Step
                </button>
              )}

              {sub === 'advancing' && (
                <button className={css.generate} disabled>Advancing…</button>
              )}

              {sub === 'warning' && (
                <>
                  <button className={css.secondary} onClick={() => setSub('processed')} disabled={busy}>
                    Back
                  </button>
                  <button className={css.generate} onClick={() => handleNext(true)} disabled={busy}>
                    Force Continue
                  </button>
                </>
              )}
            </>
          )}

          {phase === 'export' && (
            <>
              {exportSub === 'idle' && (
                <button className={css.generate} onClick={handleExportWorkloads}>Export Workloads</button>
              )}
              {exportSub === 'busy' && (
                <button className={css.generate} disabled>Exporting…</button>
              )}
              {exportSub === 'workloads-done' && (
                <button className={css.generate} onClick={handleExportLessons}>Export Lessons</button>
              )}
            </>
          )}
        </div>

      </div>
    </div>
  )
}
