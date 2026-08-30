export function ScheduleRow({ name, expression, nextRun, state = "ready" }) {
  return <div className="gs-schedule-row"><strong>{name}</strong><code className="gs-schedule-row__expression">{expression}</code><span className="gs-schedule-row__next">{nextRun || state}</span></div>;
}
