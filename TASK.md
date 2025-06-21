# TinyServer 開発タスク

このファイルは、TinyServerプロジェクトの実装タスクをフェーズ別に整理したものです。AI実装者は各タスクを順番に実装してください。

## 📊 プロジェクト進捗状況 (2025/06/21更新)

### 完了したフェーズ
- ✅ **プロジェクト初期セットアップ** - プロジェクト基盤と共通ライブラリ
- ✅ **フェーズ1: TCP基礎実装** - TCP接続、リスナー、エコーサーバ・クライアントデモ

### 現在の状況
- **実装済み**: TCP基盤、エコーサーバ・クライアント、HTTPプロトコル基盤、HTTPパーサーデモ
- **動作確認**: `./scripts/demo/run-phase1.sh` でTCP Echoデモが動作、`./scripts/demo/run-phase2.sh` でHTTPパーサーデモが動作
- **次のステップ**: フェーズ2 HTTPプロトコル実装完了（リファクタリング・ガイドブック）

### 実装ファイル一覧
```
pkg/tcp/          - TCP公開インターフェース
pkg/http/         - HTTP公開インターフェース
internal/tcp/     - TCP実装詳細
internal/http/    - HTTP実装詳細（パーサー、リクエスト・レスポンス処理）
internal/common/  - 共通基盤（エラー、ログ、定数）
demo/phase1-tcp-echo/   - TCPエコーデモ
demo/phase2-http-parser/ - HTTPパーサーデモ
scripts/          - 開発支援スクリプト
```

## 前提条件

- **必須**: `Spec.md`を必ず読み込んでから実装を開始すること
- Go標準ライブラリのみを使用（`net/http`パッケージは可能な限り避ける）
- テスト駆動開発（TDD）を実践
- SOLID原則とクリーンコード原則を適用
- 各フェーズで動作するデモを必ず作成

## プロジェクト初期セットアップ

### SETUP-1: プロジェクト基盤の構築 ✅ COMPLETED
- [x] `go.mod`ファイルの作成
- [x] `.golangci.yml`設定ファイルの作成
- [x] `Makefile`の作成（ビルド・テスト・リント・デモ実行コマンド）
- [x] ディレクトリ構造の作成（pkg/, internal/, cmd/, demo/, examples/, books/, scripts/）
- [x] 基本的な開発支援スクリプトの作成（scripts/test.sh, build.sh, lint.sh）

### SETUP-2: 共通基盤の実装 ✅ COMPLETED
- [x] `internal/common/errors.go` - カスタムエラー型の定義
- [x] `internal/common/logger.go` - シンプルなログ機能
- [x] `internal/common/constants.go` - 共通定数の定義
- [x] `internal/common/testing_helpers.go` - 基本的なテストヘルパー関数の実装
- [x] `internal/common/errors_test.go` - エラー型のテスト

## フェーズ1: TCP基礎実装

### PHASE1-1: TCP接続インターフェースの設計 ✅ COMPLETED
- [x] `pkg/tcp/interfaces.go` - TCP接続の抽象インターフェース定義
- [x] `pkg/tcp/constants.go` - TCP関連定数の定義
- [x] `internal/tcp/constants.go` - 内部TCP定数

### PHASE1-2: TCP接続の実装 ✅ COMPLETED
- [x] `internal/tcp/connection.go` - TCP接続の基本実装
- [x] `internal/tcp/connection_test.go` - TCP接続のテスト
- [x] `internal/tcp/listener.go` - TCPリスナーの実装
- [x] `internal/tcp/listener_test.go` - TCPリスナーのテスト

### PHASE1-3: TCP Echoサーバ・クライアントのデモ ✅ COMPLETED
- [x] `demo/phase1-tcp-echo/server/main.go` - TCPエコーサーバの実装
- [x] `demo/phase1-tcp-echo/client/main.go` - TCPエコークライアントの実装
- [x] `demo/phase1-tcp-echo/README.md` - デモ実行方法の説明
- [x] `scripts/demo/run-phase1.sh` - 自動デモ実行スクリプト

### PHASE1-4: フェーズ1統合テスト ✅ COMPLETED
- [x] エコーサーバ・クライアント間の通信テスト
- [x] 複数クライアント同時接続テスト
- [x] エラーハンドリングテスト

