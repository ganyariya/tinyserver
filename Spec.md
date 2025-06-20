# TinyServer プロジェクト仕様書 v1.4

## プロジェクト概要

**TinyServer**は、TCP/IP、HTTP、Webサーバ・クライアントの仕組みを実装しながら学ぶ教育的プロジェクトです。小さく段階的な実装を通じて、ネットワークプログラミングの基礎を深く理解することを目的としています。

## プロジェクト目標

- TCP/IPプロトコルの基本的な仕組みを自作実装で理解する
- HTTPプロトコルの構造と動作を実装を通じて学ぶ
- WebサーバとクライアントがTCP/IP + HTTPでどう通信するかを体験的に学習する
- テスト駆動開発の実践を通じて、堅牢なコードの書き方を身につける
- **Goのベストプラクティスに従った、実用的でクリーンなコード設計を学ぶ**
- **SOLID原則に基づいた抽象化と拡張性の高いアーキテクチャを実践する**
- **各段階で実際に動くものを作り、達成感と学習の楽しさを体験する**

## 技術仕様

### 開発言語・環境
- **言語**: Go
- **依存関係**: Go標準ライブラリのみ使用
- **制約**: `net/http`パッケージは可能な限り使用せず、低レベルな実装を行う
- **開発手法**: テスト駆動開発（TDD）
- **テストフレームワーク**: Go標準の`testing`パッケージ

### コード品質基準

#### Goベストプラクティス準拠
- **命名規則**: Go標準の命名規則に厳密に従う
- **パッケージ設計**: 単一責任の原則に基づいたパッケージ分割
- **エラーハンドリング**: Goのエラーハンドリング慣例に従った実装
- **コードフォーマット**: `gofmt`, `golint`, `go vet`を通すコード

#### SOLID原則の適用
- **単一責任原則 (SRP)**: 各型・関数は一つの責任のみを持つ
- **開放閉鎖原則 (OCP)**: インターフェースによる拡張性の確保
- **リスコフ置換原則 (LSP)**: 適切な継承関係の設計
- **インターフェース分離原則 (ISP)**: 小さく特化したインターフェース設計
- **依存性逆転原則 (DIP)**: 具象ではなく抽象に依存する設計

#### クリーンコード原則
- **マジック定数の排除**: 全ての定数は適切に命名し、定数宣言で管理
- **DRY原則**: 重複コードの排除と適切な抽象化
- **可読性**: 自己文書化されたコード
- **拡張性**: 将来の機能追加に対応できる設計

#### コメント方針
**「Code tells you how, Comments tell you why」の原則**

- **基本方針**: コメントは最小限に留める
- **自己文書化**: 型名、関数名、変数名で意図を明確に表現
- **適切な抽象化**: 型システムを活用してコードの意図を伝える

### プロジェクト構成

