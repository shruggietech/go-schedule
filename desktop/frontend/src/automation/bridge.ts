import type { AutomationBridge, ChainDraft, OperationResult, SecretResult, TriggerDraft, TriggerSetDraft, WatcherDraft } from './model'

type Result = OperationResult | SecretResult
type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<Result>> } }; runtime?: { EventsOn?(name: string, callback: (event: { kind: string }) => void): () => void } }
const unavailable = (action: string): OperationResult => ({ action, outcome: 'unavailable', message: 'Automation sources are available in the installed desktop application.' })
export function createAutomationBridge(nativeWindow: NativeWindow = window): AutomationBridge {
  const app = nativeWindow.go?.main?.App
  const operation = (name: string, action: string, ...args: unknown[]) => (app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))) as Promise<OperationResult>
  const secret = (name: string, action: string, ...args: unknown[]) => (app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))) as Promise<SecretResult>
  return {
    workspace: (daemonID) => operation('AutomationWorkspace', 'load', daemonID),
    saveJoin: (draft) => operation('SaveJoin', 'save_join', draft), deleteJoin: (id, daemonID) => operation('DeleteJoin', 'delete_join', id, daemonID), joinInstances: (id, offset, daemonID) => operation('JoinInstances', 'load_join_instances', id, offset, daemonID), joinInstance: (id, daemonID) => operation('JoinInstance', 'load_join_instance', id, daemonID),
    saveChain: (draft: ChainDraft) => operation('SaveChain', 'save_chain', draft), deleteChain: (id, daemonID) => operation('DeleteChain', 'delete_chain', id, daemonID),
    saveTrigger: (draft: TriggerDraft) => secret('SaveTrigger', 'save_trigger', draft), setTriggerEnabled: (id, enabled, daemonID) => operation('SetTriggerEnabled', 'toggle_trigger', id, enabled, daemonID), revealTrigger: (id, daemonID) => secret('RevealTrigger', 'reveal_trigger', id, daemonID), rotateTrigger: (id, daemonID) => secret('RotateTrigger', 'rotate_trigger', id, daemonID), fireTrigger: (id, daemonID) => operation('FireTrigger', 'fire_trigger', id, daemonID), deleteTrigger: (id, daemonID) => operation('DeleteTrigger', 'delete_trigger', id, daemonID),
    createTriggerSet: (draft: TriggerSetDraft) => secret('CreateTriggerSet', 'create_trigger_set', draft), retargetTriggerSet: (id, target, updatedAt, overwrite, daemonID) => operation('RetargetTriggerSet', 'retarget_trigger_set', id, target, updatedAt, overwrite, daemonID), setTriggerSetEnabled: (id, enabled, daemonID) => operation('SetTriggerSetEnabled', 'toggle_trigger_set', id, enabled, daemonID), revealTriggerSet: (id, daemonID) => secret('RevealTriggerSet', 'reveal_trigger_set', id, daemonID), rotateTriggerSet: (id, daemonID) => secret('RotateTriggerSet', 'rotate_trigger_set', id, daemonID), deleteTriggerSet: (id, daemonID) => operation('DeleteTriggerSet', 'delete_trigger_set', id, daemonID),
    saveWatcher: (draft: WatcherDraft) => operation('SaveFilesystemWatcher', 'save_watcher', draft), setWatcherEnabled: (id, enabled, daemonID) => operation('SetFilesystemWatcherEnabled', 'toggle_watcher', id, enabled, daemonID), deleteWatcher: (id, daemonID) => operation('DeleteFilesystemWatcher', 'delete_watcher', id, daemonID),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined),
  }
}
export const automationBridge = createAutomationBridge()
