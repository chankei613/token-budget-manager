import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useBudgetsStore } from '@/stores/budgets'
import { ListBudgets } from '../../wailsjs/go/main/App'

describe('budgets store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(ListBudgets).mockReset()
  })

  it('captures a failed list() as store.error and clears loading', async () => {
    vi.mocked(ListBudgets).mockRejectedValueOnce(new Error('network down'))
    const store = useBudgetsStore()

    await store.list()

    expect(store.loading).toBe(false)
    expect(store.error).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(ListBudgets).mockRejectedValueOnce(new Error('network down'))
    const store = useBudgetsStore()
    await store.list()
    expect(store.error).not.toBeNull()

    vi.mocked(ListBudgets).mockResolvedValueOnce([])
    await store.list()

    expect(store.error).toBeNull()
  })
})
