# scripts/

教材で使う雑多なシェルスクリプト置き場です。

## `gen-certs.sh`

教材 Step02 / Step03 で使う TLS 証明書を openssl で生成し、各ステップの `certs/` に配布します。

### 使い方

リポジトリのルートから実行してください (`Makefile` 経由が推奨)。

```bash
make certs
# あるいは
./scripts/gen-certs.sh
```

### 生成されるもの

スクリプトは一時ディレクトリで作業し、最終的に以下を配置します。

| 配置先 | ファイル | 役割 |
|---|---|---|
| `tutorial/step02-tls/certs/` | `ca.crt` | クライアントがサーバー証明書を検証するための CA 証明書 |
| | `server.crt` | サーバーが提示する公開鍵証明書 (CA で署名済み) |
| | `server.key` | サーバーの秘密鍵 (mode 600) |
| `tutorial/step03-mtls/certs/` | 上記すべて + | (server.key は mTLS でも必要) |
| | `client.crt` | クライアントが提示する公開鍵証明書 |
| | `client.key` | クライアントの秘密鍵 (mode 600) |

### 主な X.509 拡張

- **Server cert**
  - `extendedKeyUsage = serverAuth`
  - `subjectAltName = DNS:localhost, IP:127.0.0.1`
    Go 1.15 以降、`tls.Config.ServerName` は SAN と照合されます (CN は使われません)。`localhost` で接続するために必須です。
- **Client cert**
  - `extendedKeyUsage = clientAuth`
  - `subjectAltName = DNS:client.local`
    クライアント認証では SAN 検証は通常スキップされますが、CN/SAN は認可ロジック (= 「どの CN を許可するか」) で参照することがあるので明示しています。
- いずれも `basicConstraints = CA:FALSE`。CA としては振る舞えません (= この証明書から別の証明書に署名できません)。

### 注意

- 既存の `certs/` を**上書きします**。CA 鍵が変わるため、起動中のサーバー/クライアントは再起動してください。
- 教材以外の用途には使わないでください。鍵長 2048bit、有効期間 825 日、無パスフレーズなど、学習用に簡略化されています。
