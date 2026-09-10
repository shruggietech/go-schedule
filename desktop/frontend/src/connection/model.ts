export type ConnectionState = 'connecting' | 'connected' | 'degraded' | 'recovering' | 'unavailable' | 'access_denied' | 'unauthorized' | 'revoked' | 'forbidden' | 'incompatible' | 'trust_changed' | 'identity_changed' | 'timed_out'
export type Appearance = 'system' | 'light' | 'dark'
export type Route = 'tasks' | 'automation' | 'schedule' | 'activity' | 'notifications' | 'agentAccess' | 'connections' | 'settings'

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

export interface DesktopBridge {
  snapshot(): Promise<ConnectionSnapshot>
  retry(): Promise<ActionResult>
  connectionProfiles?(): Promise<ConnectionResult>
  selectConnection?(id: string): Promise<ConnectionResult>
  renameConnection?(id: string, label: string): Promise<ConnectionResult>
  removeConnection?(id: string): Promise<ConnectionResult>
  quit(): Promise<ActionResult>
  subscribe(listener: (event: DesktopEvent) => void): () => void
}
