import type { ConnectionState, SystemFailure, SystemRegistration } from '../connection/model'

export type SearchKind = 'task' | 'group' | 'failure' | 'schedule' | 'alert'
export type SearchAction = 'open' | 'acknowledge' | 'enable' | 'disable' | 'run_now'
export interface SearchResult { kind: SearchKind; object_id: string; task_id?: string; name: string; context?: string; occurred_at?: string; enabled?: boolean; action_hints: SearchAction[] }
export interface SearchMatch { registrationKey: string; expectedDaemonId: string; sourceLabel: string; sourceShortId: string; result: SearchResult; availableActions: SearchAction[]; disabledReason?: string }
export interface SearchObservation { registration: SystemRegistration; state: ConnectionState; observedAt?: string; truncated: boolean; matches: SearchMatch[]; failure?: SystemFailure }
export interface SearchSnapshot { generation: number; query: string; startedAt: string; completedAt?: string; complete: boolean; observations: SearchObservation[] }
export interface SearchRequest { query: string; kinds?: SearchKind[]; limit?: number }
export interface SearchSelection { registrationKey: string; expectedDaemonId: string; kind: SearchKind; objectId: string; displayName: string }
export interface SearchIntent { action: SearchAction; selections: SearchSelection[] }
export interface SearchActionOutcome extends SearchSelection { currentDaemonId?: string; outcome: string; message: string }
export interface SearchActionResult { action: SearchAction; outcome: 'accepted' | 'partial' | 'rejected' | 'unavailable'; message: string; outcomes: SearchActionOutcome[] }
export interface SearchBridge { search(request: SearchRequest): Promise<SearchSnapshot>; execute(intent: SearchIntent): Promise<SearchActionResult>; subscribe(listener: (snapshot: SearchSnapshot) => void): () => void }
