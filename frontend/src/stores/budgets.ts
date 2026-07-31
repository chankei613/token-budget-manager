import { defineStore } from 'pinia'
import { ListBudgets, CreateBudget, DeleteBudget, BudgetStatus } from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export const useBudgetsStore = defineStore('budgets', {
  state: () => ({
    budgets: [] as db.Budget[],
    statuses: {} as Record<string, api.BudgetStatus>,
    loading: false,
    error: null as string | null,
  }),
  actions: {
    async list() {
      this.loading = true
      this.error = null
      try {
        this.budgets = (await ListBudgets()) ?? []
        await Promise.all(this.budgets.map((b) => this.loadStatus(b.id)))
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async loadStatus(id: string) {
      try {
        this.statuses[id] = await BudgetStatus(id)
      } catch {
        // 個別のstatus取得失敗は一覧全体を壊さない
      }
    },
    async create(
      name: string,
      source: string,
      scopeKey: string,
      period: string,
      limitTokens: number,
      limitUsd: number,
    ) {
      this.error = null
      try {
        await CreateBudget(name, source, scopeKey, period, limitTokens, limitUsd)
        await this.list()
      } catch (e) {
        this.error = String(e)
      }
    },
    async remove(id: string) {
      this.error = null
      try {
        await DeleteBudget(id)
        await this.list()
      } catch (e) {
        this.error = String(e)
      }
    },
  },
})
