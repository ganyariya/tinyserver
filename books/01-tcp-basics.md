# 🎯 フェーズ1: TCP基礎実装ガイド

このガイドでは、TinyServerプロジェクトのフェーズ1で実装するTCP基礎機能について、理論から実装まで詳しく解説します。

## 🚀 このフェーズで作るもの

**最終目標**: 自分専用のTCP Echo Serverを作り、実際にメッセージを送受信してみましょう！

### 達成デモ
- メッセージをエコーしてくれるサーバ
- リアルタイム通信ログ
- 複数クライアント同時接続
- インタラクティブなクライアント

## 📚 1. TCP/IPプロトコルの基礎理解

### 1.1 TCP/IPとは何か

**TCP (Transmission Control Protocol)** は、信頼性の高いデータ通信を実現するプロトコルです。

#### TCPの特徴
- **接続指向**: 通信前にコネクションを確立
- **信頼性**: データの到達保証、順序保証
- **フロー制御**: 受信側の処理能力に合わせた送信
- **エラー検出・修復**: 破損したデータの再送

#### IPとの関係
- **IP (Internet Protocol)**: データを宛先まで届ける（住所のような役割）
- **TCP**: 信頼性のあるデータ配送を保証（宅配業者のような役割）

### 1.2 ソケットプログラミングの概念

**ソケット**は、ネットワーク通信のエンドポイント（通信の終端）です。

```
クライアント                    サーバー
┌─────────────┐                ┌─────────────┐
│   Socket    │ ←── TCP/IP ──→ │   Socket    │
│ (送受信口)  │                │ (送受信口)  │
└─────────────┘                └─────────────┘
```

#### ソケットの種類
- **ストリームソケット** (TCP): 信頼性のある順序付きデータストリーム
- **データグラムソケット** (UDP): 軽量だが信頼性なし（今回は使用しない）

## 🏗️ 2. 設計思想とアーキテクチャ

### 2.1 SOLID原則の適用

TinyServerでは、以下のSOLID原則を実践します：

#### 単一責任原則 (SRP)
```go
// ❌ 悪い例: 1つの型が複数の責任を持つ
type BadTCPHandler struct {
    // 接続管理、ログ、データ処理、エラーハンドリングが混在
}

// ✅ 良い例: 責任を分離
type Connection interface {      // 接続管理のみ
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
}

type Logger interface {          // ログ出力のみ
    Info(string, ...interface{})
    Error(string, ...interface{})
}
```

#### インターフェース分離原則 (ISP)
```go
// 小さく特化したインターフェース
type Connection interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
}

type Listener interface {
    Accept() (Connection, error)
    Close() error
}

type Server interface {
    Start() error
    Stop() error
    SetHandler(ConnectionHandler)
}
```

### 2.2 パッケージ構造の設計思想

```
pkg/tcp/           # 公開API - ユーザーが使うインターフェース
internal/tcp/      # 内部実装 - 実装詳細を隠蔽
internal/common/   # 共通基盤 - エラー、ログなど
```

#### なぜこの構造なのか？
1. **公開APIと実装の分離**: ユーザーは`pkg/tcp`のインターフェースのみを使用
2. **実装詳細の隠蔽**: `internal/`は外部パッケージからアクセス不可
3. **テスト容易性**: インターフェースによるモック化が可能

## 🧪 3. テスト駆動開発 (TDD) の実践

### 3.1 TDDのサイクル

```
1. Red   -> テストを書く（失敗する）
2. Green -> 最小限の実装で通す
3. Refactor -> コードを改善する
```

### 3.2 TCP実装のテスト戦略

#### 接続テスト
```go
func TestTCPConnection(t *testing.T) {
    // 1. テスト用サーバを起動
    // 2. クライアントで接続
    // 3. データ送受信を確認
    // 4. 適切にクリーンアップ
}
```

#### エラーハンドリングテスト
```go
func TestConnectionError(t *testing.T) {
    // 1. 無効なアドレスに接続試行
    // 2. 適切なエラーが返されることを確認
}
```

## 🔧 4. ステップバイステップ実装ガイド

### Step 1: インターフェース設計

まず、必要なインターフェースを`pkg/tcp/interfaces.go`に定義します。

```go
package tcp

// Connection represents a TCP connection
type Connection interface {
    Read([]byte) (int, error)
    Write([]byte) (int, error)
    Close() error
    LocalAddr() net.Addr
    RemoteAddr() net.Addr
}

// Listener represents a TCP listener
type Listener interface {
    Accept() (Connection, error)
    Close() error
    Addr() net.Addr
}
```

**設計のポイント**:
- Go標準の`net.Conn`に似せつつ、必要な機能のみを抽象化
- テストしやすいインターフェース設計

### Step 2: TCP接続の実装

`internal/tcp/connection.go`でインターフェースを実装します。

```go
type tcpConnection struct {
    conn   net.Conn
    logger *common.Logger
}

func (c *tcpConnection) Read(b []byte) (int, error) {
    n, err := c.conn.Read(b)
    if err != nil {
        return n, common.NetworkError("read failed")
    }
    return n, nil
}
```

**実装のポイント**:
- 標準ライブラリ`net.Conn`をラップ
- エラーを独自のエラー型に変換
- ログ出力で動作を追跡可能

