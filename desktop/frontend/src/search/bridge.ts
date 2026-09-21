import type { SearchActionResult, SearchBridge, SearchIntent, SearchRequest, SearchSnapshot } from './model'

type NativeWindow = Window & {
  go?: { main?: { App?: { SearchAcrossSystems?(request: SearchRequest): Promise<SearchSnapshot>; ExecuteSearchAction?(intent: SearchIntent): Promise<SearchActionResult> } } }
  runtime?: { EventsOn?(name: string, callback: (snapshot: SearchSnapshot) => void): () => void }
}

export function createSearchBridge(nativeWindow: NativeWindow = window): SearchBridge {
  const app = nativeWindow.go?.main?.App
  return {
    search: (request) => app?.SearchAcrossSystems?.(request) ?? Promise.resolve({ generation: 0, query: request.query, startedAt: '', complete: true, observations: [] }),
    execute: (intent) => app?.ExecuteSearchAction?.(intent) ?? Promise.resolve({ action: intent.action, outcome: 'unavailable', message: 'Cross-daemon actions are available in the installed desktop application.', outcomes: [] }),
    subscribe: (listener) => nativeWindow.runtime?.EventsOn?.('search:event', listener) ?? (() => undefined),
  }
}

export const searchBridge = createSearchBridge()
