import type { GroupDraft, OperationResult, TaskBridge, TaskDraft } from './model'

type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<OperationResult>> } }; runtime?: { EventsOn?(name: string, callback: (event: { kind: string }) => void): () => void } }
const unavailable = (action: string): OperationResult => ({ action, outcome: 'unavailable', message: 'Task authoring is available in the installed desktop application.' })
export function createTaskBridge(nativeWindow: NativeWindow = window): TaskBridge {
  const app = nativeWindow.go?.main?.App
  const call = (name: string, action: string, ...args: unknown[]) => app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))
  return { workspace: () => call('Workspace', 'load'), task: (id) => call('Task', 'load', id), previewTask: (draft) => call('PreviewTask', 'preview', draft), saveTask: (draft) => call('SaveTask', 'save_task', draft), runTask: (id) => call('RunTask', 'run_task', id), setTaskEnabled: (id, enabled) => call('SetTaskEnabled', 'toggle_task', id, enabled), deleteTask: (id) => call('DeleteTask', 'delete_task', id), saveGroup: (draft: GroupDraft) => call('SaveGroup', 'save_group', draft), setGroupEnabled: (id, enabled) => call('SetGroupEnabled', 'toggle_group', id, enabled), deleteGroup: (id) => call('DeleteGroup', 'delete_group', id), subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined) }
}
export const taskBridge = createTaskBridge()
export type { TaskDraft }