### PHASE1-5: フェーズ1リファクタリング ✅ COMPLETED
- [x] マジック定数の削除（ポート番号、バッファサイズなど）
- [x] その場しのぎのコードの修正
- [x] 関数・変数名の改善
- [x] 重複コードの削除
- [x] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE1-6: フェーズ1教育用ガイドブック作成 ✅ COMPLETED
- [x] `books/01-tcp-basics.md` - TCP基礎の教育用ガイド作成
- [x] TCP/IPプロトコルの仕組みと動作原理の解説
- [x] ソケットプログラミングの基本概念とAPI説明
- [x] エコーサーバ・クライアントの実装ステップバイステップガイド
- [x] トラブルシューティングとよくある問題の解決方法
- [x] ガイドに従って実装すればフェーズ1を完全に再現できる詳細度で記述
- [x] 学習者が理論と実装の両方を理解できる構成

### PHASE1-7: フェーズ1コミット ✅ COMPLETED
- [x] フェーズ1の全変更をコミット
- [x] 適切なコミットメッセージの作成
- [x] git statusで確認し、不要ファイルが含まれていないことを確認

## フェーズ2: HTTPプロトコル実装

### PHASE2-1: HTTPインターフェースの設計 ✅ COMPLETED
- [x] `pkg/http/interfaces.go` - HTTP抽象インターフェース定義
- [x] `pkg/http/constants.go` - HTTP定数（ステータスコード、ヘッダーなど）
- [x] `internal/http/constants.go` - 内部HTTP定数

### PHASE2-2: HTTPリクエスト・レスポンス構造体 ✅ COMPLETED
- [x] `pkg/http/request.go` - HTTPリクエスト型の定義
- [x] `pkg/http/response.go` - HTTPレスポンス型の定義
- [x] `internal/http/request.go` - HTTPリクエストの内部実装
- [x] `internal/http/response.go` - HTTPレスポンスの内部実装

### PHASE2-3: HTTPパーサーの実装 ✅ COMPLETED
- [x] `internal/http/parser.go` - HTTPリクエスト・レスポンスパーサー
- [x] `internal/http/parser_test.go` - パーサーのテスト
- [x] `internal/http/request_test.go` - リクエスト処理のテスト
- [x] `internal/http/response_test.go` - レスポンス生成のテスト

### PHASE2-4: HTTPパーサーデモ ✅ COMPLETED
- [x] `demo/phase2-http-parser/main.go` - HTTPパーサー・アナライザーデモ
- [x] `demo/phase2-http-parser/README.md` - デモ実行方法の説明
- [x] テスト用のHTTPリクエストサンプルファイル (`demo/phase2-http-parser/samples/`)
- [x] `scripts/demo/run-phase2.sh` - 自動デモ実行スクリプト

### PHASE2-5: フェーズ2リファクタリング
- [ ] マジック定数の削除（HTTPステータスコード、ヘッダー名など）
- [ ] その場しのぎのコードの修正
- [ ] パーサーロジックの改善
- [ ] エラーハンドリングの統一
- [ ] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE2-6: フェーズ2教育用ガイドブック作成
- [ ] `books/02-http-protocol.md` - HTTPプロトコルの教育用ガイド作成
- [ ] HTTP/1.1プロトコルの仕様と構造の詳細解説
- [ ] リクエスト・レスポンスの解析手法と実装方法
- [ ] HTTPヘッダーとステータスコードの完全解説
- [ ] パーサー実装のアルゴリズムとデータ構造設計
- [ ] ガイドに従って実装すればフェーズ2を完全に再現できる詳細度で記述
- [ ] 実際のHTTPトラフィック例と解析演習を含む

### PHASE2-7: フェーズ2コミット
- [ ] フェーズ2の全変更をコミット
- [ ] 適切なコミットメッセージの作成
- [ ] git statusで確認し、不要ファイルが含まれていないことを確認

## フェーズ3: シンプルサーバ実装

### PHASE3-1: サーバーインターフェースの設計
- [ ] `pkg/tinyserver/interfaces.go` - メインサーバーインターフェース定義
- [ ] `pkg/tinyserver/server.go` - 公開サーバー型の定義
- [ ] `internal/server/constants.go` - サーバー関連定数

### PHASE3-2: HTTPサーバーの実装
- [ ] `internal/server/server.go` - HTTPサーバーの基本実装
- [ ] `internal/server/server_test.go` - サーバーのテスト
- [ ] `internal/server/handler.go` - リクエストハンドラーの実装
- [ ] `internal/server/handler_test.go` - ハンドラーのテスト

### PHASE3-3: ミドルウェア機能
- [ ] `internal/server/middleware.go` - 基本的なミドルウェア実装
  - ログミドルウェア
  - 静的ファイル配信ミドルウェア
  - エラーハンドリングミドルウェア

