# Step 03 — mTLS (相互認証)

## ゴール

- クライアントにも証明書を提示させ、サーバー側で検証する **mTLS** を成立させる。
- `tls.Config.ClientAuth` の各モードの違いを説明できるようにする。
- ハンドシェイク完了後、`*tls.Conn.ConnectionState().PeerCertificates` から「誰が繋いできたか」を取り出して、認可ロジックの足がかりにする。

## 前提

```bash
make certs
```

`tutorial/step03-mtls/certs/` 直下に `ca.crt`, `server.crt`, `server.key`, `client.crt`, `client.key` の5つが揃っていれば OK です。

## ファイル

| ファイル | 役割 |
|---|---|
| `server/main.go` | `:9444` で mTLS listen。クライアント証明書を `ca.crt` で検証 |
| `client/main.go` | `localhost:9444` に接続。`client.crt` を提示 |

## 動かす

**ターミナル A**

```bash
go run ./tutorial/step03-mtls/server
# => step03 mTLS server listening on :9444 (mutual auth)
```

**ターミナル B**

```bash
go run ./tutorial/step03-mtls/client
# => connected to localhost:9444 over mTLS (server CN=localhost)
hi
hello alice@example.com, you said: hi
```

サーバー側のログにも `peer CN=alice@example.com` が出ているはずです。**サーバーが相手の名前を知っている**ことが、Step02 との大きな違いです。

## ハンドシェイクの追加分

Step02 のハンドシェイクに、サーバーからの `CertificateRequest` と、クライアントからの `Certificate` + `CertificateVerify` (秘密鍵で署名) が増えます。

```
Client                                  Server
  |--- ClientHello ----------------------->|
  |<-- ServerHello, Certificate -----------|
  |<-- CertificateRequest -----------------| ← Step03 から追加
  |--- Certificate, CertificateVerify ---->| ← Step03 から追加
  |--- (鍵交換) -------------------------->|
  |<== 暗号化セッション ===================|
```

サーバーは:

1. クライアントが提示した証明書チェーンが `ClientCAs` に辿り着くか
2. 署名 (`CertificateVerify`) が、提示された公開鍵に対応する秘密鍵で行われたか

を検証します。両方通って初めて接続が成立します。

## `ClientAuth` の種類

`server/main.go` のコメントにもありますが、よく使うのは下の2つです。

| 値 | 意味 | 使いどころ |
|---|---|---|
| `tls.NoClientCert` | クライアント証明書を要求しない | Step02 と同じ。Web の HTTPS など |
| `tls.RequireAndVerifyClientCert` | 必須+検証。**これが本当の mTLS** | サービス間通信、専用クライアント、社内ツール |

中間の `RequestClientCert` / `RequireAnyClientCert` / `VerifyClientCertIfGiven` は **「出てくれば見るが信頼はしない」「出てくるが署名を見ない」** など中途半端な動きをするので、認証目的では使わないこと。

## 観察ポイント (実験してみよう)

### 1. 証明書なしのクライアントは弾かれる

Step02 のクライアントを Step03 のサーバーに繋いでみます。

```bash
# ターミナル A: Step03 サーバーが起動中
# ターミナル B:
go run ./tutorial/step02-tls/client
```

ハンドシェイクが失敗します。サーバー側ログで:

```
handshake failed from 127.0.0.1:xxxxx: tls: client didn't provide a certificate
```

クライアント側ログで:

```
remote error: tls: certificate required
```

これが `ClientAuth = RequireAndVerifyClientCert` の効果です。

> 実験するときは、Step02 クライアントの接続先を `localhost:9444` に書き換える必要があります。

### 2. 別 CA の証明書を渡すと拒否される

別 CA を作って、その CA で署名した client cert を Step03 に置き換えてみます。

```bash
# 退避
mv tutorial/step03-mtls/certs/client.crt tutorial/step03-mtls/certs/client.crt.good
mv tutorial/step03-mtls/certs/client.key tutorial/step03-mtls/certs/client.key.good

# 別 CA + 別 client を一時ディレクトリで作成
TMP=$(mktemp -d)
openssl genrsa -out "$TMP/other-ca.key" 4096 2>/dev/null
openssl req -x509 -new -nodes -key "$TMP/other-ca.key" -sha256 -days 30 \
    -subj "/CN=Other CA" -out "$TMP/other-ca.crt"
openssl genrsa -out "$TMP/client.key" 2048 2>/dev/null
openssl req -new -key "$TMP/client.key" -subj "/CN=mallory" -out "$TMP/client.csr"
openssl x509 -req -in "$TMP/client.csr" \
    -CA "$TMP/other-ca.crt" -CAkey "$TMP/other-ca.key" -CAcreateserial \
    -days 30 -out "$TMP/client.crt"

cp "$TMP/client.crt" tutorial/step03-mtls/certs/client.crt
cp "$TMP/client.key" tutorial/step03-mtls/certs/client.key

go run ./tutorial/step03-mtls/client
# => 拒否される: tls: unknown certificate authority

# 後始末
rm -rf "$TMP"
mv tutorial/step03-mtls/certs/client.crt.good tutorial/step03-mtls/certs/client.crt
mv tutorial/step03-mtls/certs/client.key.good tutorial/step03-mtls/certs/client.key
```

サーバー側の `ClientCAs` に `Other CA` が入っていないため、署名を検証できず拒否されます。

### 3. peer の身元で認可してみる

`server/main.go` の `handle` で `peerCN` を取り出しているので、ここで認可ロジックが書けます。例:

```go
if peerCN != "alice@example.com" {
    log.Printf("unauthorized CN=%s", peerCN)
    return
}
```

これだけで「特定のクライアントだけ受け付ける」サーバーになります。**実運用では CN ではなく SAN (URI SAN や DNS SAN)** を見るのが一般的です — CN は人間向けの表示名で、運用で書き換えられがちなため。

## 用語ミニまとめ

- **mTLS (mutual TLS)**: クライアント/サーバー双方が証明書で身元を証明する TLS。
- **`Certificates` (サーバー側 / クライアント側)**: 自分が提示する証明書/鍵。同じフィールド名でも、サーバーでは「サーバー証明書」、クライアントでは「クライアント証明書」を意味する。
- **`ClientCAs`**: サーバーが、繋いでくるクライアントの証明書を検証するための CA 集合。
- **`PeerCertificates`**: ハンドシェイクで相手が提示してきた証明書チェーン。`[0]` が葉。

## ここまでで身についたもの

- `tls.Config` の主要フィールドが何のためにあるか、片側ずつではなく **両側のペアで** 説明できる。
- 「平文 TCP → サーバー認証 TLS → mTLS」と段階を踏んだことで、各ステップで **何が解決され、何が残るか** を語れる。

発展トピックは [`docs/00-overview.md`](../../docs/00-overview.md) の末尾に列挙してあります。
