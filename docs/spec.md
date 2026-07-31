# Token Budget Manager — 仕様書

> 作成: 2026-07-31

## 1. 製品概要

**「AIコストの可視化」ツール**。任意のAIシステムが実行のたびに消費したトークンを記録し、
プロジェクト・エージェント単位の予算に対する使用量をリアルタイムで監視、閾値を超えたらアラートする。
週次・月次のレポートも同じ集計APIで賄う（fromとtoの範囲を変えるだけ）。

## 2. コアコンセプト

- **UsageEvent**: 1回の実行の消費量記録（追記専用）。`POST /api/v1/usage` で送る
- **ModelPricing**: モデルごとの$/1Mトークン単価。Ingestion時にCostUSDを計算するために使う
- **Budget**: Source×ScopeKeyに対する期間（daily/weekly/monthly）ごとの上限。使用量は都度SUMで算出する
- **アラート**: 新しく閾値（80%=warning, 100%=critical）を超えた瞬間だけSSEで通知する。二重通知はしない

## 3. 機能一覧

### Phase 1 (Ingestion + CRUD)
- UsageEvent追加・一覧・集計API
- 価格表の取得・更新
- Budget CRUD・status取得

### Phase 2 (アラート)
- 閾値超過の一度きり通知（SSE）

### Phase 3 (UI)
- ダッシュボード・予算管理・使用量ビュー・価格設定・Help

## 4. データストア

```sql
usage_events (id, received_at, source, agent_id, scope_key, provider, model_id,
               input_tokens, output_tokens, cost_usd)
model_pricing (model_id, input_price_per_1m, output_price_per_1m)
budgets (id, name, source, scope_key, period, limit_tokens, limit_usd,
         last_alert_level, alert_period_start)
agent_keys (id, name, api_key_hash, created_at, revoked_at)
```
