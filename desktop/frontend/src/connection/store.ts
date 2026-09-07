import { useCallback, useEffect, useState } from 'react'
import { desktopBridge, unavailableSnapshot } from './bridge'
import type { ConnectionSnapshot, DesktopBridge } from './model'

export function acceptSnapshot(current: ConnectionSnapshot, incoming: ConnectionSnapshot): ConnectionSnapshot {
  if (incoming.generation < current.generation) return current
  if (incoming.generation === current.generation && incoming.revision < current.revision) return current
  return incoming
}

export function useConnection(bridge: DesktopBridge = desktopBridge) {
  const [snapshot, setSnapshot] = useState<ConnectionSnapshot>(unavailableSnapshot)
  const [announcement, setAnnouncement] = useState('')

  useEffect(() => {
    let active = true
    void bridge.snapshot().then((incoming) => { if (active) setSnapshot((current) => acceptSnapshot(current, incoming)) })
    const unsubscribe = bridge.subscribe((event) => {
      if (!active) return
      setSnapshot((current) => {
        if (event.generation < current.generation) return current
        setAnnouncement(event.message)
        return event.snapshot ? acceptSnapshot(current, event.snapshot) : current
      })
    })
    return () => { active = false; unsubscribe() }
  }, [bridge])

  const retry = useCallback(async () => {
    const result = await bridge.retry()
    setAnnouncement(result.message)
  }, [bridge])

  return { snapshot, announcement, retry }
}
