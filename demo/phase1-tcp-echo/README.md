# Phase 1: TCP Echo Server & Client Demo

このデモでは、TinyServerプロジェクトのフェーズ1で実装したTCP基盤を使用して、シンプルなエコーサーバーとクライアントを動作させます。

## 概要

- **サーバー**: 接続されたクライアントから受信したデータをそのまま送り返すエコーサーバー
- **クライアント**: サーバーに接続してメッセージを送信し、エコー応答を受信するクライアント

## 実行方法

### 1. サーバーの起動

```bash
# デフォルト設定で起動（localhost:8080）
go run demo/phase1-tcp-echo/server/main.go

# カスタムポートで起動
go run demo/phase1-tcp-echo/server/main.go -port 9090

# 詳細ログ付きで起動
go run demo/phase1-tcp-echo/server/main.go -verbose

# カスタムホストとポートで起動
go run demo/phase1-tcp-echo/server/main.go -host 0.0.0.0 -port 8080
```

### 2. クライアントの実行

別のターミナルで以下を実行：

```bash
# インタラクティブモードで接続
go run demo/phase1-tcp-echo/client/main.go

# 単一メッセージを送信
go run demo/phase1-tcp-echo/client/main.go -message "Hello, TinyServer!"

# カスタムサーバーに接続
go run demo/phase1-tcp-echo/client/main.go -host localhost -port 9090

# 詳細ログ付きで実行
go run demo/phase1-tcp-echo/client/main.go -verbose
```

## 使用例

### サーバー起動
```bash
$ go run demo/phase1-tcp-echo/server/main.go
[2024-01-15 10:30:00] INFO: Starting TCP Echo Server on localhost:8080
[2024-01-15 10:30:00] INFO: TCP Echo Server is running...
[2024-01-15 10:30:00] INFO: Press Ctrl+C to stop the server
```

### クライアント接続（インタラクティブモード）
```bash
$ go run demo/phase1-tcp-echo/client/main.go
[2024-01-15 10:30:05] INFO: Connecting to TCP Echo Server at localhost:8080
[2024-01-15 10:30:05] INFO: Connected to server successfully!
[2024-01-15 10:30:05] INFO: Interactive mode started. Type messages to echo. Type 'quit' to exit.

TCP Echo Client - Interactive Mode
=================================
Type your message and press Enter. Type 'quit' to exit.

> Hello, TinyServer!
Echo: Hello, TinyServer!

> This is my first TCP implementation!
Echo: This is my first TCP implementation!

> quit
Goodbye!
```

### 単一メッセージモード
```bash
$ go run demo/phase1-tcp-echo/client/main.go -message "Hello, World!"
[2024-01-15 10:30:10] INFO: Connecting to TCP Echo Server at localhost:8080
[2024-01-15 10:30:10] INFO: Connected to server successfully!
[2024-01-15 10:30:10] INFO: Echo response: "Hello, World!"
[2024-01-15 10:30:10] INFO: ✓ Echo successful!
```

## サーバーログ出力例

```bash
[2024-01-15 10:30:05] INFO: New client connected: 127.0.0.1:54321
[2024-01-15 10:30:06] INFO: Client disconnected: 127.0.0.1:54321
[2024-01-15 10:30:10] INFO: New client connected: 127.0.0.1:54322
[2024-01-15 10:30:11] INFO: Client disconnected: 127.0.0.1:54322
```

## コマンドラインオプション

### サーバー
- `-host`: バインドするホスト（デフォルト: localhost）
- `-port`: リスニングポート（デフォルト: 8080）
- `-verbose`: 詳細ログを有効化

### クライアント
- `-host`: 接続先ホスト（デフォルト: localhost）
- `-port`: 接続先ポート（デフォルト: 8080）
- `-message`: 単一メッセージ送信モード
- `-verbose`: 詳細ログを有効化

## テスト機能

### 複数クライアント接続テスト
複数のターミナルで同時にクライアントを実行し、サーバーが複数の接続を正しく処理できることを確認できます。

### 自動テストスクリプト
```bash
# 自動デモスクリプトを実行
./scripts/demo/run-phase1.sh
```

## 学習ポイント

1. **TCP接続の基本**: ソケットプログラミングの基礎
2. **エラーハンドリング**: ネットワーク操作のエラー処理
3. **並行処理**: 複数クライアント同時接続の処理
4. **リソース管理**: 接続の適切なクローズ処理
5. **ログ出力**: デバッグとモニタリングのためのログ

## トラブルシューティング

### 接続に失敗する場合
- サーバーが起動しているか確認
- ポート番号が正しいか確認
- ファイアウォールの設定を確認

### サーバーが起動しない場合
- ポートが他のプロセスで使用されていないか確認
- 権限不足でないか確認（特に1024以下のポート）

### エコーが正しく動作しない場合
- ネットワーク接続を確認
- 詳細ログ（-verbose）で通信内容を確認

## 次のステップ

このデモが正常に動作したら、フェーズ2のHTTPプロトコル実装に進むことができます。TCP基盤の上にHTTPプロトコルを実装していきます。