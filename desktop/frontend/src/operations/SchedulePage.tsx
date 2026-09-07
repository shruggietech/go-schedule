import { useMemo, useState } from 'react'
import { Button, Notice, StatePanel } from '../components'
import type { OperationsBridge, ScheduleOccurrence } from './model'
import { useSchedule } from './store'

const stateLabel = (value: string) => ({ upcoming: 'Upcoming', running: 'Running', success: 'Success', failure: 'Failed', skipped: 'Skipped', caught_up: 'Caught up', queued: 'Queued', unavailable: 'Unavailable' }[value] ?? 'Unavailable')
const kindLabel = (value: string) => value === 'prediction' ? 'Prediction' : 'Recorded run'
const localTime = (value: string) => new Date(value).toLocaleString()
const dayKey = (value: string | Date) => { const date = typeof value === 'string' ? new Date(value) : value; return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}` }
const dayLabel = (value: string | Date) => (typeof value === 'string' ? new Date(`${value}T12:00:00`) : value).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })

function OccurrenceDetail({ item, missing }: { item?: ScheduleOccurrence; missing: boolean }) {
  if (missing) return <aside className="panel operation-detail"><h2>Occurrence no longer in range</h2><p>The selected occurrence was not present in the latest complete snapshot. Your range and view were preserved.</p></aside>
  if (!item) return <aside className="panel operation-detail"><h2>Occurrence details</h2><p>Select an occurrence to inspect its operational identity.</p></aside>
  return <aside className="panel operation-detail"><h2>{item.taskName}</h2><dl><dt>Record</dt><dd>{kindLabel(item.kind)}</dd><dt>State</dt><dd><span className={`operation-state state-${item.state}`}>{stateLabel(item.state)}</span></dd><dt>Time</dt><dd>{localTime(item.time)}</dd><dt>Task ID</dt><dd><code>{item.taskId || 'Unavailable'}</code></dd><dt>Run ID</dt><dd><code>{item.runId || 'Not recorded'}</code></dd></dl></aside>
}

export function SchedulePage({ bridge, available, refreshToken }: { bridge: OperationsBridge; available: boolean; refreshToken: number }) {
  const [days, setDays] = useState(7)
  const [view, setView] = useState<'agenda' | 'calendar'>('agenda')
  const [selectedID, setSelectedID] = useState('')
  const [selectedDay, setSelectedDay] = useState('')
  const [month, setMonth] = useState<Date>()
  const { snapshot, status, load } = useSchedule(bridge, days, available, refreshToken)
  const selected = snapshot?.occurrences.find((item) => item.id === selectedID)
  const missing = Boolean(selectedID && snapshot && !selected)
  const daysWithItems = useMemo(() => {
    const grouped = new Map<string, ScheduleOccurrence[]>()
    for (const item of snapshot?.occurrences ?? []) { const day = dayKey(item.time); grouped.set(day, [...(grouped.get(day) ?? []), item]) }
    return grouped
  }, [snapshot])
  const calendarMonth = useMemo(() => { const value = month ?? new Date(snapshot?.occurrences[0]?.time ?? snapshot?.from ?? Date.now()); return new Date(value.getFullYear(), value.getMonth(), 1) }, [month, snapshot])
  const calendarDays = useMemo(() => { const lead = (calendarMonth.getDay() + 6) % 7; const count = new Date(calendarMonth.getFullYear(), calendarMonth.getMonth() + 1, 0).getDate(); return [...Array.from({ length: lead }, () => null), ...Array.from({ length: count }, (_, index) => new Date(calendarMonth.getFullYear(), calendarMonth.getMonth(), index + 1))] }, [calendarMonth])
  return <>
    <header className="page-header"><div><p className="eyebrow">This computer</p><h1>Schedule</h1><p>See predicted work and recorded runs without confusing one for the other.</p></div><Button variant="secondary" onClick={() => void load()} disabled={!available}>Refresh</Button></header>
    {!available && <Notice title="Read-only schedule" tone="warning">The last complete schedule remains visible while the scheduler reconnects.</Notice>}
    {status?.outcome !== 'accepted' && status && <Notice title="Schedule could not refresh" tone="error">{status.message}</Notice>}
    <section className="panel operation-filters" aria-label="Schedule controls"><label>View<select value={view} onChange={(event) => setView(event.target.value as 'agenda' | 'calendar')}><option value="agenda">Agenda</option><option value="calendar">Calendar</option></select></label><label>Window<select value={days} onChange={(event) => setDays(Number(event.target.value))}><option value={1}>1 day</option><option value={7}>7 days</option><option value={30}>30 days</option></select></label>{snapshot && <p role="status">{snapshot.occurrences.length} occurrences from {localTime(snapshot.from)} through {localTime(snapshot.to)}</p>}</section>
    {!snapshot ? <StatePanel title="Loading schedule" detail="Requesting a complete schedule snapshot." busy={available} /> : snapshot.occurrences.length === 0 ? <StatePanel title="No occurrences in this window" detail="Choose another range or create an enabled scheduled task." /> : <div className="operation-layout">
      <section className="panel operation-list">{view === 'agenda' ? <div className="table-scroll" tabIndex={0}><table><caption>Predictions and recorded runs</caption><thead><tr><th scope="col">When</th><th scope="col">Task</th><th scope="col">Record</th><th scope="col">State</th></tr></thead><tbody>{snapshot.occurrences.map((item) => <tr key={item.id} aria-selected={item.id === selectedID}><td><button className="row-select" onClick={() => setSelectedID(item.id)}>{localTime(item.time)}</button></td><td>{item.taskName}</td><td>{kindLabel(item.kind)}</td><td><span className={`operation-state state-${item.state}`}>{stateLabel(item.state)}</span></td></tr>)}</tbody></table></div> : <><div className="calendar-toolbar"><Button variant="secondary" onClick={() => setMonth(new Date(calendarMonth.getFullYear(), calendarMonth.getMonth() - 1, 1))}>Previous month</Button><h2>{calendarMonth.toLocaleDateString(undefined, { month: 'long', year: 'numeric' })}</h2><Button variant="secondary" onClick={() => setMonth(new Date(calendarMonth.getFullYear(), calendarMonth.getMonth() + 1, 1))}>Next month</Button></div><div className="calendar-grid" aria-label="Schedule calendar">{['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'].map((day) => <strong className="calendar-weekday" key={day}>{day}</strong>)}{calendarDays.map((day, index) => day ? <button key={dayKey(day)} aria-label={`${dayLabel(day)}, ${(daysWithItems.get(dayKey(day)) ?? []).length} ${(daysWithItems.get(dayKey(day)) ?? []).length === 1 ? 'occurrence' : 'occurrences'}`} aria-pressed={selectedDay === dayKey(day)} onClick={() => setSelectedDay(dayKey(day))}><strong>{day.getDate()}</strong><span>{(daysWithItems.get(dayKey(day)) ?? []).length} {(daysWithItems.get(dayKey(day)) ?? []).length === 1 ? 'occurrence' : 'occurrences'}</span></button> : <span aria-hidden="true" key={`blank-${index}`} />)}</div>{selectedDay && <div className="calendar-day"><h2>{dayLabel(selectedDay)}</h2>{(daysWithItems.get(selectedDay) ?? []).length ? (daysWithItems.get(selectedDay) ?? []).map((item) => <button className="row-select" key={item.id} onClick={() => setSelectedID(item.id)}>{localTime(item.time)} · {item.taskName} · {stateLabel(item.state)}</button>) : <p>No occurrences on this day.</p>}</div>}</>}</section>
      <OccurrenceDetail item={selected} missing={missing} />
    </div>}
  </>
}
