import type { ActionResult, ConnectionResult, ConnectionSnapshot, DesktopBridge, DesktopEvent } from './model'

type NativeWindow = Window & {
  go?: { main?: { App?: { Snapshot(): Promise<ConnectionSnapshot>; RetryConnection(): Promise<ActionResult>; Quit(): Promise<ActionResult>; ConnectionProfiles?(): Promise<ConnectionResult>; SelectConnection?(id: string): Promise<ConnectionResult>; RenameConnection?(id: string, label: string): Promise<ConnectionResult>; RemoveConnection?(id: string): Promise<ConnectionResult> } } }
  runtime?: { EventsOn?(name: string, callback: (event: DesktopEvent) => void): () => void }
}

export const unavailableSnapshot: ConnectionSnapshot = {
  generation: 0,
  revision: 0,
  state: 'unavailable',
  target: { id: 'local', displayName: 'This computer', platform: 'unknown', capabilities: [], permissions: [] },
  message: 'The native desktop connection is unavailable in this browser.',
  action: 'Open the installed desktop application.',
  stale: false,
  recovery: 'manual',
}

export function createBridge(nativeWindow: NativeWindow = window): DesktopBridge {
  const app = nativeWindow.go?.main?.App
  return {
    snapshot: () => app?.Snapshot() ?? Promise.resolve(unavailableSnapshot),
    retry: () => app?.RetryConnection() ?? Promise.resolve({ action: 'retry', outcome: 'unavailable', message: 'Retry is available in the installed desktop application.' }),
    connectionProfiles: () => app?.ConnectionProfiles?.() ?? Promise.resolve({ action: 'load_connections', outcome: 'unavailable', message: 'Connection profiles are available in the installed desktop application.' }),
    selectConnection: (id) => app?.SelectConnection?.(id) ?? Promise.resolve({ action: 'select_connection', outcome: 'unavailable', message: 'Connection selection is unavailable.' }),
    renameConnection: (id, label) => app?.RenameConnection?.(id, label) ?? Promise.resolve({ action: 'rename_connection', outcome: 'unavailable', message: 'Connection profiles are unavailable.' }),
    removeConnection: (id) => app?.RemoveConnection?.(id) ?? Promise.resolve({ action: 'remove_connection', outcome: 'unavailable', message: 'Connection profiles are unavailable.' }),
    quit: () => app?.Quit() ?? Promise.resolve({ action: 'quit', outcome: 'unavailable', message: 'Exit is available in the installed desktop application.' }),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined),
  }
}

export const desktopBridge = createBridge()
