<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useBudgetsStore } from '@/stores/budgets'
import { useI18n } from '@/i18n'
import { alertColor } from '@/alertColors'
import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime'

const { t } = useI18n()
const store = useBudgetsStore()

onMounted(() => {
  store.list()
  EventsOn('budget:alert', () => store.list())
})
onUnmounted(() => EventsOff('budget:alert'))

const rows = computed(() =>
  store.budgets.map((b) => {
    const status = store.statuses[b.id]
    return { budget: b, status }
  }),
)
</script>

<template>
  <div class="p-6 space-y-4 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('dashboard.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.list" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.budgets.length === 0" class="text-sm text-muted-foreground">{{ t('dashboard.empty') }}</div>

    <div v-else class="grid gap-3" style="grid-template-columns: repeat(auto-fill, minmax(280px, 1fr))">
      <div v-for="row in rows" :key="row.budget.id" class="border border-border rounded-lg p-4 space-y-2">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">{{ row.budget.name }}</span>
          <span class="text-xs px-1.5 py-0.5 rounded bg-gray-100 text-gray-600">{{ t('period.' + row.budget.period) }}</span>
        </div>
        <template v-if="row.status">
          <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full"
              :style="{ width: Math.min(100, row.status.percent_used * 100) + '%', background: alertColor(row.status.alert_level) }"
            />
          </div>
          <div class="flex items-center justify-between text-xs text-muted-foreground">
            <span>{{ t('dashboard.tokensUsed', { tokens: row.status.tokens_used, cost: row.status.cost_used_usd.toFixed(2) }) }}</span>
            <span :style="{ color: alertColor(row.status.alert_level) }" class="font-medium">
              {{ t('alert.' + (row.status.alert_level || 'ok')) }}
            </span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>
