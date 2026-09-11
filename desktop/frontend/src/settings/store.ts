import { useCallback, useEffect, useRef, useState } from 'react'
import type { Appearance } from '../connection/model'
import { settingsBridge } from './bridge'
import type { SettingsBridge, SettingsWorkspace } from './model'

export function useSettings(bridge: SettingsBridge = settingsBridge, refreshToken: string | number = 0) {
  const [workspace, setWorkspace] = useState<SettingsWorkspace>()
  const [message, setMessage] = useState('')
  const [pending, setPending] = useState(false)
  const [loading, setLoading] = useState(true)
  const pendingRef = useRef(false)

  useEffect(() => {
    let active = true
    void bridge.workspace().then((result) => {
      if (!active) return
      if (result.outcome === 'accepted' && result.workspace) setWorkspace(result.workspace)
      setLoading(false)
    }, () => {
      if (!active) return
      setMessage('Desktop settings are unavailable. Check access to the user configuration directory, then try again.')
      setLoading(false)
    })
    return () => { active = false }
  }, [bridge, refreshToken])

  const run = useCallback(async (operation: () => ReturnType<SettingsBridge['workspace']>) => {
    if (pendingRef.current) return undefined
    pendingRef.current = true
    setPending(true)
    try {
      const result = await operation()
      setMessage(result.message)
      if (result.outcome === 'accepted' && result.workspace) setWorkspace(result.workspace)
      return result
    } catch {
      setMessage('The desktop action could not be completed. Try again.')
      return undefined
    } finally {
      pendingRef.current = false
      setPending(false)
    }
  }, [])

  const saveAppearance = useCallback((value: Appearance) => run(() => bridge.saveAppearance(value)), [bridge, run])
  const restore = useCallback(() => run(() => bridge.restore()), [bridge, run])
  const copyStoragePath = useCallback((id: string) => run(() => bridge.copyStoragePath(id)), [bridge, run])
  const openProductLink = useCallback((key: string) => run(() => bridge.openProductLink(key)), [bridge, run])

  return { workspace, message, loading, pending, saveAppearance, restore, copyStoragePath, openProductLink }
}
