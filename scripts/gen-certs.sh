#!/usr/bin/env bash
#
# 教材用の証明書一式を openssl で生成する。
#   1. ルート CA (自己署名)
#   2. サーバー証明書 (CA で署名, EKU=serverAuth, SAN=localhost/127.0.0.1)
#   3. クライアント証明書 (CA で署名, EKU=clientAuth)
#
# Step02 (TLS)   は ca.crt / server.crt / server.key を使用。
# Step03 (mTLS)  は上記 + client.crt / client.key を使用。
#
# 冪等: 既存の certs/ を上書きする。CA を作り直すので、既に手元で
# 立ち上げ中のサーバーがあれば再起動が必要。

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORK_DIR="$(mktemp -d)"
trap 'rm -rf "$WORK_DIR"' EXIT

STEP02_DIR="$ROOT_DIR/tutorial/step02-tls/certs"
STEP03_DIR="$ROOT_DIR/tutorial/step03-mtls/certs"
mkdir -p "$STEP02_DIR" "$STEP03_DIR"

DAYS=825  # 一般的なブラウザの上限ルール (RFC ではない実務慣習) を踏襲

cd "$WORK_DIR"

echo ">> [1/3] ルート CA を生成中..."
openssl genrsa -out ca.key 4096 2>/dev/null
openssl req -x509 -new -nodes -key ca.key -sha256 -days "$DAYS" \
    -subj "/CN=go-mtls-tutorial Root CA" \
    -out ca.crt

cat > server.ext <<'EOF'
authorityKeyIdentifier = keyid, issuer
basicConstraints       = CA:FALSE
keyUsage               = digitalSignature, keyEncipherment
extendedKeyUsage       = serverAuth
subjectAltName         = @alt_names

[alt_names]
DNS.1 = localhost
IP.1  = 127.0.0.1
EOF

echo ">> [2/3] サーバー証明書を生成中 (CN=localhost, SAN=localhost/127.0.0.1)..."
openssl genrsa -out server.key 2048 2>/dev/null
openssl req -new -key server.key -subj "/CN=localhost" -out server.csr
openssl x509 -req -in server.csr \
    -CA ca.crt -CAkey ca.key -CAcreateserial \
    -days "$DAYS" -sha256 \
    -extfile server.ext \
    -out server.crt

cat > client.ext <<'EOF'
authorityKeyIdentifier = keyid, issuer
basicConstraints       = CA:FALSE
keyUsage               = digitalSignature
extendedKeyUsage       = clientAuth
subjectAltName         = @alt_names

[alt_names]
DNS.1 = client.local
EOF

echo ">> [3/3] クライアント証明書を生成中 (CN=alice@example.com)..."
openssl genrsa -out client.key 2048 2>/dev/null
openssl req -new -key client.key -subj "/CN=alice@example.com" -out client.csr
openssl x509 -req -in client.csr \
    -CA ca.crt -CAkey ca.key -CAcreateserial \
    -days "$DAYS" -sha256 \
    -extfile client.ext \
    -out client.crt

# Step02 へ配布 (サーバー認証のみなのでクライアント鍵は不要)
install -m 0644 ca.crt     "$STEP02_DIR/ca.crt"
install -m 0644 server.crt "$STEP02_DIR/server.crt"
install -m 0600 server.key "$STEP02_DIR/server.key"

# Step03 へ配布 (mTLS なのでサーバー/クライアント双方が CA を検証に使う)
install -m 0644 ca.crt     "$STEP03_DIR/ca.crt"
install -m 0644 server.crt "$STEP03_DIR/server.crt"
install -m 0600 server.key "$STEP03_DIR/server.key"
install -m 0644 client.crt "$STEP03_DIR/client.crt"
install -m 0600 client.key "$STEP03_DIR/client.key"

echo
echo "OK: 証明書の配布が完了しました。"
echo "  Step02: $STEP02_DIR"
echo "  Step03: $STEP03_DIR"
