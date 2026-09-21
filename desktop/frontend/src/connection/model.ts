export type ConnectionState = 'connecting' | 'connected' | 'degraded' | 'recovering' | 'unavailable' | 'access_denied' | 'unauthorized' | 'revoked' | 'forbidden' | 'incompatible' | 'trust_changed' | 'identity_changed' | 'timed_out'
export type Appearance = 'system' | 'light' | 'dark'
export type Route = 'systems' | 'tasks' | 'automation' | 'schedule' | 'activity' | 'notifications' | 'agentAccess' | 'connections' | 'settings'

export interface Target {
  id: string
  profileId?: string
  kind?: 'local' | 'remote'
  displayName: string
  endpoint?: string
  fingerprint?: string
  platform: string
  architecture?: string
  version?: string
  capabilities: string[]
  permissions: string[]
}

export interface ConnectionSnapshot {
  generation: number
  revision: number
  state: ConnectionState
  target: Target
  message: string
  action?: string
  lastSuccessfulAt?: string
  stale?: boolean
  retryAttempt?: number
  nextRetryAt?: string
  recovery?: 'none' | 'automatic' | 'manual'
}

export interface DesktopEvent {
  id: string
  kind: string
  message: string
  generation: number
  occurredAt: string
  entityId?: string
  snapshot?: ConnectionSnapshot
}

export interface ActionResult { action: string; outcome: 'accepted' | 'rejected' | 'unavailable' | 'uncertain'; message: string }
export interface ConnectionProfile { id: string; label: string; endpoint: string; daemonId: string; shortDaemonId: string; fingerprint: string; capability: string; platform: string; architecture?: string; productVersion?: string; lastSuccessfulAt?: string; active: boolean }
export interface ConnectionWorkspace { activeProfileId?: string; profiles: ConnectionProfile[] }
export interface ConnectionResult extends ActionResult { workspace?: ConnectionWorkspace }

export interface UpcomingSummary { task_id: string; task_name: string; scheduled_for: string }
export interface FailureSummary { run_id: string; task_id: string; task_name: string; ended_at: string }
export interface AlertSummary { alert_id: string; task_id?: string; run_id?: string; severity: string; kind: string; created_at: string }
export interface NotificationProblemSummary { delivery_id: string; task_id?: string; run_id?: string; channel_name: string; state: string; created_at: string }
export interface SystemSummary { schema: string; observed_at: string; active_task_count: number; next_occurrence?: UpcomingSummary; recent_failure_count: number; recent_failure?: FailureSummary; unacknowledged_alert_count: number; unacknowledged_alert?: AlertSummary; notification_problem_count: number; notification_problem?: NotificationProblemSummary }
export interface SystemRegistration { key: string; profileId?: string; kind: 'local' | 'remote'; label: string; endpoint?: string; daemonId?: string; shortDaemonId?: string; platform?: string; architecture?: string; version?: string }
export interface SystemFailure { state: ConnectionState; message: string; action: string }
export interface SystemObservation { registration: SystemRegistration; state: ConnectionState; observedAt?: string; stale: boolean; summary?: SystemSummary; failure?: SystemFailure }
export interface SystemsSnapshot { generation: number; startedAt: string; completedAt: string; complete: boolean; observations: SystemObservation[] }

export interface DesktopBridge {
  snapshot(): Promise<ConnectionSnapshot>
  retry(): Promise<ActionResult>
  connectionProfiles?(): Promise<ConnectionResult>
  selectConnection?(id: string): Promise<ConnectionResult>
  renameConnection?(id: string, label: string): Promise<ConnectionResult>
  removeConnection?(id: string): Promise<ConnectionResult>
  allSystems?(): Promise<SystemsSnapshot>
  subscribeSystems?(listener: (snapshot: SystemsSnapshot) => void): () => void
  quit(): Promise<ActionResult>
  subscribe(listener: (event: DesktopEvent) => void): () => void
}
