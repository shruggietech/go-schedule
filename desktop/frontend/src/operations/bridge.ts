import type { OperationResult, OperationsBridge } from './model'

type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<OperationResult>> } }; runtime?: { EventsOn?(name: string, callback: (event: { kind: string }) => void): () => void } }
const unavailable = (action: string): OperationResult => ({ action, outcome: 'unavailable', message: 'Operational data is available in the installed desktop application.' })

export function createOperationsBridge(nativeWindow: NativeWindow = window): OperationsBridge {
  const app = nativeWindow.go?.main?.App
  const call = (name: string, action: string, ...args: unknown[]) => app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))
  return {
    scheduleWindow: (days) => call('ScheduleWindow', 'load_schedule', days),
    activityWorkspace: () => call('ActivityWorkspace', 'load_activity'),
    acknowledgeAlert: (id) => call('AcknowledgeAlert', 'acknowledge_alerts', id),
    acknowledgeAlerts: (ids) => call('AcknowledgeAlerts', 'acknowledge_alerts', ids),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined),
  }
}

export const operationsBridge = createOperationsBridge()
