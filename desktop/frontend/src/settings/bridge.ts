import type { Appearance } from '../connection/model'
import type { SettingsBridge, SettingsResult } from './model'

type NativeWindow = Window & { go?: { main?: { App?: {
  SettingsWorkspace?(): Promise<SettingsResult>
  SaveAppearance?(value: string): Promise<SettingsResult>
  RestoreDesktopPreferences?(): Promise<SettingsResult>
  CopyStoragePath?(id: string): Promise<SettingsResult>
  OpenProductLink?(key: string): Promise<SettingsResult>
} } } }

const unavailable = (action: string): SettingsResult => ({ action, outcome: 'unavailable', message: 'Desktop settings are available in the installed application.' })

export function createSettingsBridge(nativeWindow: NativeWindow = window): SettingsBridge {
  const app = nativeWindow.go?.main?.App
  return {
    workspace: () => app?.SettingsWorkspace?.() ?? Promise.resolve(unavailable('load_settings')),
    saveAppearance: (value: Appearance) => app?.SaveAppearance?.(value) ?? Promise.resolve(unavailable('save_appearance')),
    restore: () => app?.RestoreDesktopPreferences?.() ?? Promise.resolve(unavailable('restore_preferences')),
    copyStoragePath: (id) => app?.CopyStoragePath?.(id) ?? Promise.resolve(unavailable('copy_storage_path')),
    openProductLink: (key) => app?.OpenProductLink?.(key) ?? Promise.resolve(unavailable('open_product_link')),
  }
}

export const settingsBridge = createSettingsBridge()
