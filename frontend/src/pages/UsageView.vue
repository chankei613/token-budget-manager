<script setup lang="ts">
import { onMounted } from 'vue'
import { useUsageStore } from '@/stores/usage'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useUsageStore()

function reload() {
  store.offset = 0
  store.search()
  store.loadSummary()
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}

onMounted(reload)
</script>

<template>
  <div class="flex h-full">
    <aside class="w-56 border-r border-border p-4 space-y-3 shrink-0">
      <h2 class="text-sm font-semibold">{{ t('usage.title') }}</h2>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('usage.filter.source') }}</label>
        <input v-model="store.filters.source" @change="reload" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('usage.filter.agent') }}</label>
        <input v-model="store.filters.agentId" @change="reload" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('usage.filter.scopeKey') }}</label>
        <input v-model="store.filters.scopeKey" @change="reload" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>
      <div class="space-y-1">
        <label class="text-xs text-muted-foreground">{{ t('usage.filter.model') }}</label>
        <input v-model="store.filters.modelId" @change="reload" class="w-full text-sm border border-border rounded px-2 py-1" />
      </div>

      <div v-if="store.summary" class="pt-2 border-t border-border space-y-1 text-xs">
        <div class="flex justify-between"><span class="text-muted-foreground">Tokens</span><span class="tabular-nums">{{ store.summary.total_input_tokens + store.summary.total_output_tokens }}</span></div>
        <div class="flex justify-between"><span class="text-muted-foreground">Cost</span><span class="tabular-nums">${{ store.summary.total_cost_usd.toFixed(2) }}</span></div>
      </div>
    </aside>

    <main class="flex-1 overflow-y-auto p-4">
      <div v-if="store.error" class="text-sm border rounded px-3 py-2 mb-3 border-red-300 text-red-600">
        {{ t('error.prefix') }}{{ store.error }}
        <button @click="reload" class="ml-2 underline">{{ t('error.retry') }}</button>
      </div>

      <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
      <div v-else-if="store.items.length === 0" class="text-sm text-muted-foreground">{{ t('usage.empty') }}</div>

      <table v-else class="w-full text-xs">
        <thead>
          <tr class="text-left text-muted-foreground border-b border-border">
            <th class="py-1.5 pr-3 font-medium">Time</th>
            <th class="py-1.5 pr-3 font-medium">Agent</th>
            <th class="py-1.5 pr-3 font-medium">Model</th>
            <th class="py-1.5 pr-3 font-medium text-right">Input</th>
            <th class="py-1.5 pr-3 font-medium text-right">Output</th>
            <th class="py-1.5 font-medium text-right">Cost</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="e in store.items" :key="e.id" class="border-b border-border/50">
            <td class="py-1.5 pr-3 text-muted-foreground">{{ fmt(e.received_at) }}</td>
            <td class="py-1.5 pr-3 font-medium">{{ e.agent_id }}</td>
            <td class="py-1.5 pr-3 font-mono">{{ e.model_id }}</td>
            <td class="py-1.5 pr-3 text-right tabular-nums">{{ e.input_tokens }}</td>
            <td class="py-1.5 pr-3 text-right tabular-nums">{{ e.output_tokens }}</td>
            <td class="py-1.5 text-right tabular-nums">${{ e.cost_usd.toFixed(4) }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="store.total > 0" class="flex items-center justify-between pt-3 text-xs text-muted-foreground">
        <button :disabled="store.offset === 0" @click="store.prevPage" class="disabled:opacity-30 underline">Prev</button>
        <span>{{ store.offset + 1 }}–{{ Math.min(store.offset + store.limit, store.total) }} / {{ store.total }}</span>
        <button :disabled="store.offset + store.limit >= store.total" @click="store.nextPage" class="disabled:opacity-30 underline">Next</button>
      </div>
    </main>
  </div>
</template>