```
tinyserver/
├── README.md
├── go.mod
├── go.sum
├── Makefile                  # ビルド・テスト・リント自動化
├── .golangci.yml            # linter設定
├── books/                    # 教育用ガイド
│   ├── 01-tcp-basics.md
│   ├── 02-socket-programming.md
│   ├── 03-http-protocol.md
│   ├── 04-simple-server.md
│   ├── 05-client-implementation.md
│   ├── 06-full-integration.md
│   └── 07-refactoring-best-practices.md
├── pkg/                      # 公開API
│   ├── tinyserver/
│   │   ├── interfaces.go    # 主要インターフェース定義
│   │   ├── server.go
│   │   └── client.go
│   ├── http/
│   │   ├── interfaces.go
│   │   ├── request.go
│   │   ├── response.go
│   │   └── constants.go     # HTTP定数定義
│   └── tcp/
│       ├── interfaces.go
│       ├── connection.go
│       └── constants.go     # TCP定数定義
├── internal/
│   ├── tcp/                  # TCP実装詳細
│   │   ├── connection.go
│   │   ├── connection_test.go
│   │   ├── listener.go
│   │   ├── listener_test.go
│   │   └── constants.go
│   ├── http/                 # HTTP実装詳細
│   │   ├── request.go
│   │   ├── request_test.go
│   │   ├── response.go
│   │   ├── response_test.go
│   │   ├── parser.go
│   │   ├── parser_test.go
│   │   └── constants.go
│   ├── server/               # サーバ実装詳細
│   │   ├── server.go
│   │   ├── server_test.go
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── middleware.go
│   │   └── constants.go
│   └── common/               # 共通ユーティリティ
│       ├── errors.go
│       ├── logger.go
│       └── constants.go
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── demo/                     # 各フェーズの実行可能デモ
│   ├── phase1-tcp-echo/      # TCP Echo サーバ・クライアント
│   │   ├── server/
│   │   │   └── main.go
│   │   ├── client/
│   │   │   └── main.go
│   │   └── README.md
│   ├── phase2-http-parser/   # HTTP パース デモ
│   │   ├── main.go
│   │   └── README.md
│   ├── phase3-simple-server/ # シンプル Web サーバ
│   │   ├── main.go
│   │   ├── static/
│   │   │   └── index.html
│   │   └── README.md
│   ├── phase4-http-client/   # HTTP クライアント
│   │   ├── main.go
│   │   └── README.md
│   └── phase5-full-stack/    # フルスタック Web アプリ
│       ├── server/
│       │   └── main.go
│       ├── client/
│       │   └── main.go
│       ├── static/
│       │   └── index.html
│       └── README.md
├── examples/                 # 実行例・デモ
│   ├── simple-get/
│   ├── file-server/
│   └── echo-server/
└── scripts/                  # 開発支援スクリプト
    ├── test.sh
    ├── build.sh
    ├── lint.sh
    └── demo/                 # デモ実行スクリプト
        ├── run-phase1.sh
        ├── run-phase2.sh
        ├── run-phase3.sh
        ├── run-phase4.sh
        └── run-phase5.sh
```

## 学習体験設計：各フェーズの「動くもの」

### 学習体験の基本コンセプト
- **即座に動く**: 各フェーズの最後には必ず動かせるプログラムがある
- **シンプルな出力**: ターミナルでの基本的な入出力に集中
- **達成感の演出**: 「Hello, World!」から始まり、段階的に複雑になる
- **実用的価値**: おもちゃではなく、実際に役立つツールを作る

## 実装段階（フェーズ）と達成デモ

### フェーズ1: TCP基礎実装 - TCPエコーサーバ
**目標**: 基本的なTCP接続とデータ送受信

**実装内容**:
- TCP接続の確立・切断
- バイナリデータの送受信
- 基本的なエラーハンドリング

**🎯 達成デモ: TCP Echo Server & Client**
```
📂 demo/phase1-tcp-echo/

🚀 実行方法:
$ make demo-phase1
# または
$ ./scripts/demo/run-phase1.sh

✨ 体験内容:
1. TCP Echo サーバが起動 (localhost:8080)
2. クライアントから任意の文字列を送信
3. サーバが同じ文字列をエコーバック
4. 送受信ログをコンソールに表示
```

**実際の体験例**:
```bash
# ターミナル1: サーバ起動
$ go run demo/phase1-tcp-echo/server/main.go
TCP Echo Server started on :8080
Waiting for connections...

# ターミナル2: クライアント実行
$ go run demo/phase1-tcp-echo/client/main.go
Connected to server!
> Hello, TinyServer!
Echo: Hello, TinyServer!
> This is my first TCP implementation!
Echo: This is my first TCP implementation!
```

### フェーズ2: HTTPプロトコル実装 - HTTPパーサー
**目標**: HTTP/1.1の基本的な要素を自作実装

**実装内容**:
- HTTPリクエストのパース（GET, POST）
- HTTPレスポンスの生成
- ヘッダーの処理

**🎯 達成デモ: HTTP Request Parser & Analyzer**
```
📂 demo/phase2-http-parser/

🚀 実行方法:
$ make demo-phase2

✨ 体験内容:
1. 様々なHTTPリクエストをファイルから読み込み
2. 自作パーサーで解析してテキスト表示
3. 不正なリクエストのエラー検出デモ
4. リクエスト/レスポンスの構造をプレーンテキストで出力
```