### PHASE3-4: 静的ファイルサーバーデモ
- [ ] `demo/phase3-simple-server/main.go` - 静的ファイルサーバー実装
- [ ] `demo/phase3-simple-server/static/index.html` - デモ用HTMLファイル
- [ ] `demo/phase3-simple-server/README.md` - デモ実行方法の説明
- [ ] `scripts/demo/run-phase3.sh` - 自動デモ実行スクリプト

### PHASE3-5: フェーズ3リファクタリング
- [ ] マジック定数の削除（ポート番号、ファイルパス、MIMEタイプなど）
- [ ] その場しのぎのコードの修正
- [ ] ミドルウェアの抽象化改善
- [ ] ハンドラーロジックの整理
- [ ] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE3-6: フェーズ3教育用ガイドブック作成
- [ ] `books/03-simple-server.md` - Webサーバー実装の教育用ガイド作成
- [ ] Webサーバーアーキテクチャと設計パターンの解説
- [ ] リクエストルーティングとハンドラーパターンの実装手法
- [ ] ミドルウェアシステムの設計と実装方法
- [ ] 静的ファイル配信とMIMEタイプ処理の詳細
- [ ] エラーハンドリングとログ機能の実装
- [ ] ガイドに従って実装すればフェーズ3を完全に再現できる詳細度で記述
- [ ] 実際のWebサーバーとの比較とベンチマーク手法

### PHASE3-7: フェーズ3コミット
- [ ] フェーズ3の全変更をコミット
- [ ] 適切なコミットメッセージの作成
- [ ] git statusで確認し、不要ファイルが含まれていないことを確認

## フェーズ4: クライアント実装

### PHASE4-1: HTTPクライアントインターフェース設計
- [ ] `pkg/tinyserver/client.go` - 公開クライアント型の定義
- [ ] HTTPクライアントのインターフェース定義

### PHASE4-2: HTTPクライアントの実装
- [ ] HTTPリクエスト送信機能の実装
- [ ] レスポンス受信・解析機能の実装
- [ ] タイムアウト処理の実装
- [ ] HTTPクライアントのテスト

### PHASE4-3: TinyCurl HTTPクライアントツールデモ
- [ ] `demo/phase4-http-client/main.go` - curlライクなHTTPクライアント実装
- [ ] 様々なHTTPメソッド対応（GET, POST, PUT, DELETE）
- [ ] レスポンス表示とパフォーマンス測定
- [ ] `demo/phase4-http-client/README.md` - デモ実行方法の説明
- [ ] `scripts/demo/run-phase4.sh` - 自動デモ実行スクリプト

### PHASE4-4: フェーズ4リファクタリング
- [ ] マジック定数の削除（タイムアウト値、リトライ回数など）
- [ ] その場しのぎのコードの修正
- [ ] クライアントロジックの整理
- [ ] エラーハンドリングの統一
- [ ] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE4-5: フェーズ4教育用ガイドブック作成
- [ ] `books/04-client-implementation.md` - HTTPクライアント実装の教育用ガイド作成
- [ ] HTTPクライアントアーキテクチャと通信フローの解説
- [ ] リクエスト送信とレスポンス受信の実装詳細
- [ ] タイムアウト処理とエラーハンドリングの設計
- [ ] 非同期通信とパフォーマンス最適化手法
- [ ] 実用的なcurlライクツールの設計と実装
- [ ] ガイドに従って実装すればフェーズ4を完全に再現できる詳細度で記述
- [ ] 実際のHTTPクライアントライブラリとの比較分析

### PHASE4-6: フェーズ4コミット
- [ ] フェーズ4の全変更をコミット
- [ ] 適切なコミットメッセージの作成
- [ ] git statusで確認し、不要ファイルが含まれていないことを確認

## フェーズ5: 統合・最適化

### PHASE5-1: 最適化機能の実装
- [ ] コネクションプールの実装
- [ ] ロングポーリング機能の実装
- [ ] JSON API処理機能の実装
- [ ] パフォーマンス監視機能

### PHASE5-2: フルスタックチャットアプリケーション
- [ ] `demo/phase5-full-stack/server/main.go` - チャットサーバー実装
- [ ] `demo/phase5-full-stack/client/main.go` - チャットクライアント実装
- [ ] `demo/phase5-full-stack/static/index.html` - チャットアプリフロントエンド
- [ ] メッセージ履歴の保存・読み込み機能
- [ ] 複数ユーザー対応

