// アラートレベル配色は dataviz スキルの検証済みパレット（status palette）をそのまま使う。
export type AlertLevel = '' | 'warning' | 'critical'

export const ALERT_COLORS: Record<AlertLevel, string> = {
  '': '#0ca30c',
  warning: '#fab219',
  critical: '#d03b3b',
}

export function alertColor(level: string): string {
  return ALERT_COLORS[(level as AlertLevel) ?? ''] ?? ALERT_COLORS['']
}
