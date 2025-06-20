# TinyServer 開発タスク

このファイルは、TinyServerプロジェクトの実装タスクをフェーズ別に整理したものです。AI実装者は各タスクを順番に実装してください。

## 前提条件

- **必須**: `Spec.md`を必ず読み込んでから実装を開始すること
- Go標準ライブラリのみを使用（`net/http`パッケージは可能な限り避ける）
- テスト駆動開発（TDD）を実践
- SOLID原則とクリーンコード原則を適用
- 各フェーズで動作するデモを必ず作成

## プロジェクト初期セットアップ

### SETUP-1: プロジェクト基盤の構築
- [ ] `go.mod`ファイルの作成
- [ ] `.golangci.yml`設定ファイルの作成
- [ ] `Makefile`の作成（ビルド・テスト・リント・デモ実行コマンド）
- [ ] ディレクトリ構造の作成（pkg/, internal/, cmd/, demo/, examples/, books/, scripts/）
- [ ] 基本的な開発支援スクリプトの作成（scripts/test.sh, build.sh, lint.sh）

### SETUP-2: 共通基盤の実装
- [ ] `internal/common/errors.go` - カスタムエラー型の定義
- [ ] `internal/common/logger.go` - シンプルなログ機能
- [ ] `internal/common/constants.go` - 共通定数の定義
- [ ] 基本的なテストヘルパー関数の実装

## フェーズ1: TCP基礎実装

### PHASE1-1: TCP接続インターフェースの設計
- [ ] `pkg/tcp/interfaces.go` - TCP接続の抽象インターフェース定義
- [ ] `pkg/tcp/constants.go` - TCP関連定数の定義
- [ ] `internal/tcp/constants.go` - 内部TCP定数

### PHASE1-2: TCP接続の実装
- [ ] `internal/tcp/connection.go` - TCP接続の基本実装
- [ ] `internal/tcp/connection_test.go` - TCP接続のテスト
- [ ] `internal/tcp/listener.go` - TCPリスナーの実装
- [ ] `internal/tcp/listener_test.go` - TCPリスナーのテスト

### PHASE1-3: TCP Echoサーバ・クライアントのデモ
- [ ] `demo/phase1-tcp-echo/server/main.go` - TCPエコーサーバの実装
- [ ] `demo/phase1-tcp-echo/client/main.go` - TCPエコークライアントの実装
- [ ] `demo/phase1-tcp-echo/README.md` - デモ実行方法の説明
- [ ] `scripts/demo/run-phase1.sh` - 自動デモ実行スクリプト

### PHASE1-4: フェーズ1統合テスト
- [ ] エコーサーバ・クライアント間の通信テスト
- [ ] 複数クライアント同時接続テスト
- [ ] エラーハンドリングテスト

### PHASE1-5: フェーズ1リファクタリング
- [ ] マジック定数の削除（ポート番号、バッファサイズなど）
- [ ] その場しのぎのコードの修正
- [ ] 関数・変数名の改善
- [ ] 重複コードの削除
- [ ] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE1-6: フェーズ1コミット
- [ ] フェーズ1の全変更をコミット
- [ ] 適切なコミットメッセージの作成
- [ ] git statusで確認し、不要ファイルが含まれていないことを確認

## フェーズ2: HTTPプロトコル実装

### PHASE2-1: HTTPインターフェースの設計
- [ ] `pkg/http/interfaces.go` - HTTP抽象インターフェース定義
- [ ] `pkg/http/constants.go` - HTTP定数（ステータスコード、ヘッダーなど）
- [ ] `internal/http/constants.go` - 内部HTTP定数

### PHASE2-2: HTTPリクエスト・レスポンス構造体
- [ ] `pkg/http/request.go` - HTTPリクエスト型の定義
- [ ] `pkg/http/response.go` - HTTPレスポンス型の定義
- [ ] `internal/http/request.go` - HTTPリクエストの内部実装
- [ ] `internal/http/response.go` - HTTPレスポンスの内部実装

### PHASE2-3: HTTPパーサーの実装
- [ ] `internal/http/parser.go` - HTTPリクエスト・レスポンスパーサー
- [ ] `internal/http/parser_test.go` - パーサーのテスト
- [ ] `internal/http/request_test.go` - リクエスト処理のテスト
- [ ] `internal/http/response_test.go` - レスポンス生成のテスト

### PHASE2-4: HTTPパーサーデモ
- [ ] `demo/phase2-http-parser/main.go` - HTTPパーサー・アナライザーデモ
- [ ] `demo/phase2-http-parser/README.md` - デモ実行方法の説明
- [ ] テスト用のHTTPリクエストサンプルファイル
- [ ] `scripts/demo/run-phase2.sh` - 自動デモ実行スクリプト

### PHASE2-5: フェーズ2リファクタリング
- [ ] マジック定数の削除（HTTPステータスコード、ヘッダー名など）
- [ ] その場しのぎのコードの修正
- [ ] パーサーロジックの改善
- [ ] エラーハンドリングの統一
- [ ] コードフォーマットの統一（gofmt, golint, go vet）

### PHASE2-6: フェーズ2コミット
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

### PHASE3-6: フェーズ3コミット
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

### PHASE4-5: フェーズ4コミット
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

### PHASE5-5: フェーズ5コミット
- [ ] フェーズ5の全変更をコミット
- [ ] 適切なコミットメッセージの作成
- [ ] git statusで確認し、不要ファイルが含まれていないことを確認

## 実行可能デモとコマンドライン
- [ ] `cmd/server/main.go` - メインサーバーアプリケーション
- [ ] `cmd/client/main.go` - メインクライアントアプリケーション

## 教育用ガイドブック（オプション）
- [ ] `books/01-tcp-basics.md` - TCP基礎のガイド
- [ ] `books/02-socket-programming.md` - ソケットプログラミングガイド
- [ ] `books/03-http-protocol.md` - HTTPプロトコルガイド
- [ ] `books/04-simple-server.md` - シンプルサーバーガイド
- [ ] `books/05-client-implementation.md` - クライアント実装ガイド
- [ ] `books/06-full-integration.md` - 統合実装ガイド
- [ ] `books/07-refactoring-best-practices.md` - リファクタリングガイド

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

### 禁止事項
1. `net/http`パッケージの使用（低レベル実装の学習目的のため）
2. 外部ライブラリの使用（Go標準ライブラリのみ）
3. マジック定数の使用（全て定数として定義）
4. 複雑な抽象化（教育目的でシンプルさを重視）

このタスクリストに従って順次実装し、各フェーズで動作確認を行ってください。