### Step 3: エラーハンドリング

`internal/common/errors.go`で統一されたエラー処理を実装します。

```go
type NetworkError struct {
    message string
    cause   error
}

func (e *NetworkError) Error() string {
    if e.cause != nil {
        return fmt.Sprintf("%s: %v", e.message, e.cause)
    }
    return e.message
}
```

### Step 4: ログ機能

`internal/common/logger.go`でデバッグしやすいログ機能を実装します。

```go
func (l *Logger) Info(format string, args ...interface{}) {
    message := fmt.Sprintf(format, args...)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] INFO: %s\n", timestamp, message)
}
```

### Step 5: TCP Echoサーバの実装

`demo/phase1-tcp-echo/server/main.go`で実際に動くサーバを作成します。

```go
func main() {
    server, err := tcp.NewServer("tcp", ":8080")
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    // エコーハンドラを設定
    server.SetHandler(func(conn tcp.Connection) {
        buffer := make([]byte, 1024)
        for {
            n, err := conn.Read(buffer)
            if err != nil {
                break
            }
            conn.Write(buffer[:n]) // エコーバック
        }
    })

    fmt.Println("TCP Echo Server started on :8080")
    server.Start()
}
```

### Step 6: TCP Echoクライアントの実装

`demo/phase1-tcp-echo/client/main.go`でインタラクティブなクライアントを作成します。

```go
func main() {
    conn, err := tcp.Dial("tcp", "localhost:8080")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("> ")
    
    for scanner.Scan() {
        message := scanner.Text()
        
        // メッセージ送信
        conn.Write([]byte(message))
        
        // エコー受信
        buffer := make([]byte, 1024)
        n, _ := conn.Read(buffer)
        fmt.Printf("Echo: %s\n", string(buffer[:n]))
        fmt.Print("> ")
    }
}
```

## 🎉 5. 達成の瞬間 - デモ実行

実装が完了したら、以下のコマンドで動作確認しましょう：

```bash
# 自動デモ実行
$ make demo-phase1
# または
$ ./scripts/demo/run-phase1.sh
```

### 手動でのデモ実行

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
> My first TCP implementation!
Echo: My first TCP implementation!
```

## 🎮 Fun Challenge

基本実装ができたら、以下にチャレンジしてみましょう：

### 初級チャレンジ
- [ ] エコーメッセージに現在時刻を追加
- [ ] 接続数をカウントして表示
- [ ] ログ出力の改善

### 中級チャレンジ
- [ ] 特定のキーワードに特別な応答
- [ ] 接続履歴の保存
- [ ] 簡易的なルーム機能

### 上級チャレンジ
- [ ] バイナリデータの送受信
- [ ] 圧縮機能の追加
- [ ] 暗号化通信（簡易版）

## 📝 実装チェックリスト

### 基本機能
- [ ] TCP接続を確立できる
- [ ] メッセージを送受信できる
- [ ] 複数クライアント同時接続
- [ ] 適切なエラーハンドリング
- [ ] ログ出力機能

### コード品質
- [ ] インターフェース駆動設計
- [ ] 適切なエラー型定義
- [ ] テストカバレッジ
- [ ] `gofmt`, `go vet`を通すコード
- [ ] マジック定数の排除

### デモ機能
- [ ] サーバが正常に起動
- [ ] クライアントが接続可能
- [ ] エコー機能が動作
- [ ] インタラクティブな操作
- [ ] ログ出力の確認

## 🔍 トラブルシューティング

### よくある問題と解決方法

#### 1. 接続エラー
```
Error: dial tcp :8080: connect: connection refused
```
**解決方法**: サーバが起動しているか確認

#### 2. ポートが使用中
```
Error: listen tcp :8080: bind: address already in use
```
**解決方法**: 別のポートを使用するか、プロセスを終了

#### 3. テストがタイムアウト
```
Error: test timeout
```
**解決方法**: 適切なタイムアウト設定とリソースクリーンアップ

## 📊 学習チェック

### 理解度確認
- [ ] TCPとUDPの違いを説明できる
- [ ] ソケットプログラミングの基本概念を理解している
- [ ] SOLID原則の適用方法を理解している
- [ ] Go言語でのインターフェース設計ができる
- [ ] エラーハンドリングの重要性を理解している

### 実装スキル
- [ ] TCP接続の確立・切断ができる
- [ ] データの送受信処理を実装できる
- [ ] 適切なエラーハンドリングを実装できる
- [ ] ログ機能を活用してデバッグできる
- [ ] テスト駆動開発を実践できる

## ➡️ 次の段階への準備

フェーズ1が完了したら、次はフェーズ2「HTTPプロトコル実装」に進みます。

### フェーズ2で学ぶこと
- HTTP/1.1プロトコルの構造
- リクエスト・レスポンスのパース
- ヘッダーの処理
- ステータスコードの理解

### 準備すること
- TCP基礎の理解を確実にする
- インターフェース設計の感覚を掴む
- テスト駆動開発の習慣化

**おめでとうございます！** あなたは今、低レベルなTCP通信を理解し、実装できるようになりました。これは現代のWeb技術の土台となる重要なスキルです。

フェーズ2では、このTCP基盤の上にHTTPプロトコルを実装し、Webの世界へ踏み出しましょう！