export type Condition = 'new' | 'empty' | 'loading' | 'connected' | 'disconnected' | 'degraded' | 'destructive' | 'validation-error' | 'success'
export type Page = 'tasks' | 'editor' | 'schedule' | 'activity' | 'targets'
export type Appearance = 'system' | 'light' | 'dark'

export interface TargetIdentity {
  id: string
  displayName: string
  platform: string
  connection: 'connected' | 'disconnected' | 'degraded' | 'loading'
  detail: string
}

export interface TaskSummary {
  id: string
  name: string
  enabled: boolean
  state: string
  schedule: string
  nextRun: string
  lastResult: string
}

export interface ProofSnapshot {
  target: TargetIdentity
  health: { status: TargetIdentity['connection']; version?: string; message: string }
  tasks: TaskSummary[]
  generatedAt: string
}

export interface ProofEvent {
  id: string
  kind: string
  message: string
  occurredAt: string
  taskId?: string
}

type WailsWindow = Window & {
  go?: { main?: { ProofApp?: { Snapshot(): Promise<ProofSnapshot>; ShowAbout(): Promise<{ outcome: string; message: string }> } } }
  runtime?: { EventsOn?(name: string, callback: (event: ProofEvent) => void): () => void }
}

export const connectedFixture: ProofSnapshot = {
  target: { id: 'local', displayName: 'This computer', platform: 'windows', connection: 'connected', detail: 'Local scheduler service' },
  health: { status: 'connected', version: 'v1.1.1', message: 'Scheduler service is available.' },
  tasks: [
    { id: 'backup', name: 'Nightly backup', enabled: true, state: 'active', schedule: 'Every day at 23:00', nextRun: 'Today, 23:00', lastResult: 'success' },
    { id: 'reports', name: 'Publish weekly reports', enabled: true, state: 'active', schedule: 'Monday at 08:30', nextRun: 'Mon, 08:30', lastResult: 'failed' },
    { id: 'cleanup', name: 'Clear temporary files', enabled: false, state: 'paused', schedule: 'Every Friday at 18:00', nextRun: '', lastResult: 'success' },
  ],
  generatedAt: '2026-09-07T05:30:00Z',
}

export function fixtureFor(condition: Condition): ProofSnapshot {
  const snapshot = structuredClone(connectedFixture)
  snapshot.target.connection = condition === 'loading' ? 'loading' : condition === 'disconnected' ? 'disconnected' : condition === 'degraded' ? 'degraded' : 'connected'
  snapshot.health.status = snapshot.target.connection
  if (condition === 'empty' || condition === 'new') snapshot.tasks = []
  if (condition === 'disconnected') snapshot.health.message = 'The scheduler service is unavailable. Start the service, then try again.'
  if (condition === 'degraded') snapshot.health.message = 'Live updates are delayed. Displayed data may be out of date.'
  if (condition === 'loading') snapshot.health.message = 'Connecting to the scheduler service.'
  return snapshot
}

export async function loadSnapshot(): Promise<ProofSnapshot> {
  const bridge = (window as WailsWindow).go?.main?.ProofApp
  return bridge ? bridge.Snapshot() : connectedFixture
}

export async function showAbout(): Promise<void> {
  const bridge = (window as WailsWindow).go?.main?.ProofApp
  if (bridge) await bridge.ShowAbout()
}

export function subscribeToProofEvents(callback: (event: ProofEvent) => void): () => void {
  return (window as WailsWindow).runtime?.EventsOn?.('proof:event', callback) ?? (() => undefined)
}
