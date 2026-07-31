// Package db はToken Budget ManagerのGORMモデルとSQLite初期化を提供する。
// docs/spec.md参照。
package db

import "time"

// UsageEvent は1回の実行で消費したトークン・コストの記録（追記専用）。
type UsageEvent struct {
	ID         string    `gorm:"primaryKey" json:"id"`
	ReceivedAt time.Time `gorm:"index" json:"received_at"`

	Source   string `gorm:"index" json:"source"`
	AgentID  string `gorm:"index" json:"agent_id"`
	ScopeKey string `gorm:"index" json:"scope_key"`

	Provider string `json:"provider"`
	ModelID  string `gorm:"index" json:"model_id"`

	InputTokens  int64   `json:"input_tokens"`
	OutputTokens int64   `json:"output_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

// ModelPricing はモデルごとの$/1Mトークン単価。
type ModelPricing struct {
	ModelID          string  `gorm:"primaryKey" json:"model_id"`
	InputPricePer1M  float64 `json:"input_price_per_1m"`
	OutputPricePer1M float64 `json:"output_price_per_1m"`
}

type BudgetPeriod string

const (
	PeriodDaily   BudgetPeriod = "daily"
	PeriodWeekly  BudgetPeriod = "weekly"
	PeriodMonthly BudgetPeriod = "monthly"
)

type AlertLevel string

const (
	AlertNone     AlertLevel = ""
	AlertWarning  AlertLevel = "warning"
	AlertCritical AlertLevel = "critical"
)

// Budget はSource×ScopeKeyに対する期間ごとの上限。使用量は都度SUMで算出し、
// カウンタは持たない。LastAlertLevel/AlertPeriodStartで「今期で既に出した
// アラートの最高値」を憶え、二重アラートを防ぐ。
type Budget struct {
	ID       string       `gorm:"primaryKey" json:"id"`
	Name     string       `json:"name"`
	Source   string       `gorm:"index" json:"source"`
	ScopeKey string       `gorm:"index" json:"scope_key"`
	Period   BudgetPeriod `json:"period"`

	LimitTokens int64   `json:"limit_tokens"`
	LimitUSD    float64 `json:"limit_usd"`

	LastAlertLevel   AlertLevel `json:"last_alert_level"`
	AlertPeriodStart time.Time  `json:"alert_period_start"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AgentKey — Ingestion/CRUD APIを叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

// DefaultPricing は主要モデルの目安単価（2026-07時点、$/1Mトークン）。
// ユーザーがUIから編集できる前提のデフォルト値であり、正確な最新価格の追従はしない。
func DefaultPricing() []ModelPricing {
	return []ModelPricing{
		{ModelID: "claude-opus-5", InputPricePer1M: 15, OutputPricePer1M: 75},
		{ModelID: "claude-sonnet-5", InputPricePer1M: 3, OutputPricePer1M: 15},
		{ModelID: "claude-haiku-4-5", InputPricePer1M: 0.8, OutputPricePer1M: 4},
		{ModelID: "gpt-4o", InputPricePer1M: 2.5, OutputPricePer1M: 10},
		{ModelID: "gpt-4o-mini", InputPricePer1M: 0.15, OutputPricePer1M: 0.6},
		{ModelID: "gemini-1.5-pro", InputPricePer1M: 1.25, OutputPricePer1M: 5},
	}
}