**実際の体験例**:
```bash
$ go run demo/phase2-http-parser/main.go

HTTP Request Parser Demo
========================

Parsing Sample Request:
GET /api/users?id=123 HTTP/1.1
Host: example.com
User-Agent: TinyClient/1.0

Parse Result:
Method: GET
Path: /api/users
Query: id=123
Headers: 
  Host: example.com
  User-Agent: TinyClient/1.0

Generated Response:
HTTP/1.1 200 OK
Content-Type: text/plain
Content-Length: 42

Successfully parsed request!
```

### フェーズ3: シンプルサーバ実装 - 静的ファイルサーバ
**目標**: 基本的なWebサーバの動作

**実装内容**:
- リクエストの受信とルーティング
- 静的ファイルの配信
- 基本的なミドルウェア機能

**🎯 達成デモ: Basic File Server**
```
📂 demo/phase3-simple-server/

🚀 実行方法:
$ make demo-phase3

✨ 体験内容:
1. 基本的なHTMLファイルを配信するWebサーバ
2. ブラウザで http://localhost:8080 にアクセス
3. シンプルなHTMLページが表示
4. アクセスログをコンソールに表示
```

**実際の体験例**:
```bash
$ go run demo/phase3-simple-server/main.go
TinyServer Static File Server
Server started on http://localhost:8080
Serving files from: ./demo/phase3-simple-server/static

# ブラウザでアクセスすると...
Access Log:
[2024-01-15 10:30:45] GET / -> 200 (index.html, 2.1KB)
[2024-01-15 10:30:46] GET /about -> 404 (not found)
```

**配信される内容**:
- シンプルなHTMLページ
- プレーンテキストベースのコンテンツ
- TinyServerプロジェクトの進捗確認ページ

### フェーズ4: クライアント実装 - HTTP クライアントツール
**目標**: HTTPクライアントの実装

**実装内容**:
- HTTPリクエストの送信
- レスポンスの受信と解析
- タイムアウト処理

**🎯 達成デモ: TinyCurl - HTTP Client Tool**
```
📂 demo/phase4-http-client/

🚀 実行方法:
$ make demo-phase4

✨ 体験内容:
1. 自作curlライクなHTTPクライアント
2. 様々なWebサイトにリクエスト送信
3. レスポンスをプレーンテキストで表示
4. 基本的なパフォーマンス測定
```

**実際の体験例**:
```bash
$ go run demo/phase4-http-client/main.go

TinyCurl - HTTP Client Demo
===========================

Sending request to https://httpbin.org/json
Response time: 234ms

Response:
Status: 200 OK
Headers:
  Content-Type: application/json
  Content-Length: 429
  
Body:
{
  "slideshow": {
    "author": "Yours Truly",
    "date": "date of publication",
    "slides": [...]
  }
}

Performance:
  Total time: 234ms
```

### フェーズ5: 統合・最適化 - シンプルチャットアプリ
**目標**: 実用的なWebアプリケーションの完成

**実装内容**:
- パフォーマンス最適化
- ロングポーリング
- JSON API
- フロントエンド連携

**🎯 達成デモ: TinyChat - Simple Chat Application**
```
📂 demo/phase5-full-stack/

🚀 実行方法:
$ make demo-phase5

✨ 体験内容:
1. シンプルなチャットアプリケーション
2. 複数タブで同時接続してチャット
3. メッセージ履歴の保存・読み込み
4. プレーンHTMLフロントエンド
```

**実際の体験例**:
```bash
$ go run demo/phase5-full-stack/server/main.go
TinyChat Server
Server started on http://localhost:8080
Long Polling enabled
Message history: enabled

# ブラウザで複数タブを開いてチャット！
Message Log:
[10:45:12] User-7f3a: Hello from TinyServer!
[10:45:18] User-2b8c: This is working!
[10:45:25] User-7f3a: We built this from scratch!
```

**フロントエンド機能**:
- シンプルなメッセージ表示
- 基本的なユーザー識別
- メッセージ履歴表示
- 送信ボタン・Enterキー対応
- プレーンHTMLとシンプルなCSS

