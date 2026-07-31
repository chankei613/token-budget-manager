package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/chankei613/token-budget-manager/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IngestUsageInput struct {
	Source       string `json:"source"`
	AgentID      string `json:"agent_id"`
	ScopeKey     string `json:"scope_key"`
	Provider     string `json:"provider"`
	ModelID      string `json:"model_id"`
	InputTokens  int64  `json:"input_tokens"`
	OutputTokens int64  `json:"output_tokens"`
}

// IngestUsage は1件のUsageEventを追加する（HTTP・ネイティブバインディング共用）。
// モデル価格が登録されていれば自動でCostUSDを計算する。追加後、影響するBudgetの
// アラート状態を再評価する（Phase 2）。
func (s *Server) IngestUsage(in IngestUsageInput) (db.UsageEvent, error) {
	if in.AgentID == "" || in.ModelID == "" {
		return db.UsageEvent{}, &apiError{"agent_id and model_id are required"}
	}

	cost := 0.0
	if price, ok := s.lookupPricing(in.ModelID); ok {
		cost = computeCost(in.InputTokens, in.OutputTokens, price)
	}

	event := db.UsageEvent{
		ID:           uuid.NewString(),
		ReceivedAt:   time.Now(),
		Source:       in.Source,
		AgentID:      in.AgentID,
		ScopeKey:     in.ScopeKey,
		Provider:     in.Provider,
		ModelID:      in.ModelID,
		InputTokens:  in.InputTokens,
		OutputTokens: in.OutputTokens,
		CostUSD:      cost,
	}
	if err := s.DB.Create(&event).Error; err != nil {
		return db.UsageEvent{}, err
	}

	s.evaluateBudgetsForEvent(event)
	return event, nil
}

type ListUsageFilters struct {
	Source   string
	AgentID  string
	ScopeKey string
	ModelID  string
	From     *time.Time
	To       *time.Time
}

func usageFiltersFromQuery(q url.Values) ListUsageFilters {
	f := ListUsageFilters{
		Source:   q.Get("source"),
		AgentID:  q.Get("agent_id"),
		ScopeKey: q.Get("scope_key"),
		ModelID:  q.Get("model_id"),
	}
	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.From = &t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.To = &t
		}
	}
	return f
}

func (f ListUsageFilters) apply(q *gorm.DB) *gorm.DB {
	if f.Source != "" {
		q = q.Where("source = ?", f.Source)
	}
	if f.AgentID != "" {
		q = q.Where("agent_id = ?", f.AgentID)
	}
	if f.ScopeKey != "" {
		q = q.Where("scope_key = ?", f.ScopeKey)
	}
	if f.ModelID != "" {
		q = q.Where("model_id = ?", f.ModelID)
	}
	if f.From != nil {
		q = q.Where("received_at >= ?", *f.From)
	}
	if f.To != nil {
		q = q.Where("received_at <= ?", *f.To)
	}
	return q
}

type ListUsageResult struct {
	Items []db.UsageEvent `json:"items"`
	Total int64           `json:"total"`
}

func (s *Server) ListUsage(f ListUsageFilters, limit, offset int) (ListUsageResult, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}

	q := f.apply(s.DB.Model(&db.UsageEvent{}))

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return ListUsageResult{}, err
	}

	var items []db.UsageEvent
	if err := q.Order("received_at desc").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return ListUsageResult{}, err
	}
	return ListUsageResult{Items: items, Total: total}, nil
}

type UsageSummary struct {
	TotalInputTokens  int64              `json:"total_input_tokens"`
	TotalOutputTokens int64              `json:"total_output_tokens"`
	TotalCostUSD      float64            `json:"total_cost_usd"`
	ByAgent           map[string]float64 `json:"by_agent_cost_usd"`
	ByModel           map[string]float64 `json:"by_model_cost_usd"`
	DailySeries       []DailyUsagePoint  `json:"daily_series"`
}

type DailyUsagePoint struct {
	Date    string  `json:"date"`
	CostUSD float64 `json:"cost_usd"`
	Tokens  int64   `json:"tokens"`
}

// Summary は週次・月次レポートを兼ねる集計。fromとtoの範囲を変えるだけで
// 週次・月次どちらの粒度のレポートにも使える。
func (s *Server) Summary(f ListUsageFilters) (UsageSummary, error) {
	summary := UsageSummary{ByAgent: map[string]float64{}, ByModel: map[string]float64{}}

	q := f.apply(s.DB.Model(&db.UsageEvent{}))

	var events []db.UsageEvent
	if err := q.Find(&events).Error; err != nil {
		return UsageSummary{}, err
	}

	byDate := map[string]*DailyUsagePoint{}
	for _, e := range events {
		summary.TotalInputTokens += e.InputTokens
		summary.TotalOutputTokens += e.OutputTokens
		summary.TotalCostUSD += e.CostUSD
		summary.ByAgent[e.AgentID] += e.CostUSD
		summary.ByModel[e.ModelID] += e.CostUSD

		date := e.ReceivedAt.Format("2006-01-02")
		if byDate[date] == nil {
			byDate[date] = &DailyUsagePoint{Date: date}
		}
		byDate[date].CostUSD += e.CostUSD
		byDate[date].Tokens += e.InputTokens + e.OutputTokens
	}

	for _, p := range byDate {
		summary.DailySeries = append(summary.DailySeries, *p)
	}
	sortDailySeries(summary.DailySeries)

	return summary, nil
}

func sortDailySeries(points []DailyUsagePoint) {
	for i := 1; i < len(points); i++ {
		for j := i; j > 0 && points[j].Date < points[j-1].Date; j-- {
			points[j], points[j-1] = points[j-1], points[j]
		}
	}
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpIngestUsage(w http.ResponseWriter, r *http.Request) {
	var body IngestUsageInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	event, err := s.IngestUsage(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, event)
}

func (s *Server) httpListUsage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	result, err := s.ListUsage(usageFiltersFromQuery(q), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpSummary(w http.ResponseWriter, r *http.Request) {
	result, err := s.Summary(usageFiltersFromQuery(r.URL.Query()))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
