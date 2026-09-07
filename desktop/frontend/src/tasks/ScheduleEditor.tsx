import { Field } from '../components'
import type { TaskDraft } from './model'

function wallTimeValue(value: string, timezone: string) {
  if (!/(?:Z|[+-]\d\d:\d\d)$/.test(value)) return value.slice(0, 16)
  try {
    const parts = new Intl.DateTimeFormat('en-CA', { timeZone: timezone && timezone !== 'Local' ? timezone : undefined, year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hourCycle: 'h23' }).formatToParts(new Date(value))
    const part = (type: Intl.DateTimeFormatPartTypes) => parts.find((candidate) => candidate.type === type)?.value ?? ''
    return `${part('year')}-${part('month')}-${part('day')}T${part('hour')}:${part('minute')}`
  } catch {
    return value.slice(0, 16)
  }
}

export function ScheduleEditor({ draft, update, errorField, errorMessage }: { draft: TaskDraft; update(values: Partial<TaskDraft>): void; errorField?: string; errorMessage?: string }) {
  const error = (field: string) => errorField === field ? errorMessage : undefined
  return <fieldset><legend>Timing</legend><Field label="Timing mode" error={error('mode')}><select name="mode" value={draft.mode} onChange={(event) => update({ mode: event.target.value as TaskDraft['mode'] })}><option value="recurring">Recurring</option><option value="one_off">One time</option><option value="manual">Manual only</option></select></Field>{draft.mode === 'recurring' && <><Field label="Schedule" help="Use a supported phrase or cron expression." error={error('schedule')}><input name="schedule" value={draft.schedule} onChange={(event) => update({ schedule: event.target.value })} /></Field><Field label="Schedule syntax" error={error('schedule_syntax')}><select name="schedule_syntax" value={draft.scheduleSyntax} onChange={(event) => update({ scheduleSyntax: event.target.value })}><option value="">Detect automatically</option><option value="human">Human phrase</option><option value="cron">Cron</option></select></Field></>}{draft.mode === 'one_off' && <Field label="Run at" help={`Interpreted in ${draft.timezone || 'the local timezone'}.`} error={error('at')}><input name="at" type="datetime-local" value={wallTimeValue(draft.at, draft.timezone)} onChange={(event) => update({ at: event.target.value })} /></Field>}<Field label="Timezone" error={error('timezone')}><input name="timezone" value={draft.timezone} onChange={(event) => update({ timezone: event.target.value })} /></Field>{draft.nextRuns.length > 0 && <ol aria-label="Next five runs">{draft.nextRuns.map((run) => <li key={run}>{new Date(run).toLocaleString()}</li>)}</ol>}</fieldset>
}
