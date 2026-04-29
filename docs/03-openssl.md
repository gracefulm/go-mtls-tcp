# openssl — 基本的な使い方

本教材で扱う TLS 通信や証明書（X.509）の仕組みを理解・デバッグする上で欠かせないコマンドラインツール `openssl` の代表的な使い方をまとめます。

## 何ができるツールか

- 暗号化・復号、ハッシュ計算
- 秘密鍵や公開鍵の生成
- 証明書署名要求 (CSR) の作成
- 自己署名証明書 (オレオレ証明書) の発行
- 作成した証明書や CSR の中身の確認
- TLS サーバーへのテスト接続 (`s_client`)

---

## 1. 証明書・鍵の中身を確認する

拡張子が `.crt` や `.pem` のファイルは多くの場合 Base64 エンコードされているため、そのまま `cat` で見ても読めません。`openssl` を使ってパースします。

### 証明書 (X.509) の中身を確認する
```bash
openssl x509 -in server.crt -text -noout
```
- `-text`: 発行者 (Issuer)、所有者 (Subject)、有効期限、SAN (Subject Alternative Name) などをテキストで出力します。
- `-noout`: `-----BEGIN CERTIFICATE-----` のようなエンコードされた本体を出力しません。

### 証明書署名要求 (CSR) の中身を確認する
```bash
openssl req -in server.csr -text -noout
```
- CSR に記載されているドメイン名 (CN) や SAN が正しく設定されているか確認する際によく使います。

### 秘密鍵の内容・整合性を確認する
```bash
openssl rsa -in server.key -check -noout
```

---

## 2. 鍵と証明書を作る

### 秘密鍵を作る (RSA)
```bash
openssl genrsa -out server.key 2048
```
- 2048ビットのRSA秘密鍵を `server.key` として生成します。

### CSR (証明書署名要求) を作る
```bash
openssl req -new -key server.key -out server.csr
```
- この段階で Country Name や Common Name (CN) などを対話式で聞かれます。

### オレオレ証明書 (自己署名証明書) を手っ取り早く作る
```bash
openssl req -x509 -new -nodes -key server.key -sha256 -days 365 -out server.crt
```
- CA を立てずに、単一のサーバー証明書だけをテスト用に作りたい場合に便利です。

---

## 3. TLS サーバーにテスト接続する (`s_client`)

ブラウザや Go のクライアントコードを書く前に、サーバーが正しく TLS で待ち受けているか、どんな証明書を返してくるかを確認できる強力なデバッグツールです。

### 基本的な接続テスト
```bash
openssl s_client -connect localhost:9443
```
- サーバーが提示した証明書チェーン（`Certificate chain`）や、選択された暗号化方式（`Cipher`）などがズラッと表示されます。
- 入力待ち状態になるので、`GET / HTTP/1.0` と打って Enter を2回押すと、HTTP リクエストを手動で送ることもできます。

### SNI (Server Name Indication) を指定して接続
```bash
openssl s_client -connect 127.0.0.1:9443 -servername localhost
```
- 1つのIPアドレスで複数のドメインを運用しているサーバー (バーチャルホスト) に対して、どのドメイン宛の通信かを明示的に伝えるオプションです。

### 独自の CA 証明書を指定して検証する
```bash
openssl s_client -connect localhost:9443 -CAfile ca.crt
```
- OS にインストールされていない「オレオレ CA」で署名したサーバーに繋ぐ場合、`-CAfile` を指定しないと `Verify return code: 21 (unable to verify the first certificate)` のようなエラーになります。このオプションを指定して検証が成功するか確認できます。
