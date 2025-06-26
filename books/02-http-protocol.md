# フェーズ2: HTTPプロトコル実装ガイド

このガイドでは、HTTPプロトコルの仕様を理解し、ゼロからHTTPパーサーを実装する方法を学びます。このガイドに従って実装すれば、TinyServerのフェーズ2を完全に再現できます。

## 🎯 学習目標

- HTTP/1.1プロトコルの仕様と構造を完全に理解する
- HTTPリクエスト・レスポンスの解析手法を習得する
- 堅牢なHTTPパーサーを実装できるようになる
- HTTPヘッダーとステータスコードの体系的知識を獲得する
- 実際のHTTPトラフィックを解析できるスキルを身につける

## 📚 目次

1. [HTTP/1.1プロトコルの基礎](#http11プロトコルの基礎)
2. [HTTPメッセージの構造](#httpメッセージの構造)
3. [HTTPリクエストの詳細](#httpリクエストの詳細)
4. [HTTPレスポンスの詳細](#httpレスポンスの詳細)
5. [HTTPヘッダーの完全解説](#httpヘッダーの完全解説)
6. [HTTPステータスコード](#httpステータスコード)
7. [HTTPパーサーの設計と実装](#httpパーサーの設計と実装)
8. [実装演習](#実装演習)
9. [トラブルシューティング](#トラブルシューティング)

## HTTP/1.1プロトコルの基礎

### HTTPとは

HTTP（HyperText Transfer Protocol）は、WWW上でデータを転送するためのアプリケーション層プロトコルです。HTTP/1.1は1997年にRFC 2068で標準化され、1999年にRFC 2616で改訂されました。

### HTTPの特徴

1. **ステートレス**: 各リクエストは独立しており、サーバーは前のリクエストの状態を保持しません
2. **テキストベース**: HTTPメッセージは人間が読める形式のテキストです
3. **リクエスト・レスポンス型**: クライアントがリクエストを送信し、サーバーがレスポンスを返します
4. **階層化**: TCP/IPプロトコルスタックの上位に位置します

### HTTP通信の流れ

```
1. クライアント → サーバー: HTTPリクエスト
2. サーバー → クライアント: HTTPレスポンス
```

## HTTPメッセージの構造

HTTPメッセージは、リクエストとレスポンスの共通構造を持ちます：

```
開始行（Request Line / Status Line）
ヘッダーフィールド1: 値1
ヘッダーフィールド2: 値2
...
空行（CRLF）
メッセージボディ（オプション）
```

### 重要な規則

1. **改行コード**: HTTP/1.1では必ずCRLF（`\r\n`）を使用
2. **ヘッダーとボディの境界**: 空行（`\r\n\r\n`）で区切る
3. **大文字小文字**: ヘッダー名は大文字小文字を区別しない
4. **エンコーディング**: ASCII文字セットを基本とする

## HTTPリクエストの詳細

### リクエストラインの構造

```
METHOD SP REQUEST-URI SP HTTP-VERSION CRLF
```

例:
```
GET /api/users?id=123 HTTP/1.1
```

#### HTTPメソッド

| メソッド | 用途 | 冪等性 | 安全性 |
|---------|------|-------|-------|
| GET | リソースの取得 | ○ | ○ |
| POST | リソースの作成・処理 | × | × |
| PUT | リソースの作成・更新 | ○ | × |
| DELETE | リソースの削除 | ○ | × |
| HEAD | ヘッダーのみ取得 | ○ | ○ |
| OPTIONS | サポートメソッド確認 | ○ | ○ |
| PATCH | リソースの部分更新 | × | × |

#### URI（Uniform Resource Identifier）

URIの構成要素：
```
scheme://authority/path?query#fragment
```

HTTPリクエストでは通常、pathとqueryの部分のみを送信：
```
/path/to/resource?param1=value1&param2=value2
```

#### HTTPバージョン

- `HTTP/1.0`: シンプルな実装、接続ごとに切断
- `HTTP/1.1`: 持続的接続、パイプライニング、チャンク転送

### リクエストヘッダー

必須ヘッダー（HTTP/1.1）：
- `Host`: ターゲットサーバーのホスト名

一般的なヘッダー：
```
Host: example.com
User-Agent: TinyClient/1.0
Accept: application/json
Content-Type: application/json
Content-Length: 42
Authorization: Bearer token123
```

### リクエストボディ

POST、PUT、PATCHリクエストで使用されることが多い：

```http
POST /api/users HTTP/1.1
Host: example.com
Content-Type: application/json
Content-Length: 58

{"name":"John Doe","email":"john@example.com","age":30}
```

## HTTPレスポンスの詳細

### ステータスラインの構造

```
HTTP-VERSION SP STATUS-CODE SP REASON-PHRASE CRLF
```

例:
```
HTTP/1.1 200 OK
```

### レスポンスヘッダー

一般的なレスポンスヘッダー：
```
HTTP/1.1 200 OK
Server: TinyServer/1.0
Date: Mon, 26 Jun 2025 12:00:00 GMT
Content-Type: application/json
Content-Length: 55
Connection: close
```

### レスポンスボディ

サーバーからクライアントに返されるデータ：

```http
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 55

{"id":123,"name":"John Doe","email":"john@example.com"}
```

## HTTPヘッダーの完全解説

### ヘッダーの分類

#### 1. 一般ヘッダー
リクエスト・レスポンス両方で使用：
- `Date`: メッセージの生成日時
- `Connection`: 接続管理
- `Cache-Control`: キャッシュ制御
- `Transfer-Encoding`: 転送エンコーディング

#### 2. リクエストヘッダー
クライアントの情報や要求：
- `Host`: ターゲットホスト（必須）
- `User-Agent`: クライアント情報
- `Accept`: 受け入れ可能なメディアタイプ
- `Accept-Language`: 受け入れ可能な言語
- `Accept-Encoding`: 受け入れ可能な圧縮方式
- `Authorization`: 認証情報
- `If-Modified-Since`: 条件付きリクエスト

#### 3. レスポンスヘッダー
サーバーの情報：
- `Server`: サーバー情報
- `Location`: リダイレクト先
- `WWW-Authenticate`: 認証要求

#### 4. エンティティヘッダー
メッセージボディの情報：
- `Content-Type`: メディアタイプ
- `Content-Length`: ボディのバイト数
- `Content-Encoding`: 圧縮方式
- `Last-Modified`: 最終更新日時
- `ETag`: エンティティタグ

### 重要なヘッダーの詳細

#### Content-Length
メッセージボディのバイト数を指定：
```
Content-Length: 1234
```

**重要**: バイト数であり、文字数ではありません（UTF-8などでは異なる）

#### Content-Type
メディアタイプとオプションの文字セット：
```
Content-Type: text/html; charset=UTF-8
Content-Type: application/json
Content-Type: multipart/form-data; boundary=something
```

#### Transfer-Encoding
データの転送方式：
```
Transfer-Encoding: chunked
```

チャンク転送では、データを複数の「チャンク」に分割して送信：
```
5\r\n
Hello\r\n
6\r\n
 World\r\n
0\r\n
\r\n
```

## HTTPステータスコード

### ステータスコードの分類

#### 1xx: 情報的レスポンス
- `100 Continue`: リクエストの続行を許可
- `101 Switching Protocols`: プロトコル切り替え

#### 2xx: 成功
- `200 OK`: リクエスト成功
- `201 Created`: リソース作成成功
- `202 Accepted`: リクエスト受理（処理中）
- `204 No Content`: 成功、ボディなし

#### 3xx: リダイレクション
- `301 Moved Permanently`: 永続的移動
- `302 Found`: 一時的移動
- `304 Not Modified`: 未変更
- `307 Temporary Redirect`: 一時的リダイレクト

#### 4xx: クライアントエラー
- `400 Bad Request`: 不正なリクエスト
- `401 Unauthorized`: 認証が必要
- `403 Forbidden`: アクセス禁止
- `404 Not Found`: リソースが見つからない
- `405 Method Not Allowed`: メソッドが許可されていない
- `409 Conflict`: リソースの競合
- `422 Unprocessable Entity`: 処理不可能なエンティティ

#### 5xx: サーバーエラー
- `500 Internal Server Error`: サーバー内部エラー
- `501 Not Implemented`: 実装されていない機能
- `502 Bad Gateway`: 不正なゲートウェイ
- `503 Service Unavailable`: サービス利用不可
- `504 Gateway Timeout`: ゲートウェイタイムアウト

## HTTPパーサーの設計と実装

### 設計原則

1. **段階的解析**: リクエストライン → ヘッダー → ボディの順で解析
2. **堅牢性**: 不正な形式に対する適切なエラーハンドリング
3. **効率性**: 大きなデータでもメモリ効率よく処理
4. **拡張性**: 新しいヘッダーやメソッドに対応しやすい設計

### パーサーのアーキテクチャ

```go
type HTTPParser interface {
    Parse(r io.Reader) (Request, error)
    ParseBytes(data []byte) (Request, error)
    Validate(req Request) error
}
```

### 実装手順

#### 1. メッセージの分割

HTTPメッセージをヘッダー部とボディ部に分割：

```go
func ParseRequest(r io.Reader, remoteAddr net.Addr) (Request, error) {
    // 全体を読み込み
    buf := &bytes.Buffer{}
    if _, err := io.Copy(buf, r); err != nil {
        return nil, HTTPError("failed to read request: " + err.Error())
    }
    
    data := buf.Bytes()
    
    // ヘッダー・ボディ分離点を検索
    headerEndIndex := bytes.Index(data, []byte("\r\n\r\n"))
    if headerEndIndex == -1 {
        return nil, HTTPError(ErrInvalidRequestLine)
    }
    
    headerData := data[:headerEndIndex]
    bodyData := data[headerEndIndex+4:] // \r\n\r\nをスキップ
    
    // ... 解析処理
}
```

#### 2. リクエストラインの解析

```go
func parseRequestLine(line string) (Method, string, Version, error) {
    if line == "" {
        return "", "", "", HTTPError(ErrInvalidRequestLine)
    }
    
    // スペースで分割
    parts := strings.SplitN(line, " ", 3)
    if len(parts) != 3 {
        return "", "", "", HTTPError(ErrInvalidRequestLine)
    }
    
    methodStr := parts[0]
    path := parts[1]
    versionStr := parts[2]
    
    // バリデーション
    method := Method(methodStr)
    if !isValidMethod(method) {
        return "", "", "", HTTPError(ErrInvalidMethod)
    }
    
    if !isValidPath(path) {
        return "", "", "", HTTPError(ErrInvalidPath)
    }
    
    version := Version(versionStr)
    if !isValidVersion(version) {
        return "", "", "", HTTPError(ErrInvalidVersion)
    }
    
    return method, path, version, nil
}
```

#### 3. ヘッダーの解析

```go
func parseHeaders(scanner *bufio.Scanner) (Header, error) {
    headers := make(Header)
    headerCount := 0
    
    for scanner.Scan() {
        line := scanner.Text()
        
        // 空行でヘッダー終了
        if line == "" {
            break
        }
        
        // ヘッダー数制限チェック
        headerCount++
        if headerCount > MaxHeaderLines {
            return nil, HTTPError(ErrHeaderTooLarge)
        }
        
        // コロンで分割
        colonIndex := strings.Index(line, ":")
        if colonIndex == -1 {
            return nil, HTTPError(ErrInvalidHeader)
        }
        
        name := strings.TrimSpace(line[:colonIndex])
        value := strings.TrimSpace(line[colonIndex+1:])
        
        // ヘッダー名のバリデーション
        if !isValidHeaderName(name) {
            return nil, HTTPError(ErrInvalidHeader)
        }
        
        headers[name] = append(headers[name], value)
    }
    
    return headers, nil
}
```

#### 4. ボディの解析

Content-Lengthに基づく解析：

```go
// Content-Lengthと実際のボディサイズの整合性チェック
contentLength := req.ContentLength()
if contentLength > 0 {
    if int64(len(bodyData)) != contentLength {
        return nil, HTTPError(ErrUnexpectedEOF)
    }
    req.SetBody(bytes.NewReader(bodyData))
}
```

### バリデーション

#### メソッドのバリデーション

```go
func isValidMethod(method Method) bool {
    switch method {
    case MethodGet, MethodPost, MethodPut, 
         MethodDelete, MethodHead, MethodOptions, 
         MethodPatch:
        return true
    default:
        return false
    }
}
```

#### パスのバリデーション

```go
func isValidPath(path string) bool {
    if path == "" {
        return false
    }
    
    // パスは/で始まる必要がある
    if !strings.HasPrefix(path, "/") {
        return false
    }
    
    // 制御文字の禁止
    for _, r := range path {
        if r < 32 || r == 127 {
            return false
        }
    }
    
    return true
}
```

#### ヘッダー名のバリデーション

```go
func isValidHeaderName(name string) bool {
    if name == "" {
        return false
    }
    
    // RFC 7230に準拠: 文字、数字、ハイフンのみ
    for _, r := range name {
        if !((r >= 'a' && r <= 'z') ||
             (r >= 'A' && r <= 'Z') ||
             (r >= '0' && r <= '9') ||
             r == '-') {
            return false
        }
    }
    
    return true
}
```

### エラーハンドリング

エラー分類システム：

```go
type ErrorType int

const (
    ErrorTypeNetwork ErrorType = iota
    ErrorTypeProtocol
    ErrorTypeServer
    ErrorTypeClient
    ErrorTypeIO
    ErrorTypeTimeout
    ErrorTypeInvalidInput
)
```

統一エラー処理：

```go
func HTTPError(message string) error {
    return NewTinyServerError(ErrorTypeProtocol, message)
}
```

## 実装演習

### 演習1: 基本的なGETリクエストの解析

以下のHTTPリクエストを解析するコードを書いてください：

```http
GET /hello HTTP/1.1
Host: example.com
User-Agent: TinyClient/1.0

```

期待される結果：
- Method: GET
- Path: /hello
- Version: HTTP/1.1
- Headers: Host, User-Agent

### 演習2: POSTリクエストのボディ解析

以下のHTTPリクエストを解析し、ボディを正しく取得してください：

```http
POST /api/users HTTP/1.1
Host: api.example.com
Content-Type: application/json
Content-Length: 58

{"name":"Jane Doe","email":"jane@example.com","age":25}
```

### 演習3: 不正なリクエストのエラーハンドリング

以下の不正なリクエストに対して適切なエラーを返すコードを実装してください：

1. メソッドが不正: `INVALID /path HTTP/1.1`
2. バージョンが不正: `GET /path HTTP/2.0`
3. ヘッダー形式が不正: `Host example.com` (コロンなし)
4. Content-Lengthとボディサイズの不整合

### 演習4: レスポンスパーサーの実装

HTTPレスポンスを解析するパーサーを実装してください：

```http
HTTP/1.1 200 OK
Server: TinyServer/1.0
Content-Type: text/plain
Content-Length: 13

Hello, World!
```

## 実際のHTTPトラフィック例

### 例1: REST API のやり取り

#### リクエスト
```http
GET /api/users/123 HTTP/1.1
Host: api.example.com
Accept: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
User-Agent: MyApp/1.0
```

#### レスポンス
```http
HTTP/1.1 200 OK
Server: nginx/1.18.0
Date: Mon, 26 Jun 2025 12:00:00 GMT
Content-Type: application/json
Content-Length: 156
Access-Control-Allow-Origin: *

{
  "id": 123,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-01-15T10:30:00Z",
  "last_login": "2025-06-26T11:45:00Z"
}
```

### 例2: フォームデータの送信

#### リクエスト
```http
POST /submit-form HTTP/1.1
Host: www.example.com
Content-Type: application/x-www-form-urlencoded
Content-Length: 47

username=johndoe&password=secret123&remember=on
```

#### レスポンス
```http
HTTP/1.1 302 Found
Location: /dashboard
Set-Cookie: session_id=abc123; Path=/; HttpOnly
Content-Length: 0

```

### 例3: ファイルアップロード

#### リクエスト
```http
POST /upload HTTP/1.1
Host: upload.example.com
Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Length: 234

------WebKitFormBoundary7MA4YWxkTrZu0gW
Content-Disposition: form-data; name="file"; filename="test.txt"
Content-Type: text/plain

Hello, this is test file content.
------WebKitFormBoundary7MA4YWxkTrZu0gW--
```

### 例4: エラーレスポンス

#### リクエスト
```http
GET /nonexistent HTTP/1.1
Host: example.com
```

#### レスポンス
```http
HTTP/1.1 404 Not Found
Server: TinyServer/1.0
Content-Type: text/html
Content-Length: 196

<!DOCTYPE html>
<html>
<head>
    <title>404 Not Found</title>
</head>
<body>
    <h1>404 Not Found</h1>
    <p>The requested resource was not found</p>
    <hr>
    <p><em>TinyServer</em></p>
</body>
</html>
```

## トラブルシューティング

### よくある問題と解決方法

#### 1. Content-Lengthの不整合

**問題**: ボディのサイズとContent-Lengthヘッダーの値が一致しない

**原因**:
- 文字エンコーディングの考慮不足（バイト数 vs 文字数）
- 改行コードの処理誤り

**解決方法**:
```go
// 正しい: バイト数を計算
body := `{"name":"太郎"}`
contentLength := len([]byte(body)) // UTF-8でのバイト数

// 間違い: 文字数を計算
contentLength := len(body) // 文字数（マルチバイト文字で問題）
```

#### 2. 改行コードの問題

**問題**: `\n`と`\r\n`の混在

**解決方法**:
- HTTP/1.1では必ず`\r\n`を使用
- パーサーでは両方に対応する場合もある

#### 3. ヘッダーの大文字小文字

**問題**: ヘッダー名の大文字小文字が統一されていない

**解決方法**:
```go
// 正規化して比較
headerName := strings.ToLower(name)
if headerName == "content-type" {
    // 処理
}
```

#### 4. 無限ループの回避

**問題**: 不正なデータで解析が無限ループに陥る

**解決方法**:
- 最大読み込みサイズの制限
- タイムアウトの設定
- ループカウンタの設置

#### 5. メモリ使用量の制御

**問題**: 大きなリクエストでメモリ不足

**解決方法**:
```go
const MaxRequestSize = 10 << 20 // 10MB

func ParseRequest(r io.Reader, remoteAddr net.Addr) (Request, error) {
    limitedReader := io.LimitReader(r, MaxRequestSize)
    // ...
}
```

### デバッグ手法

#### 1. リクエストのダンプ

```go
func dumpRequest(data []byte) {
    fmt.Printf("Raw request (%d bytes):\n", len(data))
    fmt.Printf("%q\n", string(data))
    fmt.Printf("Hex dump:\n%s\n", hex.Dump(data))
}
```

#### 2. 段階的解析

```go
func debugParse(data []byte) {
    fmt.Printf("1. Finding header boundary...\n")
    headerEnd := bytes.Index(data, []byte("\r\n\r\n"))
    fmt.Printf("   Found at position: %d\n", headerEnd)
    
    if headerEnd != -1 {
        headerData := data[:headerEnd]
        bodyData := data[headerEnd+4:]
        
        fmt.Printf("2. Header data (%d bytes):\n%q\n", 
                   len(headerData), string(headerData))
        fmt.Printf("3. Body data (%d bytes):\n%q\n", 
                   len(bodyData), string(bodyData))
    }
}
```

#### 3. パフォーマンス測定

```go
func benchmarkParse(data []byte) {
    start := time.Now()
    req, err := ParseRequest(bytes.NewReader(data), nil)
    duration := time.Since(start)
    
    fmt.Printf("Parse time: %v\n", duration)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Success: %s %s\n", req.Method(), req.Path())
    }
}
```

## まとめ

このガイドでは、HTTP/1.1プロトコルの詳細な仕様から実装方法まで、包括的に学習しました。

### 重要なポイント

1. **プロトコル理解**: HTTP/1.1の仕様を正確に理解することが重要
2. **段階的実装**: リクエストライン → ヘッダー → ボディの順で実装
3. **堅牢性**: 不正な入力に対する適切なエラーハンドリング
4. **効率性**: メモリとCPU使用量を考慮した実装
5. **テスト**: 様々なケースでのテストが重要

### 次のステップ

フェーズ2を完了したら、次はフェーズ3でHTTPサーバーの実装に進みます：

1. **サーバーアーキテクチャ**: リクエストハンドラーとミドルウェア
2. **ルーティング**: URLパターンマッチング
3. **静的ファイル配信**: ファイルシステムとの連携
4. **エラーハンドリング**: 適切なエラーレスポンス

このガイドで学んだHTTPパーサーの知識は、フェーズ3以降でも活用されます。しっかりと理解を深めて次のステップに進みましょう。