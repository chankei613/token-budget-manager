import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'Token Budget Manager',
    'lang.toggle': 'JA',
    'nav.dashboard': 'Dashboard',
    'nav.budgets': 'Budgets',
    'nav.usage': 'Usage',
    'nav.pricing': 'Pricing',
    'nav.help': 'Help',
    'nav.settings': 'Settings',

    'error.prefix': 'Error: ',
    'error.retry': 'Retry',
    'loading': 'Loading…',

    'period.daily': 'Daily',
    'period.weekly': 'Weekly',
    'period.monthly': 'Monthly',
    'alert.ok': 'OK',
    'alert.warning': 'Warning',
    'alert.critical': 'Critical',

    'dashboard.title': 'Dashboard',
    'dashboard.empty': 'No budgets yet. Create one to start tracking spend.',
    'dashboard.tokensUsed': '{tokens} tokens · ${cost}',

    'budgets.title': 'Budgets',
    'budgets.empty': 'No budgets yet.',
    'budgets.new': 'New budget',
    'budgets.new.name': 'Budget name',
    'budgets.new.source': 'Source (blank = any)',
    'budgets.new.scopeKey': 'Scope key (blank = global)',
    'budgets.new.period': 'Period',
    'budgets.new.limitTokens': 'Token limit (0 = none)',
    'budgets.new.limitUsd': 'USD limit (0 = none)',
    'budgets.new.create': 'Create',
    'budgets.delete': 'Delete',
    'budgets.delete.confirm': 'Delete this budget? This cannot be undone.',

    'usage.title': 'Usage',
    'usage.empty': 'No usage recorded yet.',
    'usage.filter.source': 'Source',
    'usage.filter.agent': 'Agent',
    'usage.filter.scopeKey': 'Scope key',
    'usage.filter.model': 'Model',

    'pricing.title': 'Model pricing',
    'pricing.desc': 'Per-1M-token prices used to compute cost when usage is ingested. Edit to match your actual rates.',
    'pricing.model': 'Model',
    'pricing.input': 'Input $/1M',
    'pricing.output': 'Output $/1M',
    'pricing.save': 'Save',
    'pricing.new': 'Add model',

    'help.title': 'Help',
    'help.intro': 'How budgets, usage, and alerts fit together.',
    'help.what.title': 'What this app does',
    'help.what.body': 'AI systems POST their token usage here after each run. Cost is computed automatically from the model pricing table. Set budgets per project or agent, and get alerted the moment one crosses 80% (warning) or 100% (critical) of its limit — only once per period, not every time.',
    'help.start.title': 'Getting started',
    'help.start.1': 'Check Pricing and adjust the default $/1M rates if they don\'t match your actual contract.',
    'help.start.2': 'Create a Budget — leave Source/Scope key blank for a budget that covers everything, or fill them in to scope it to one project or system.',
    'help.start.3': 'Have your AI system POST to /api/v1/usage after each run (see Settings for the endpoint and API key).',
    'help.start.4': 'Watch Dashboard — a budget turns amber at 80% and red at 100%.',
    'help.stuck.title': 'Common snags',
    'help.stuck.1': 'Cost always shows $0 → the model_id sent doesn\'t match any row in Pricing; add it there.',
    'help.stuck.2': 'A budget never alerts → check Source/Scope key match what your AI system actually sends; an exact match is required unless left blank.',
    'help.stuck.3': 'Want a weekly or monthly report → open Usage and adjust the date range; the same summary powers both.',

    'settings.title': 'Settings',
    'settings.api.title': 'API endpoint',
    'settings.api.desc': 'AI systems POST usage here.',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet.',
    'settings.version': 'Version',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app? Ingestion will stop until you reopen it.',
  },
  ja: {
    'app.subtitle': 'Token Budget Manager',
    'lang.toggle': 'EN',
    'nav.dashboard': 'ダッシュボード',
    'nav.budgets': '予算',
    'nav.usage': '使用量',
    'nav.pricing': '価格設定',
    'nav.help': 'ヘルプ',
    'nav.settings': '設定',

    'error.prefix': 'エラー: ',
    'error.retry': '再試行',
    'loading': '読み込み中…',

    'period.daily': '日次',
    'period.weekly': '週次',
    'period.monthly': '月次',
    'alert.ok': '正常',
    'alert.warning': '警告',
    'alert.critical': '危険',

    'dashboard.title': 'ダッシュボード',
    'dashboard.empty': 'まだ予算がありません。作成して使用量を追跡しましょう。',
    'dashboard.tokensUsed': '{tokens}トークン · ${cost}',

    'budgets.title': '予算',
    'budgets.empty': 'まだ予算がありません。',
    'budgets.new': '新規作成',
    'budgets.new.name': '予算名',
    'budgets.new.source': 'Source（空欄=全て）',
    'budgets.new.scopeKey': 'Scope key（空欄=全体）',
    'budgets.new.period': '期間',
    'budgets.new.limitTokens': 'トークン上限（0=無制限）',
    'budgets.new.limitUsd': 'USD上限（0=無制限）',
    'budgets.new.create': '作成',
    'budgets.delete': '削除',
    'budgets.delete.confirm': 'この予算を削除しますか？元に戻せません。',

    'usage.title': '使用量',
    'usage.empty': 'まだ使用量が記録されていません。',
    'usage.filter.source': 'Source',
    'usage.filter.agent': 'Agent',
    'usage.filter.scopeKey': 'Scope key',
    'usage.filter.model': 'モデル',

    'pricing.title': 'モデル価格',
    'pricing.desc': '使用量記録時のコスト計算に使う$/1Mトークン単価。実際の契約に合わせて編集してください。',
    'pricing.model': 'モデル',
    'pricing.input': '入力 $/1M',
    'pricing.output': '出力 $/1M',
    'pricing.save': '保存',
    'pricing.new': 'モデルを追加',

    'help.title': 'ヘルプ',
    'help.intro': '予算・使用量・アラートがどう連動するかをまとめました。',
    'help.what.title': 'このアプリでできること',
    'help.what.body': 'AIシステムは実行のたびに使用量をここへPOSTします。コストはモデル価格表から自動計算されます。プロジェクトやエージェント単位で予算を設定すると、80%（警告）・100%（危険）を超えた瞬間にアラートされます — 毎回ではなく、今期で一度だけ。',
    'help.start.title': 'はじめに',
    'help.start.1': '価格設定を確認し、実際の契約と異なる場合はデフォルトの$/1M単価を調整してください。',
    'help.start.2': '予算を作成 — Source/Scope keyを空欄にすると全体を対象にした予算に、入力すると特定のプロジェクト/システムに絞った予算になります。',
    'help.start.3': 'AIシステム側から実行のたびに/api/v1/usageへPOSTするよう設定してください（エンドポイントとAPIキーは設定画面）。',
    'help.start.4': 'ダッシュボードを確認 — 80%で黄色、100%で赤になります。',
    'help.stuck.title': 'よくある詰まりどころ',
    'help.stuck.1': 'コストが常に$0になる → 送られたmodel_idが価格設定のどの行とも一致していません。価格設定に追加してください。',
    'help.stuck.2': '予算が全くアラートしない → AIシステム側が実際に送っているSource/Scope keyと一致しているか確認してください。空欄にしない限り完全一致が必要です。',
    'help.stuck.3': '週次・月次レポートが欲しい → 使用量ビューで日付範囲を調整してください。同じ集計が両方に使えます。',

    'settings.title': '設定',
    'settings.api.title': 'APIエンドポイント',
    'settings.api.desc': 'AIシステムはここへ使用量をPOSTします。',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。',
    'settings.version': 'バージョン',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？再度開くまでIngestionは停止します。',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