### PHASE5-3: 統合デモ
- [ ] `demo/phase5-full-stack/README.md` - フルスタックデモの説明
- [ ] `scripts/demo/run-phase5.sh` - 自動デモ実行スクリプト
- [ ] 全フェーズの統合テスト

### PHASE5-4: フェーズ5リファクタリング
- [ ] マジック定数の削除（設定値、リソース制限など）
- [ ] その場しのぎのコードの修正
- [ ] 全体アーキテクチャの最適化
- [ ] パフォーマンス改善
- [ ] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE5-5: フェーズ5教育用ガイドブック作成
- [ ] `books/05-full-integration.md` - 統合・最適化の教育用ガイド作成
- [ ] フルスタックWebアプリケーションアーキテクチャの解説
- [ ] パフォーマンス最適化技術とボトルネック分析手法
- [ ] ロングポーリングとリアルタイム通信の実装
- [ ] JSON API設計とデータ交換プロトコル
- [ ] チャットアプリケーションの設計と実装手法
- [ ] ガイドに従って実装すればフェーズ5を完全に再現できる詳細度で記述
- [ ] 実用的なWebアプリケーション開発における応用例

### PHASE5-6: 最終統合ガイドブック作成
- [ ] `books/06-refactoring-best-practices.md` - リファクタリングとベストプラクティスガイド
- [ ] SOLID原則とクリーンコード原則の実践的適用
- [ ] プロジェクト全体の設計思想と意思決定の解説
- [ ] パフォーマンス改善とスケーラビリティ向上手法
- [ ] テスト戦略とQA手法の包括的解説
- [ ] 実用的Webサーバー開発への発展と応用
- [ ] TinyServerプロジェクト完全再現ガイドとしての総仕上げ

### PHASE5-7: フェーズ5コミット
- [ ] フェーズ5の全変更をコミット
- [ ] 適切なコミットメッセージの作成
- [ ] git statusで確認し、不要ファイルが含まれていないことを確認

## 実行可能デモとコマンドライン
- [ ] `cmd/server/main.go` - メインサーバーアプリケーション
- [ ] `cmd/client/main.go` - メインクライアントアプリケーション

## 教育用ガイドブック（各フェーズで作成済み）
- [x] `books/01-tcp-basics.md` - TCP基礎のガイド（フェーズ1で作成）
- [x] `books/02-http-protocol.md` - HTTPプロトコルガイド（フェーズ2で作成）
- [x] `books/03-simple-server.md` - シンプルサーバーガイド（フェーズ3で作成）
- [x] `books/04-client-implementation.md` - クライアント実装ガイド（フェーズ4で作成）
- [x] `books/05-full-integration.md` - 統合実装ガイド（フェーズ5で作成）
- [x] `books/06-refactoring-best-practices.md` - リファクタリングガイド（フェーズ5で作成）

## 品質保証とドキュメント
- [ ] 全パッケージのテストカバレッジ80%以上
- [ ] `go vet`、`golint`、`gofmt`を通すコード品質
- [ ] パフォーマンステスト（ベンチマーク）
- [ ] 各デモの動作確認ドキュメント
- [ ] トラブルシューティングガイド

## 実装時の注意事項

### 必須要件
1. **教育的価値**: 各実装は学習目的を重視し、低レベルな理解を促進すること
2. **段階的構築**: 前フェーズの実装に依存し、段階的に機能を追加すること
3. **動作デモ**: 各フェーズで必ず動作するデモを作成し、達成感を提供すること
4. **テスト重視**: 全ての機能にテストを書き、TDDを実践すること
5. **SOLID原則**: インターフェース設計と依存性注入を適切に活用すること
6. **📊 進捗管理**: **タスクを完了するたびに、必ずTASK.mdの該当タスクを更新し、進捗状況を反映すること**

### 禁止事項
1. `net/http`パッケージの使用（低レベル実装の学習目的のため）
2. 外部ライブラリの使用（Go標準ライブラリのみ）
3. マジック定数の使用（全て定数として定義）
4. 複雑な抽象化（教育目的でシンプルさを重視）

### 📋 タスク管理の運用ルール
- **必須**: 各タスク完了時にTASK.mdの進捗欄（[ ] → [x]）を更新すること
- **必須**: 各フェーズ完了時に「📊 プロジェクト進捗状況」セクションを最新化すること
- **推奨**: 次に取り組むタスクを明確にし、優先順位を整理すること

このタスクリストに従って順次実装し、各フェーズで動作確認を行ってください。