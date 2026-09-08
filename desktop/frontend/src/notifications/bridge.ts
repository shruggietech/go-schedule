import type { ChannelDraft, NotificationBridge, NotificationResult, PolicyDraft } from './model'

type NativeWindow = Window & { go?: { main?: { App?: Record<string, (...args: unknown[]) => Promise<NotificationResult>> } }; runtime?: { EventsOn?(name: string, callback: (event: { kind: string }) => void): () => void } }
const unavailable = (action: string): NotificationResult => ({ action, outcome: 'unavailable', message: 'Notifications are available in the installed desktop application.' })

export function createNotificationBridge(nativeWindow: NativeWindow = window): NotificationBridge {
  const app = nativeWindow.go?.main?.App
  const call = (name: string, action: string, ...args: unknown[]) => app?.[name]?.(...args) ?? Promise.resolve(unavailable(action))
  return {
    workspace: () => call('NotificationWorkspace', 'load_notifications'),
    saveChannel: (draft: ChannelDraft) => call('SaveNotificationChannel', 'save_notification_channel', draft),
    setChannelEnabled: (id, enabled) => call('SetNotificationChannelEnabled', 'toggle_notification_channel', id, enabled),
    testChannel: (id) => call('TestNotificationChannel', 'test_notification_channel', id),
    deleteChannel: (id) => call('DeleteNotificationChannel', 'delete_notification_channel', id),
    policy: (scopeType, scopeId) => call('NotificationPolicy', 'load_notification_policy', scopeType, scopeId),
    savePolicy: (draft: PolicyDraft) => call('SaveNotificationPolicy', 'save_notification_policy', draft),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('desktop:event', listener) ?? (() => undefined),
  }
}

export const notificationBridge = createNotificationBridge()
