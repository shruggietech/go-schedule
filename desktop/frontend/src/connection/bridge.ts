import type { ActionResult, ConnectionSnapshot, DesktopBridge, DesktopEvent } from './model'

type NativeWindow = Window & {
  go?: { main?: { App?: { Snapshot(): Promise<ConnectionSnapshot>; RetryConnection(): Promise<ActionResult>; Quit(): Promise<ActionResult> } } }
  runtime?: { EventsOn?(name: string, callback: (event: DesktopEvent) => void): () => void }
}

export const unavailableSnapshot: ConnectionSnapshot = {
  generation: 0,
  state: 'unavailable',
  target: { id: 'local', displayName: 'This computer', platform: 'unknown', capabilities: [], permissions: [] },
  message: 'The native desktop connection is unavailable in this browser.',
  action: 'Open the installed desktop application.',
}

export function createBridge(nativeWindow: NativeWindow = window): DesktopBridge {
  const app = nativeWindow.go?.main?.App
  return {
    snapshot: () => app?.Snapshot() ?? Promise.resolve(unavailableSnapshot),
    retry: () => app?.RetryConnection() ?? Promise.resolve({ action: 'retry', outcome: 'unavailable', message: 'Retry is available in the installed desktop application.' }),
    quit: () => app?.Quit() ?? Promise.resolve({ action: 'quit', outcome: 'unavailable', message: 'Exit is available in the installed desktop application.' }),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined),
  }
}

export const desktopBridge = createBridge()
