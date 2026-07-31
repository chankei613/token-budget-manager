# 開発計画

| Phase | 内容 |
|---|---|
| Phase 0 | プロジェクト立ち上げ |
| Phase 1 | データモデル・Ingestion/CRUD API |
| Phase 2 | アラート・SSE |
| Phase 3 | Wails + Vue3 UI（Help含む） |
| Phase 4 | 仕上げ・署名・配布・LP |

集計の正しさが最重要のため、Ingestion→集計→アラートまでをUIより先に固め、
smoketestで検証してからUIに進む（execution-ledger/ai-decision-reviewerと同じ優先順位判断）。
