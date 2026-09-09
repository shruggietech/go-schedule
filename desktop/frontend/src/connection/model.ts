export type ConnectionState = 'connecting' | 'connected' | 'degraded' | 'recovering' | 'unavailable' | 'access_denied' | 'incompatible' | 'timed_out'
export type Appearance = 'system' | 'light' | 'dark'
export type Route = 'tasks' | 'automation' | 'schedule' | 'activity' | 'notifications' | 'agentAccess' | 'connections' | 'settings'

export interface Target {
  id: string
  displayName: string
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

export interface ActionResult { action: string; outcome: 'accepted' | 'rejected' | 'unavailable'; message: string }

export interface DesktopBridge {
  snapshot(): Promise<ConnectionSnapshot>
  retry(): Promise<ActionResult>
  quit(): Promise<ActionResult>
  subscribe(listener: (event: DesktopEvent) => void): () => void
}
