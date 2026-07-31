// cmd/smoketest はToken Budget ManagerのAPIを一時DBで自前起動し、
// ブートストラップ鍵発行 → 価格確認 → Budget作成 → 使用量投入 → 集計 →
// アラート発火（一度きり） → SSE配信、の一連が通しで動くことを確認する。
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/chankei613/token-budget-manager/internal/api"
	"github.com/chankei613/token-budget-manager/internal/db"
)

func main() {
	dbPath := "smoketest.db"
	_ = os.Remove(dbPath)
	defer func() { _ = os.Remove(dbPath) }()

	conn, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("db init: %v", err)
	}

	srv := httptest.NewServer(api.NewRouter(conn))
	defer srv.Close()

	// 1. bootstrap key issuance
	issueBody, _ := json.Marshal(map[string]string{"name": "smoketest"})
	resp, err := http.Post(srv.URL+"/api/v1/keys", "application/json", bytes.NewReader(issueBody))
	if err != nil {
		log.Fatal(err)
	}
	var issued api.IssueKeyResult
	if err := json.NewDecoder(resp.Body).Decode(&issued); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if issued.APIKey == "" {
		log.Fatal("FAIL: bootstrap key issuance returned empty key")
	}
	fmt.Println("PASS: bootstrap key issued")

	authed := func(method, path string, body []byte) *http.Response {
		req, _ := http.NewRequest(method, srv.URL+path, bytes.NewReader(body))
		req.Header.Set("Authorization", "Bearer "+issued.APIKey)
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		return resp
	}

	// 2. default pricing is seeded
	resp = authed(http.MethodGet, "/api/v1/pricing", nil)
	var pricing []db.ModelPricing
	if err := json.NewDecoder(resp.Body).Decode(&pricing); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if len(pricing) == 0 {
		log.Fatal("FAIL: expected default pricing to be seeded")
	}
	fmt.Println("PASS: default pricing seeded")

	// 3. start an SSE subscriber before crossing a threshold
	sseReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/api/v1/events", nil)
	sseReq.Header.Set("Authorization", "Bearer "+issued.APIKey)
	sseResp, err := http.DefaultClient.Do(sseReq)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = sseResp.Body.Close() }()
	sseLines := make(chan string, 16)
	go func() {
		scanner := bufio.NewScanner(sseResp.Body)
		for scanner.Scan() {
			sseLines <- scanner.Text()
		}
	}()

	// 4. create a small budget so a single usage event crosses "critical"
	createBudgetBody, _ := json.Marshal(api.CreateBudgetInput{
		Name:        "smoketest budget",
		Source:      "smoketest",
		ScopeKey:    "project:demo",
		Period:      db.PeriodDaily,
		LimitTokens: 1000,
	})
	resp = authed(http.MethodPost, "/api/v1/budgets", createBudgetBody)
	var budget db.Budget
	if err := json.NewDecoder(resp.Body).Decode(&budget); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if budget.ID == "" {
		log.Fatal("FAIL: budget creation returned empty ID")
	}
	fmt.Println("PASS: budget created")

	// 5. ingest usage that exceeds the budget (2000 tokens > 1000 limit)
	usageBody, _ := json.Marshal(api.IngestUsageInput{
		Source:       "smoketest",
		AgentID:      "claude-01",
		ScopeKey:     "project:demo",
		Provider:     "anthropic",
		ModelID:      "claude-sonnet-5",
		InputTokens:  1500,
		OutputTokens: 500,
	})
	resp = authed(http.MethodPost, "/api/v1/usage", usageBody)
	var event db.UsageEvent
	if err := json.NewDecoder(resp.Body).Decode(&event); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if event.CostUSD <= 0 {
		log.Fatalf("FAIL: expected non-zero cost from seeded pricing, got %+v", event)
	}
	fmt.Println("PASS: usage ingested with computed cost")

	// 6. budget status reflects the overage
	resp = authed(http.MethodGet, "/api/v1/budgets/"+budget.ID+"/status", nil)
	var status api.BudgetStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if status.AlertLevel != db.AlertCritical {
		log.Fatalf("FAIL: expected critical alert level, got %+v", status)
	}
	fmt.Println("PASS: budget status shows critical overage")

	// 7. confirm the SSE subscriber saw exactly one budget:alert (no duplicate on a second ingest)
	usageBody2, _ := json.Marshal(api.IngestUsageInput{
		Source: "smoketest", AgentID: "claude-01", ScopeKey: "project:demo",
		ModelID: "claude-sonnet-5", InputTokens: 100, OutputTokens: 10,
	})
	resp = authed(http.MethodPost, "/api/v1/usage", usageBody2)
	_ = resp.Body.Close()

	alertCount := 0
	deadline := time.After(2 * time.Second)
loop:
	for {
		select {
		case line := <-sseLines:
			// SSEの "event: budget:alert" 行だけを数える（"data: ..." 行にも
			// JSONのtypeフィールドとして同じ文字列が含まれるため二重カウントしない）
			if strings.HasPrefix(line, "event: budget:alert") {
				alertCount++
			}
		case <-deadline:
			break loop
		}
	}
	if alertCount != 1 {
		log.Fatalf("FAIL: expected exactly 1 budget:alert (no duplicate), got %d", alertCount)
	}
	fmt.Println("PASS: exactly one alert fired despite two threshold-crossing events")

	// 8. summary reflects both events
	resp = authed(http.MethodGet, "/api/v1/usage/summary?source=smoketest", nil)
	var summary api.UsageSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		log.Fatal(err)
	}
	_ = resp.Body.Close()
	if summary.TotalInputTokens != 1600 {
		log.Fatalf("FAIL: expected 1600 total input tokens, got %d", summary.TotalInputTokens)
	}
	fmt.Println("PASS: summary aggregates both usage events")

	fmt.Println("SMOKE TEST OK")
}
