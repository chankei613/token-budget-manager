# Token Budget Manager

「AIコストの可視化」ツール — comet-taskAI ロードマップ Product F。

プロジェクト・エージェントごとにトークン/コストの予算を設定し、リアルタイムで使用量を監視、
閾値（80%=warning, 100%=critical）を新しく超えた瞬間にアラートする。週次・月次レポートは
同じ集計APIの期間指定を変えるだけで得られる。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: Phase 1-2（Ingestion/CRUD API + アラート/SSE）完了

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・Ingestion/CRUD API
- [x] Phase 2: アラート・SSE
- [ ] Phase 3: Wails + Vue3 UI
- [ ] Phase 4: 仕上げ・署名・配布・LP

## 使い方（開発用ヘッドレスサーバー）

```bash
go mod tidy
go run .      # :8424 でAPIサーバー起動
go run ./cmd/smoketest
```

### 使用量の記録・予算作成

```bash
curl -X POST localhost:8424/api/v1/budgets \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"name":"Photo Editor月次予算","scope_key":"project:photo-editor","period":"monthly","limit_usd":50}'

curl -X POST localhost:8424/api/v1/usage \
  -H "Authorization: Bearer $API_KEY" -H "Content-Type: application/json" \
  -d '{"agent_id":"claude-01","scope_key":"project:photo-editor","model_id":"claude-sonnet-5","input_tokens":12000,"output_tokens":3000}'
```

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST/GET/DELETE | `/api/v1/keys` | APIキー管理 |
| POST | `/api/v1/usage` | 使用量記録（自動でコスト計算） |
| GET | `/api/v1/usage` | 一覧・フィルタ |
| GET | `/api/v1/usage/summary` | 集計（週次/月次レポートを兼ねる） |
| GET/PUT | `/api/v1/pricing` | モデル価格の取得・更新 |
| POST/GET/PATCH/DELETE | `/api/v1/budgets` | 予算のCRUD |
| GET | `/api/v1/budgets/{id}/status` | 今期の使用量・percent_used・alert_level |
| GET | `/api/v1/events` | SSE（budget:alert） |

## ディレクトリ構成

```
internal/db/       GORMモデル・デフォルト価格シード
internal/api/       REST API（usage/pricing/budgets/events）+ 認証ミドルウェア
internal/events/    SSE配信用pub/sub
cmd/smoketest/      通しスモークテスト
docs/                設計ドキュメント
```
