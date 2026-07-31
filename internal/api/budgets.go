package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/token-budget-manager/internal/db"
	"github.com/chankei613/token-budget-manager/internal/events"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var errBudgetNotFound = &apiError{"budget not found"}

type CreateBudgetInput struct {
	Name        string          `json:"name"`
	Source      string          `json:"source"`
	ScopeKey    string          `json:"scope_key"`
	Period      db.BudgetPeriod `json:"period"`
	LimitTokens int64           `json:"limit_tokens"`
	LimitUSD    float64         `json:"limit_usd"`
}

func validPeriod(p db.BudgetPeriod) bool {
	switch p {
	case db.PeriodDaily, db.PeriodWeekly, db.PeriodMonthly:
		return true
	default:
		return false
	}
}

func (s *Server) CreateBudget(in CreateBudgetInput) (db.Budget, error) {
	if in.Name == "" {
		return db.Budget{}, errNameRequired
	}
	if !validPeriod(in.Period) {
		return db.Budget{}, &apiError{"period must be one of: daily, weekly, monthly"}
	}

	now := time.Now()
	b := db.Budget{
		ID:               uuid.NewString(),
		Name:             in.Name,
		Source:           in.Source,
		ScopeKey:         in.ScopeKey,
		Period:           in.Period,
		LimitTokens:      in.LimitTokens,
		LimitUSD:         in.LimitUSD,
		AlertPeriodStart: periodStart(in.Period, now),
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if err := s.DB.Create(&b).Error; err != nil {
		return db.Budget{}, err
	}
	return b, nil
}

func (s *Server) ListBudgets() ([]db.Budget, error) {
	var rows []db.Budget
	err := s.DB.Order("created_at desc").Find(&rows).Error
	return rows, err
}

func (s *Server) GetBudget(id string) (db.Budget, error) {
	var b db.Budget
	if err := s.DB.First(&b, "id = ?", id).Error; err != nil {
		return db.Budget{}, errBudgetNotFound
	}
	return b, nil
}

type UpdateBudgetInput struct {
	Name        *string  `json:"name"`
	LimitTokens *int64   `json:"limit_tokens"`
	LimitUSD    *float64 `json:"limit_usd"`
}

func (s *Server) UpdateBudget(id string, in UpdateBudgetInput) (db.Budget, error) {
	b, err := s.GetBudget(id)
	if err != nil {
		return db.Budget{}, err
	}
	if in.Name != nil {
		b.Name = *in.Name
	}
	if in.LimitTokens != nil {
		b.LimitTokens = *in.LimitTokens
	}
	if in.LimitUSD != nil {
		b.LimitUSD = *in.LimitUSD
	}
	b.UpdatedAt = time.Now()
	if err := s.DB.Save(&b).Error; err != nil {
		return db.Budget{}, err
	}
	return b, nil
}

func (s *Server) DeleteBudget(id string) error {
	res := s.DB.Delete(&db.Budget{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errBudgetNotFound
	}
	return nil
}

type BudgetStatus struct {
	Budget      db.Budget     `json:"budget"`
	TokensUsed  int64         `json:"tokens_used"`
	CostUsedUSD float64       `json:"cost_used_usd"`
	PercentUsed float64       `json:"percent_used"`
	AlertLevel  db.AlertLevel `json:"alert_level"`
	PeriodStart time.Time     `json:"period_start"`
}

// periodStart はUTC基準で期間の開始時刻を返す。
func periodStart(period db.BudgetPeriod, now time.Time) time.Time {
	now = now.UTC()
	switch period {
	case db.PeriodDaily:
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	case db.PeriodWeekly:
		// ISO週の月曜0時を週の開始とする
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return d.AddDate(0, 0, -(weekday - 1))
	case db.PeriodMonthly:
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	default:
		return now
	}
}

func (s *Server) usageWithinBudget(b db.Budget, start time.Time) (tokens int64, costUSD float64, err error) {
	q := s.DB.Model(&db.UsageEvent{}).Where("received_at >= ?", start)
	if b.Source != "" {
		q = q.Where("source = ?", b.Source)
	}
	if b.ScopeKey != "" {
		q = q.Where("scope_key = ?", b.ScopeKey)
	}

	var events []db.UsageEvent
	if err := q.Find(&events).Error; err != nil {
		return 0, 0, err
	}
	for _, e := range events {
		tokens += e.InputTokens + e.OutputTokens
		costUSD += e.CostUSD
	}
	return tokens, costUSD, nil
}

func alertLevelFor(percentUsed float64) db.AlertLevel {
	switch {
	case percentUsed >= 1.0:
		return db.AlertCritical
	case percentUsed >= 0.8:
		return db.AlertWarning
	default:
		return db.AlertNone
	}
}

func alertRank(l db.AlertLevel) int {
	switch l {
	case db.AlertCritical:
		return 2
	case db.AlertWarning:
		return 1
	default:
		return 0
	}
}

func (s *Server) BudgetStatus(id string) (BudgetStatus, error) {
	b, err := s.GetBudget(id)
	if err != nil {
		return BudgetStatus{}, err
	}
	return s.computeStatus(b)
}

func (s *Server) computeStatus(b db.Budget) (BudgetStatus, error) {
	start := periodStart(b.Period, time.Now())
	tokens, cost, err := s.usageWithinBudget(b, start)
	if err != nil {
		return BudgetStatus{}, err
	}

	percent := 0.0
	if b.LimitTokens > 0 {
		percent = max(percent, float64(tokens)/float64(b.LimitTokens))
	}
	if b.LimitUSD > 0 {
		percent = max(percent, cost/b.LimitUSD)
	}

	return BudgetStatus{
		Budget:      b,
		TokensUsed:  tokens,
		CostUsedUSD: cost,
		PercentUsed: percent,
		AlertLevel:  alertLevelFor(percent),
		PeriodStart: start,
	}, nil
}

// evaluateBudgetsForEvent は新しいUsageEventの影響を受けるBudgetを再評価し、
// 新しく閾値を超えた（今期でまだ出していないレベルの）場合だけアラートを発火する。
// 期間が変わっていたらLastAlertLevelを先にリセットする。永続化エラーはログせず無視する
// （アラート評価はIngestion自体を失敗させない付加機能のため）。
func (s *Server) evaluateBudgetsForEvent(event db.UsageEvent) {
	var budgets []db.Budget
	q := s.DB.Where("source = ? OR source = ''", event.Source)
	q = q.Where("scope_key = ? OR scope_key = ''", event.ScopeKey)
	if err := q.Find(&budgets).Error; err != nil {
		return
	}

	for _, b := range budgets {
		start := periodStart(b.Period, time.Now())
		if !start.Equal(b.AlertPeriodStart) {
			b.AlertPeriodStart = start
			b.LastAlertLevel = db.AlertNone
		}

		status, err := s.computeStatus(b)
		if err != nil {
			continue
		}

		if alertRank(status.AlertLevel) > alertRank(b.LastAlertLevel) {
			b.LastAlertLevel = status.AlertLevel
			s.Events.Publish(events.Event{
				Type:     events.EventBudgetAlert,
				BudgetID: b.ID,
				Level:    string(status.AlertLevel),
				At:       time.Now(),
			})
		}

		_ = s.DB.Save(&b).Error
	}
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpCreateBudget(w http.ResponseWriter, r *http.Request) {
	var body CreateBudgetInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b, err := s.CreateBudget(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (s *Server) httpListBudgets(w http.ResponseWriter, r *http.Request) {
	rows, err := s.ListBudgets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, rows)
}

func (s *Server) httpGetBudget(w http.ResponseWriter, r *http.Request) {
	b, err := s.GetBudget(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) httpUpdateBudget(w http.ResponseWriter, r *http.Request) {
	var body UpdateBudgetInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b, err := s.UpdateBudget(chi.URLParam(r, "id"), body)
	if err != nil {
		if err == errBudgetNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) httpDeleteBudget(w http.ResponseWriter, r *http.Request) {
	if err := s.DeleteBudget(chi.URLParam(r, "id")); err != nil {
		if err == errBudgetNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) httpBudgetStatus(w http.ResponseWriter, r *http.Request) {
	status, err := s.BudgetStatus(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, status)
}
