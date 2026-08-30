# Token Budget Manager

「AIコストの可視化」ツール — comet-taskAI ロードマップ Product F。

プロジェクト・エージェントごとにトークン/コストの予算を設定し、リアルタイムで使用量を監視、
閾値（80%=warning, 100%=critical）を新しく超えた瞬間にアラートする。週次・月次レポートは
同じ集計APIの期間指定を変えるだけで得られる。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: v0.1.0 リリース済み

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・Ingestion/CRUD API
- [x] Phase 2: アラート・SSE
- [x] Phase 3: Wails + Vue3 UI（ダッシュボード・予算・使用量・価格設定・Help・設定）
- [x] Phase 4: 仕上げ・署名・配布・LP

macOSアプリ（署名・公証済み、Apple Silicon / Intel 共通のUniversalバイナリ）は
[GitHub Releases](https://github.com/chankei613/token-budget-manager/releases) から、
ランディングページは https://token-budget-manager-psi.vercel.app/ から入手できる。
アプリ内のHelpタブに使い方の説明がある。

## 使い方（デスクトップアプリ）

1. [Releases](../../releases) から自分のOS用のビルドをダウンロードして起動する
2. Pricing画面でモデル単価を実際の契約に合わせて調整する（未登録モデルはコスト$0のまま記録される）
3. Budgets画面で予算を作成（Source/Scope keyを空欄にすると全体対象、指定するとそのプロジェクト/システムに絞れる）
4. AIシステム側にSettings画面のAPIエンドポイント・APIキーを設定し、実行のたびに使用量をPOSTしてもらう
5. Dashboardで進捗バーを確認。80%で黄色、100%で赤になる

## 使い方（開発・ヘッドレスサーバー）

```bash
go mod tidy
make run      # :8424 でAPIサーバー起動（cmd/tbmserve）
make ui       # frontend/ の vite dev サーバー起動
make smoke    # 通しスモークテスト
```

デスクトップアプリとしてビルドするには `wails build`。

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
