import type { AutomationBridge, ChainDraft, OperationResult, SecretResult, TriggerDraft, TriggerSetDraft, WatcherDraft } from './model'

type Result = OperationResult | SecretResult
type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<Result>> } }; runtime?: { EventsOn?(name: string, callback: (event: { kind: string }) => void): () => void } }
const unavailable = (action: string): OperationResult => ({ action, outcome: 'unavailable', message: 'Automation sources are available in the installed desktop application.' })
export function createAutomationBridge(nativeWindow: NativeWindow = window): AutomationBridge {
  const app = nativeWindow.go?.main?.App
  const operation = (name: string, action: string, ...args: unknown[]) => (app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))) as Promise<OperationResult>
  const secret = (name: string, action: string, ...args: unknown[]) => (app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))) as Promise<SecretResult>
  return {
    workspace: () => operation('AutomationWorkspace', 'load'),
    saveChain: (draft: ChainDraft) => operation('SaveChain', 'save_chain', draft), deleteChain: (id) => operation('DeleteChain', 'delete_chain', id),
    saveTrigger: (draft: TriggerDraft) => secret('SaveTrigger', 'save_trigger', draft), setTriggerEnabled: (id, enabled) => operation('SetTriggerEnabled', 'toggle_trigger', id, enabled), revealTrigger: (id) => secret('RevealTrigger', 'reveal_trigger', id), rotateTrigger: (id) => secret('RotateTrigger', 'rotate_trigger', id), fireTrigger: (id) => operation('FireTrigger', 'fire_trigger', id), deleteTrigger: (id) => operation('DeleteTrigger', 'delete_trigger', id),
    createTriggerSet: (draft: TriggerSetDraft) => secret('CreateTriggerSet', 'create_trigger_set', draft), retargetTriggerSet: (id, target, updatedAt, overwrite) => operation('RetargetTriggerSet', 'retarget_trigger_set', id, target, updatedAt, overwrite), setTriggerSetEnabled: (id, enabled) => operation('SetTriggerSetEnabled', 'toggle_trigger_set', id, enabled), revealTriggerSet: (id) => secret('RevealTriggerSet', 'reveal_trigger_set', id), rotateTriggerSet: (id) => secret('RotateTriggerSet', 'rotate_trigger_set', id), deleteTriggerSet: (id) => operation('DeleteTriggerSet', 'delete_trigger_set', id),
    saveWatcher: (draft: WatcherDraft) => operation('SaveFilesystemWatcher', 'save_watcher', draft), setWatcherEnabled: (id, enabled) => operation('SetFilesystemWatcherEnabled', 'toggle_watcher', id, enabled), deleteWatcher: (id) => operation('DeleteFilesystemWatcher', 'delete_watcher', id),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined),
  }
}
export const automationBridge = createAutomationBridge()
