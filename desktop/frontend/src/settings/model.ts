import type { Appearance } from '../connection/model'

export interface PreferenceTransition { status: 'migrated' | 'not_found' | 'invalid' | 'unreadable' | 'not_required'; retired: string[] }
export interface DesktopPreferences { version: number; appearance: Appearance; transition: PreferenceTransition }
export interface StorageRecord { id: string; label: string; path?: string; owner: 'go-schedule' | 'desktop' | 'operating-system' | 'external'; scope: 'machine' | 'user' | 'runtime' | 'external'; existence: 'present' | 'absent' | 'unavailable'; normalRemoval: string; explicitWipe: string; copyable: boolean }
export interface ProductLink { key: string; label: string; destination: string }
export interface ProductInformation { name: string; version: string; publisher: string; links: ProductLink[] }
export interface SettingsWorkspace { preferences: DesktopPreferences; preferencePath: string; storage: StorageRecord[]; product: ProductInformation; daemonAvailable: boolean; loadedAt: string }
export interface SettingsResult { action: string; outcome: 'accepted' | 'rejected' | 'unavailable'; message: string; workspace?: SettingsWorkspace }
export interface SettingsBridge {
  workspace(): Promise<SettingsResult>
  saveAppearance(value: Appearance): Promise<SettingsResult>
  restore(): Promise<SettingsResult>
  copyStoragePath(id: string): Promise<SettingsResult>
  openProductLink(key: string): Promise<SettingsResult>
}
