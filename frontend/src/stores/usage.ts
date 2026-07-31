import { defineStore } from 'pinia'
import { ListUsage, Summary } from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export const useUsageStore = defineStore('usage', {
  state: () => ({
    filters: { source: '', agentId: '', scopeKey: '', modelId: '' },
    items: [] as db.UsageEvent[],
    total: 0,
    limit: 50,
    offset: 0,
    summary: null as api.UsageSummary | null,
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async search() {
      this.loading = true
      this.error = null
      try {
        const result = await ListUsage(
          this.filters.source,
          this.filters.agentId,
          this.filters.scopeKey,
          this.filters.modelId,
          this.limit,
          this.offset,
        )
        this.items = result.items ?? []
        this.total = result.total
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async loadSummary() {
      this.error = null
      try {
        this.summary = await Summary(this.filters.source, this.filters.agentId, this.filters.scopeKey, this.filters.modelId)
      } catch (e) {
        this.error = String(e)
      }
    },
    nextPage() {
      if (this.offset + this.limit < this.total) {
        this.offset += this.limit
        this.search()
      }
    },
    prevPage() {
      if (this.offset > 0) {
        this.offset = Math.max(0, this.offset - this.limit)
        this.search()
      }
    },
  },
})