## 教育用ガイド（/books）の学習体験設計

### 各章の拡張構成
1. **概念説明**: 技術概念 + 設計原則
2. **設計思想**: なぜこの抽象化を選ぶのか
3. **実装方針**: SOLID原則の適用方法
4. **📝 ハンズオン実装**: ステップバイステップの実装ガイド
5. **🧪 テスト戦略**: 何をどのようにテストするか
6. **🔧 リファクタリング**: より良い設計への改善
7. **🎯 達成デモ**: 動くものを作って達成感を得る
8. **🎉 実行体験**: 実際に動かしてみる楽しさ
9. **📊 学習チェック**: 理解度確認とトラブルシューティング
10. **➡️ 次の段階への準備**: 次のフェーズへの橋渡し

### 学習体験の工夫

#### 01-tcp-basics.md
```markdown
# 🎯 フェーズ1の達成目標
あなた専用のTCP Echo Serverを作り、実際にメッセージを送受信してみましょう！

## 🚀 最終的に作るもの
- メッセージをエコーしてくれるサーバ
- リアルタイム通信ログ
- 複数クライアント同時接続

## 📝 実装チェックリスト
- [ ] TCP接続を確立できる
- [ ] メッセージを送受信できる
- [ ] エラーを適切に処理できる
- [ ] ログを出力できる

## 🎉 達成の瞬間
実装が完了したら、以下のコマンドを実行してください：

$ make demo-phase1

自分で作ったTCPサーバと会話してみましょう！
```

#### 各章に「🎮 Fun Challenge」セクション
```markdown
## 🎮 Fun Challenge
基本実装ができたら、以下にチャレンジしてみましょう：

### 初級
- [ ] エコーメッセージに現在時刻を追加
- [ ] 接続数をカウントして表示
- [ ] ログ出力の改善

### 中級  
- [ ] 特定のキーワードに特別な応答
- [ ] 接続履歴の保存
- [ ] 簡易的なルーム機能

### 上級
- [ ] バイナリデータの送受信
- [ ] 圧縮機能の追加
- [ ] 暗号化通信（簡易版）
```

## デモ実行スクリプト

### 自動化されたデモ体験
```bash
# scripts/demo/run-phase1.sh
#!/bin/bash
echo "Phase 1: TCP Echo Server Demo"
echo "=============================="
echo ""
echo "Starting server in background..."
go run demo/phase1-tcp-echo/server/main.go &
SERVER_PID=$!

sleep 2
echo "Server started! (PID: $SERVER_PID)"
echo ""
echo "Now sending test messages..."
echo ""

# 自動的にテストメッセージを送信
echo "Hello, TinyServer!" | go run demo/phase1-tcp-echo/client/main.go
echo "My first TCP implementation!" | go run demo/phase1-tcp-echo/client/main.go
echo "This is working!" | go run demo/phase1-tcp-echo/client/main.go

echo ""
echo "Demo completed successfully!"
echo "Try running the client manually for interactive experience:"
echo "  go run demo/phase1-tcp-echo/client/main.go"

# サーバを停止
kill $SERVER_PID
```

## 学習効果の最大化

### 達成感の設計
1. **小さな成功の積み重ね**: 各フェーズで確実に動くものを作る
2. **シンプルな出力**: コンソール出力に集中し、低レイヤーを理解
3. **実用性**: おもちゃではなく実際に使えるツール
4. **カスタマイズ性**: 学習者が改造して遊べる余地

## プロジェクト管理指針

### 📋 タスク管理の必須要件
TinyServerは教育プロジェクトとして、学習者の進捗管理と達成感の維持が重要です：

1. **必須**: `TASK.md`を使用した段階的進捗管理
   - 各タスク完了時に進捗欄（[ ] → [x]）を必ず更新
   - 各フェーズ完了時に「📊 プロジェクト進捗状況」セクションを最新化
   - 実装済みファイル一覧を常に最新状態に保持

2. **推奨**: 進捗の可視化
   - 完了したタスクを明確に識別
   - 次に取り組むタスクの明確化
   - 学習者のモチベーション維持
