import { vi } from 'vitest'

// Node 25+ の実験的 localStorage がhappy-domのwindow.localStorageと衝突し
// `--localstorage-file` 未指定だと getItem 等がthrowする既知の非互換。
// メモリ上の簡易実装に差し替えて回避する。
class MemoryStorage implements Storage {
  private store = new Map<string, string>()
  get length() { return this.store.size }
  clear() { this.store.clear() }
  getItem(key: string) { return this.store.has(key) ? this.store.get(key)! : null }
  key(index: number) { return Array.from(this.store.keys())[index] ?? null }
  removeItem(key: string) { this.store.delete(key) }
  setItem(key: string, value: string) { this.store.set(key, String(value)) }
}
const memoryStorage = new MemoryStorage()
Object.defineProperty(globalThis, 'localStorage', { value: memoryStorage, configurable: true })
Object.defineProperty(window, 'localStorage', { value: memoryStorage, configurable: true })

// Wails ランタイムのモック — テスト環境では no-op
vi.mock('../../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn(),
  EventsOff: vi.fn(),
  EventsEmit: vi.fn(),
  BrowserOpenURL: vi.fn(),
}))

vi.mock('../../wailsjs/go/main/App', () => ({
  GetAppVersion: vi.fn().mockResolvedValue('0.1.0'),
  GetAPIURL: vi.fn().mockResolvedValue('http://127.0.0.1:8424'),
  Quit: vi.fn().mockResolvedValue(undefined),
  ListBudgets: vi.fn().mockResolvedValue([]),
  CreateBudget: vi.fn().mockResolvedValue({ id: '1', name: 'test' }),
  UpdateBudget: vi.fn().mockResolvedValue({ id: '1', name: 'test' }),
  DeleteBudget: vi.fn().mockResolvedValue(undefined),
  BudgetStatus: vi.fn().mockResolvedValue({ budget: { id: '1' }, tokens_used: 0, cost_used_usd: 0, percent_used: 0, alert_level: '' }),
  ListUsage: vi.fn().mockResolvedValue({ items: [], total: 0 }),
  Summary: vi.fn().mockResolvedValue({ total_input_tokens: 0, total_output_tokens: 0, total_cost_usd: 0, by_agent_cost_usd: {}, by_model_cost_usd: {}, daily_series: [] }),
  ListPricing: vi.fn().mockResolvedValue([]),
  SetPricing: vi.fn().mockResolvedValue({ model_id: 'x', input_price_per_1m: 0, output_price_per_1m: 0 }),
  ListKeys: vi.fn().mockResolvedValue([]),
  IssueKey: vi.fn().mockResolvedValue({ id: '1', name: 'test', api_key: 'test-key' }),
  RevokeKey: vi.fn().mockResolvedValue(undefined),
}))
