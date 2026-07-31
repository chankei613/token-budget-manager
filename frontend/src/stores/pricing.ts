import { defineStore } from 'pinia'
import { ListPricing, SetPricing } from '../../wailsjs/go/main/App'
import { db } from '../../wailsjs/go/models'

export const usePricingStore = defineStore('pricing', {
  state: () => ({
    rows: [] as db.ModelPricing[],
    loading: false,
    saving: false,
    error: null as string | null,
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = null
      try {
        this.rows = (await ListPricing()) ?? []
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async save(modelId: string, inputPricePer1M: number, outputPricePer1M: number) {
      this.saving = true
      this.error = null
      try {
        await SetPricing(modelId, inputPricePer1M, outputPricePer1M)
        await this.load()
      } catch (e) {
        this.error = String(e)
      } finally {
        this.saving = false
      }
    },
  },
})
