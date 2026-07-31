package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankei613/token-budget-manager/internal/api"
	"github.com/chankei613/token-budget-manager/internal/db"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// apiAddr はUsage Ingestion/CRUD APIの待ち受けアドレス。外部プロセスがアプリ起動中
// いつでもPOSTできるよう、ウインドウの表示/非表示に関わらずこのHTTPサーバーは動き続ける。
const apiAddr = "127.0.0.1:8424"

// App はWailsのバインディング。実処理は internal/api.Server が持っている。
type App struct {
	ctx    context.Context
	server *api.Server
	srv    *http.Server
	ready  bool
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir := appDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		runtime.LogErrorf(ctx, "data dir error: %s", err)
		return
	}

	conn, err := db.Init(filepath.Join(dataDir, "token-budget-manager.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "db init error: %s", err)
		return
	}
	a.server = api.New(conn)

	a.srv = &http.Server{Addr: apiAddr, Handler: a.server.Router()}
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			runtime.LogErrorf(ctx, "api server error: %s", err)
		}
	}()

	// SSEブローカーの内容をWailsのネイティブイベントとして転送する。
	stream, _ := a.server.Events.Subscribe()
	go func() {
		for ev := range stream {
			runtime.EventsEmit(ctx, ev.Type, ev)
		}
	}()

	a.ready = true
	runtime.LogInfof(ctx, "Token Budget Manager ready (api: http://%s, data: %s)", apiAddr, dataDir)
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		_ = a.srv.Close()
	}
	if a.server != nil {
		if sqlDB, err := a.server.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

var errNotReady = &notReadyError{}

type notReadyError struct{}

func (e *notReadyError) Error() string { return "app not ready — check startup logs" }

// ─── フロントエンドへ公開するメソッド ──────────────────────────────────────────

func (a *App) GetAppVersion() string {
	return AppVersion
}

func (a *App) GetAPIURL() string {
	return "http://" + apiAddr
}

func (a *App) ListBudgets() ([]db.Budget, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListBudgets()
}

func (a *App) CreateBudget(name, source, scopeKey string, period db.BudgetPeriod, limitTokens int64, limitUSD float64) (db.Budget, error) {
	if !a.ready {
		return db.Budget{}, errNotReady
	}
	return a.server.CreateBudget(api.CreateBudgetInput{
		Name: name, Source: source, ScopeKey: scopeKey, Period: period,
		LimitTokens: limitTokens, LimitUSD: limitUSD,
	})
}

func (a *App) UpdateBudget(id string, name *string, limitTokens *int64, limitUSD *float64) (db.Budget, error) {
	if !a.ready {
		return db.Budget{}, errNotReady
	}
	return a.server.UpdateBudget(id, api.UpdateBudgetInput{Name: name, LimitTokens: limitTokens, LimitUSD: limitUSD})
}

func (a *App) DeleteBudget(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.DeleteBudget(id)
}

func (a *App) BudgetStatus(id string) (api.BudgetStatus, error) {
	if !a.ready {
		return api.BudgetStatus{}, errNotReady
	}
	return a.server.BudgetStatus(id)
}

func (a *App) ListUsage(source, agentID, scopeKey, modelID string, limit, offset int) (api.ListUsageResult, error) {
	if !a.ready {
		return api.ListUsageResult{}, errNotReady
	}
	f := api.ListUsageFilters{Source: source, AgentID: agentID, ScopeKey: scopeKey, ModelID: modelID}
	return a.server.ListUsage(f, limit, offset)
}

func (a *App) Summary(source, agentID, scopeKey, modelID string) (api.UsageSummary, error) {
	if !a.ready {
		return api.UsageSummary{}, errNotReady
	}
	f := api.ListUsageFilters{Source: source, AgentID: agentID, ScopeKey: scopeKey, ModelID: modelID}
	return a.server.Summary(f)
}

func (a *App) ListPricing() ([]db.ModelPricing, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListPricing()
}

func (a *App) SetPricing(modelID string, inputPricePer1M, outputPricePer1M float64) (db.ModelPricing, error) {
	if !a.ready {
		return db.ModelPricing{}, errNotReady
	}
	return a.server.SetPricing(api.SetPricingInput{
		ModelID: modelID, InputPricePer1M: inputPricePer1M, OutputPricePer1M: outputPricePer1M,
	})
}

func (a *App) ListKeys() ([]db.AgentKey, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListKeys()
}

func (a *App) IssueKey(name string) (api.IssueKeyResult, error) {
	if !a.ready {
		return api.IssueKeyResult{}, errNotReady
	}
	return a.server.IssueKey(name)
}

func (a *App) RevokeKey(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.RevokeKey(id)
}

// Quit はアプリを完全終了する（Settings 画面から呼ぶ）。
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".token-budget-manager")
}
