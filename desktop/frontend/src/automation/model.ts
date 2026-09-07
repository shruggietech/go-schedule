export type TaskChoice = { id: string; name: string; readiness: string; reason: string }
export type ChainSummary = { id: string; sourceTaskId: string; sourceTaskName: string; targetTaskId: string; targetTaskName: string; onOutcome: string; readiness: string; reason: string; updatedAt: string }
export type TriggerSummary = { id: string; name: string; targetTaskId: string; targetTaskName: string; enabled: boolean; readiness: string; reason: string; updatedAt: string }
export type TriggerSetMember = { id: string; name: string; position: number; enabled: boolean; readiness: string; reason: string }
export type TriggerSetSummary = { id: string; name: string; targetTaskId: string; targetTaskName: string; memberCount: number; enabledCount: number; members: TriggerSetMember[]; readiness: string; reason: string; updatedAt: string }
export type WatcherSummary = { id: string; name: string; kind: 'file' | 'directory'; path: string; pattern: string; recursive: boolean; debounce: string; stability: string; targetTaskId: string; targetTaskName: string; enabled: boolean; health: string; healthReason: string; readiness: string; reason: string; updatedAt: string }
export type AutomationWorkspace = { tasks: TaskChoice[]; chains: ChainSummary[]; triggers: TriggerSummary[]; triggerSets: TriggerSetSummary[]; watchers: WatcherSummary[]; loadedAt: string }
export type ChainDraft = { id: string; sourceTaskId: string; targetTaskId: string; onOutcome: string; originalUpdatedAt: string; isNew: boolean; overwriteStale: boolean }
export type TriggerDraft = { id: string; name: string; targetTaskId: string; enabled: boolean; originalUpdatedAt: string; isNew: boolean; overwriteStale: boolean }
export type TriggerSetDraft = { id: string; name: string; targetTaskId: string; count: number; enabled: boolean; originalUpdatedAt: string; isNew: boolean; overwriteStale: boolean }
export type WatcherDraft = { id: string; name: string; kind: 'file' | 'directory'; path: string; pattern: string; recursive: boolean; debounce: string; stability: string; targetTaskId: string; enabled: boolean; originalUpdatedAt: string; isNew: boolean; overwriteStale: boolean }
export type OperationResult = { action: string; outcome: 'accepted' | 'rejected' | 'stale' | 'unavailable'; message: string; field?: string; entityId?: string; workspace?: AutomationWorkspace }
export type SecretResult = { action: string; outcome: OperationResult['outcome']; message: string; field?: string; entityId?: string; title?: string; secrets?: Array<{ label: string; key: string; command: string }> }
export interface AutomationBridge {
  workspace(): Promise<OperationResult>
  saveChain(draft: ChainDraft): Promise<OperationResult>; deleteChain(id: string): Promise<OperationResult>
  saveTrigger(draft: TriggerDraft): Promise<SecretResult>; setTriggerEnabled(id: string, enabled: boolean): Promise<OperationResult>; revealTrigger(id: string): Promise<SecretResult>; rotateTrigger(id: string): Promise<SecretResult>; fireTrigger(id: string): Promise<OperationResult>; deleteTrigger(id: string): Promise<OperationResult>
  createTriggerSet(draft: TriggerSetDraft): Promise<SecretResult>; retargetTriggerSet(id: string, target: string, updatedAt: string, overwrite: boolean): Promise<OperationResult>; setTriggerSetEnabled(id: string, enabled: boolean): Promise<OperationResult>; revealTriggerSet(id: string): Promise<SecretResult>; rotateTriggerSet(id: string): Promise<SecretResult>; deleteTriggerSet(id: string): Promise<OperationResult>
  saveWatcher(draft: WatcherDraft): Promise<OperationResult>; setWatcherEnabled(id: string, enabled: boolean): Promise<OperationResult>; deleteWatcher(id: string): Promise<OperationResult>
  subscribe?(listener: (event: { kind: string }) => void): () => void
}
export const blankChain = (): ChainDraft => ({ id: '', sourceTaskId: '', targetTaskId: '', onOutcome: 'success', originalUpdatedAt: '', isNew: true, overwriteStale: false })
export const blankTrigger = (): TriggerDraft => ({ id: '', name: '', targetTaskId: '', enabled: true, originalUpdatedAt: '', isNew: true, overwriteStale: false })
export const blankSet = (): TriggerSetDraft => ({ id: '', name: '', targetTaskId: '', count: 2, enabled: true, originalUpdatedAt: '', isNew: true, overwriteStale: false })
export const blankWatcher = (): WatcherDraft => ({ id: '', name: '', kind: 'file', path: '', pattern: '', recursive: false, debounce: '500ms', stability: '1s', targetTaskId: '', enabled: true, originalUpdatedAt: '', isNew: true, overwriteStale: false })
