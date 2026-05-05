# Go mTLS TCP チュートリアル

Go の `net` / `crypto/tls` / `crypto/x509` パッケージだけを使い、**素の TCP → TLS（サーバー認証） → mTLS（相互認証）** と段階的にエコーサーバー/クライアントを進化させるハンズオン教材です。

## 学べること

| Step | テーマ | 解決される / 学べること |
|---|---|---|
| [Step 01](tutorial/step01-tcp-echo/README.md) | 素の TCP エコー | — |
| [Step 02](tutorial/step02-tls/README.md) | TLS（サーバー認証のみ） | 盗聴・改ざん・サーバーのなりすまし |
| [Step 03](tutorial/step03-mtls/README.md) | mTLS（相互認証） | クライアントのなりすまし |
| [Step 04](tutorial/step04-envoy-mtls/README.md) | Envoy を sidecar にして mTLS を肩代わり | アプリから TLS 設定を剥がす / Envoy の基本 (listener / cluster / transport_socket / admin) |

## プロジェクトの構成

```txt
.
├── docs/                          # 補助ドキュメント
│   ├── 00-overview.md             #   全体マップ・用語集
│   ├── 01-tcpdump.md              #   tcpdump の基本的な使い方
│   ├── 02-x509-and-ca.md          #   X.509 と CA の仕組み
│   ├── 03-openssl.md              #   openssl コマンドの基本的な使い方
│   ├── 04-certificate-flow.md     #   CA・署名・証明書の発行フロー
│   ├── 05-tls-versions.md         #   TLS 1.2 と TLS 1.3 のフロー比較
│   └── 06-envoy.md                #   Envoy 概論・sidecar 化の利点と欠点・他製品比較
├── scripts/
│   ├── gen-certs.sh               # CA・サーバー・クライアント証明書を一括生成
│   └── README.md                  # スクリプトの詳細説明
├── tutorial/
│   ├── step01-tcp-echo/           # Step 01: 素の TCP
│   │   ├── server/main.go
│   │   ├── client/main.go
│   │   └── README.md
│   ├── step02-tls/                # Step 02: TLS (サーバー認証のみ)
│   │   ├── server/main.go
│   │   ├── client/main.go
│   │   ├── certs/                 #   make certs で生成される証明書
│   │   └── README.md
│   ├── step03-mtls/               # Step 03: mTLS (相互認証)
│   │   ├── server/main.go
│   │   ├── client/main.go
│   │   ├── certs/                 #   make certs で生成される証明書
│   │   └── README.md
│   └── step04-envoy-mtls/         # Step 04: Envoy を sidecar にして mTLS を肩代わり
│       ├── client/main.go         #   平文 TCP しか喋らないクライアント (ホスト側)
│       ├── envoy/envoy.yaml       #   Envoy の静的設定
│       ├── compose.yaml           #   Envoy + Step03 サーバーを compose で同居
│       └── README.md
├── go.mod
├── Makefile
└── README.md
```

## セットアップ

```bash
# Git のコミットテンプレートを設定
make init

# Step 02 / 03 に必要な証明書を生成
make certs
```

## 進め方

1. `make certs` で証明書を生成します（Step 02 以降で必要）。
2. `tutorial/` 内の各ステップを **Step 01 → 02 → 03** の順に進めてください。
3. 各ステップの `README.md` に、動かし方・観察ポイント・実験課題が書いてあります。
4. `tcpdump` や Wireshark でパケットを観察すると、各ステップの違いをより深く理解できます。

詳細は [docs/00-overview.md](docs/00-overview.md) を参照してください。

## 補助ドキュメント

| ドキュメント | 内容 |
|---|---|
| [00-overview.md](docs/00-overview.md) | 全体マップ・用語集・発展トピック |
| [01-tcpdump.md](docs/01-tcpdump.md) | tcpdump の基本的な使い方 |
| [02-x509-and-ca.md](docs/02-x509-and-ca.md) | X.509 証明書と CA の仕組み |
| [03-openssl.md](docs/03-openssl.md) | openssl コマンドの基本的な使い方 |
| [04-certificate-flow.md](docs/04-certificate-flow.md) | CA・デジタル署名・証明書の発行フロー |
| [05-tls-versions.md](docs/05-tls-versions.md) | TLS 1.2 と 1.3 のハンドシェイクフロー比較 |
| [06-envoy.md](docs/06-envoy.md) | Envoy 概論・sidecar 化のメリット/デメリット・他製品比較 |

## Web で読む

このリポジトリの内容は GitHub Pages 上でも公開されています（`main` への push で自動デプロイ）。

- 公開URL: https://gracefulm.github.io/go-mtls-tcp/
- ローカルプレビュー: [Hugo extended](https://gohugo.io/installation/) を入れた上で `make site-serve` → http://localhost:1313/
- 静的サイトビルド: `make site-build`（`site/public/` に出力）

サイトの仕組みは `site/`（Hugo プロジェクト）と `.github/workflows/pages.yml` を参照してください。初回のみ、リポジトリ Settings → Pages → **Source: GitHub Actions** に切り替える必要があります。
