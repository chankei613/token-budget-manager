<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { usePricingStore } from '@/stores/pricing'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = usePricingStore()

const newModelId = ref('')
const newInput = ref(0)
const newOutput = ref(0)

onMounted(() => store.load())

async function addModel() {
  if (!newModelId.value.trim()) return
  await store.save(newModelId.value.trim(), newInput.value, newOutput.value)
  newModelId.value = ''
  newInput.value = 0
  newOutput.value = 0
}
</script>

<template>
  <div class="p-6 space-y-4 max-w-2xl overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('pricing.title') }}</h2>
    <p class="text-xs text-muted-foreground">{{ t('pricing.desc') }}</p>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.load" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

    <table v-else class="w-full text-sm">
      <thead>
        <tr class="text-left text-xs text-muted-foreground border-b border-border">
          <th class="py-1.5 pr-3 font-medium">{{ t('pricing.model') }}</th>
          <th class="py-1.5 pr-3 font-medium">{{ t('pricing.input') }}</th>
          <th class="py-1.5 font-medium">{{ t('pricing.output') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="row in store.rows" :key="row.model_id" class="border-b border-border/50">
          <td class="py-1.5 pr-3 font-mono text-xs">{{ row.model_id }}</td>
          <td class="py-1.5 pr-3">
            <input
              type="number" min="0" step="0.01"
              :value="row.input_price_per_1m"
              @change="store.save(row.model_id, ($event.target as HTMLInputElement).valueAsNumber, row.output_price_per_1m)"
              class="w-24 text-sm border border-border rounded px-2 py-1"
            />
          </td>
          <td class="py-1.5">
            <input
              type="number" min="0" step="0.01"
              :value="row.output_price_per_1m"
              @change="store.save(row.model_id, row.input_price_per_1m, ($event.target as HTMLInputElement).valueAsNumber)"
              class="w-24 text-sm border border-border rounded px-2 py-1"
            />
          </td>
        </tr>
      </tbody>
    </table>

    <div class="flex items-end gap-2 pt-2 border-t border-border">
      <input v-model="newModelId" placeholder="model_id" class="flex-1 text-sm border border-border rounded px-2 py-1.5 font-mono" />
      <input v-model.number="newInput" type="number" min="0" step="0.01" class="w-24 text-sm border border-border rounded px-2 py-1.5" />
      <input v-model.number="newOutput" type="number" min="0" step="0.01" class="w-24 text-sm border border-border rounded px-2 py-1.5" />
      <button @click="addModel" class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white shrink-0">{{ t('pricing.new') }}</button>
    </div>
  </div>
</template>
