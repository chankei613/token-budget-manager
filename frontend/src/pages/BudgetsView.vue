<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useBudgetsStore } from '@/stores/budgets'
import { useI18n } from '@/i18n'
import { alertColor } from '@/alertColors'

const { t } = useI18n()
const store = useBudgetsStore()

const name = ref('')
const source = ref('')
const scopeKey = ref('')
const period = ref('monthly')
const limitTokens = ref(0)
const limitUsd = ref(0)
const creating = ref(false)

onMounted(() => store.list())

async function create() {
  if (!name.value.trim()) return
  creating.value = true
  await store.create(name.value.trim(), source.value.trim(), scopeKey.value.trim(), period.value, limitTokens.value, limitUsd.value)
  creating.value = false
  name.value = ''
  source.value = ''
  scopeKey.value = ''
  limitTokens.value = 0
  limitUsd.value = 0
}

async function remove(id: string) {
  if (!confirm(t('budgets.delete.confirm'))) return
  await store.remove(id)
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('budgets.title') }}</h2>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.list" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div class="border border-border rounded-lg p-4 space-y-3 max-w-lg">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('budgets.new') }}</h3>
      <input v-model="name" :placeholder="t('budgets.new.name')" class="w-full text-sm border border-border rounded px-2 py-1.5" />
      <div class="grid grid-cols-2 gap-2">
        <input v-model="source" :placeholder="t('budgets.new.source')" class="text-sm border border-border rounded px-2 py-1.5" />
        <input v-model="scopeKey" :placeholder="t('budgets.new.scopeKey')" class="text-sm border border-border rounded px-2 py-1.5" />
      </div>
      <div class="grid grid-cols-3 gap-2">
        <select v-model="period" class="text-sm border border-border rounded px-2 py-1.5">
          <option value="daily">{{ t('period.daily') }}</option>
          <option value="weekly">{{ t('period.weekly') }}</option>
          <option value="monthly">{{ t('period.monthly') }}</option>
        </select>
        <input v-model.number="limitTokens" type="number" min="0" :placeholder="t('budgets.new.limitTokens')" class="text-sm border border-border rounded px-2 py-1.5" />
        <input v-model.number="limitUsd" type="number" min="0" step="0.01" :placeholder="t('budgets.new.limitUsd')" class="text-sm border border-border rounded px-2 py-1.5" />
      </div>
      <button @click="create" :disabled="creating || !name.trim()" class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40">
        {{ t('budgets.new.create') }}
      </button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.budgets.length === 0" class="text-sm text-muted-foreground">{{ t('budgets.empty') }}</div>

    <div v-else class="space-y-1.5">
      <div
        v-for="b in store.budgets"
        :key="b.id"
        class="flex items-center gap-3 border border-border rounded-lg px-4 py-2.5 border-l-4"
        :style="{ borderLeftColor: store.statuses[b.id] ? alertColor(store.statuses[b.id].alert_level) : '#e5e5e5' }"
      >
        <div class="flex-1">
          <div class="text-sm font-medium">{{ b.name }}</div>
          <div class="text-xs text-muted-foreground">
            {{ t('period.' + b.period) }}
            <span v-if="b.source"> · source={{ b.source }}</span>
            <span v-if="b.scope_key"> · scope={{ b.scope_key }}</span>
            <span v-if="b.limit_tokens"> · {{ b.limit_tokens }} tokens</span>
            <span v-if="b.limit_usd"> · ${{ b.limit_usd }}</span>
          </div>
        </div>
        <span v-if="store.statuses[b.id]" class="text-xs tabular-nums text-muted-foreground">
          {{ (store.statuses[b.id].percent_used * 100).toFixed(0) }}%
        </span>
        <button @click="remove(b.id)" class="text-xs text-red-600 hover:underline">{{ t('budgets.delete') }}</button>
      </div>
    </div>
  </div>
</template>
