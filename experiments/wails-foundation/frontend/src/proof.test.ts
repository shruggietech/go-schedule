import { describe, expect, it } from 'vitest'
import { fixtureFor } from './proof'

describe('proof fixtures', () => {
  it('preserves target identity for every connection condition', () => {
    for (const condition of ['connected', 'disconnected', 'degraded', 'loading'] as const) {
      const snapshot = fixtureFor(condition)
      expect(snapshot.target.displayName).toBe('This computer')
      expect(snapshot.target.connection).toBe(condition)
    }
  })

  it('provides deterministic empty data', () => {
    expect(fixtureFor('empty').tasks).toEqual([])
  })
})
