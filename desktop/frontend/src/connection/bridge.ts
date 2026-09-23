import type { ActionResult, ConnectionResult, ConnectionSnapshot, DesktopBridge, DesktopEvent, LocalServiceActionResult, LocalServiceSnapshot, PopupActivationIntent, SystemsSnapshot } from './model'

type NativeWindow = Window & {
  go?: { main?: { App?: { Snapshot(): Promise<ConnectionSnapshot>; RetryConnection(): Promise<ActionResult>; Quit(): Promise<ActionResult>; ConnectionProfiles?(): Promise<ConnectionResult>; SelectConnection?(id: string): Promise<ConnectionResult>; RenameConnection?(id: string, label: string): Promise<ConnectionResult>; RemoveConnection?(id: string): Promise<ConnectionResult>; AllSystems?(): Promise<SystemsSnapshot>; LocalServiceSnapshot?(): Promise<LocalServiceSnapshot>; ControlLocalService?(action: string, confirmed: boolean): Promise<LocalServiceActionResult> } } }
  runtime?: { EventsOn?(name: string, callback: (event: DesktopEvent | SystemsSnapshot | PopupActivationIntent) => void): () => void }
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
    allSystems: () => app?.AllSystems?.() ?? Promise.resolve({ generation: 0, startedAt: '', completedAt: '', complete: true, observations: [] }),
    localServiceSnapshot: () => app?.LocalServiceSnapshot?.() ?? Promise.resolve({ state: 'unsupported', scmState: '', detail: 'Local service controls are unavailable in this browser.', observedAt: '' }),
    controlLocalService: (action, confirmed) => app?.ControlLocalService?.(action, confirmed) ?? Promise.resolve({ action, outcome: 'unavailable', message: 'Local service controls are unavailable.', snapshot: { state: 'unsupported', scmState: '', detail: '', observedAt: '' } }),
    subscribeSystems: (listener) => nativeWindow.runtime?.EventsOn?.('systems:event', listener as (event: DesktopEvent | SystemsSnapshot | PopupActivationIntent) => void) ?? (() => undefined),
    quit: () => app?.Quit() ?? Promise.resolve({ action: 'quit', outcome: 'unavailable', message: 'Exit is available in the installed desktop application.' }),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener as (event: DesktopEvent | SystemsSnapshot | PopupActivationIntent) => void) ?? (() => undefined),
    subscribePopup: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:popup-activated', listener as (event: DesktopEvent | SystemsSnapshot | PopupActivationIntent) => void) ?? (() => undefined),
  }
}

export const desktopBridge = createBridge()
