# Go mTLS TCP ハンズオン — 全体マップ

## 何を学ぶか

Go の `net` / `crypto/tls` / `crypto/x509` パッケージだけを使って、

```
素の TCP   →   TLS (サーバー認証のみ)   →   mTLS (相互認証)
```

と段階的にエコーサーバー/クライアントを進化させ、**各層で「何が解決され、何が残るか」を自分の言葉で説明できる**状態を目指します。

## 前提知識

- Go の基本構文と `go run` / `go build` の動かし方。
- ターミナルを2つ開いて作業できること (サーバー1つ、クライアント1つ)。
- TCP の名前を聞いたことがある程度で十分です。TLS / X.509 はゼロから扱います。

## 進め方

リポジトリのルートで作業してください。各ステップは独立したディレクトリに分かれており、それぞれ `README.md` に**動かし方・観察ポイント・実験課題**が書いてあります。

1. **証明書の準備 (Step02 以降に必要)**

    ```bash
    make certs
    ```

    `tutorial/step02-tls/certs/` と `tutorial/step03-mtls/certs/` に必要な `.crt`/`.key` が配布されます。詳細は [`scripts/README.md`](../scripts/README.md)。

2. **ステップを順に実施**
    - [Step 01 — 素の TCP エコー](../tutorial/step01-tcp-echo/README.md)
    - [Step 02 — TLS (サーバー認証のみ)](../tutorial/step02-tls/README.md)
    - [Step 03 — mTLS (相互認証)](../tutorial/step03-mtls/README.md)
    - [Step 04 — Envoy を sidecar にして mTLS を肩代わり](../tutorial/step04-envoy-mtls/README.md)

3. 詰まったら本ドキュメントの **用語集** に戻ってきてください。

## 補助資料

- [tcpdump — 基本的な使い方](01-tcpdump.md): 各ステップで実際に流れるバイト列を観察するための入門。Step 01 で平文が丸見えになる様子、Step 02 / 03 でTLS化後に暗号化される様子を比較できます。
- [Envoy 概論 — sidecar 化のメリット・デメリットと他製品比較](06-envoy.md): Step 04 の周辺知識。なぜ Envoy なのか、sidecar 化で何が嬉しくて何が痛いか、Istio / Linkerd / nginx / Traefik など他製品との位置づけ。

## 各ステップの早見表

| Step | 目的 | 増えるもの | 解決される / 学べること |
|---|---|---|---|
| 01 | TCP の素の挙動 | (TLS なし) | — |
| 02 | サーバー認証 TLS | `tls.Config`: `Certificates` (サーバー側), `RootCAs` / `ServerName` (クライアント側) | 盗聴・改ざん・サーバーのなりすまし |
| 03 | mTLS | `tls.Config`: `ClientCAs` / `ClientAuth` (サーバー側), `Certificates` (クライアント側) | クライアントのなりすまし |
| 04 | Envoy sidecar で mTLS を肩代わり | Envoy YAML: `listener` / `cluster` / `transport_socket` (`UpstreamTlsContext`) / `admin` | アプリから TLS 設定を剥がす / Envoy の基本構成 |

## 用語集 (アルファベット順)

- **CA (認証局)**: 証明書に署名する権威。商用 CA (Let's Encrypt 等) もあれば、社内用にプライベート CA を立てることもある。本教材では `gen-certs.sh` で自分用ルート CA を作る。
- **証明書チェーン**: 葉の証明書 → 中間 CA の証明書 → ルート CA の証明書、と署名がつながった鎖。検証側はルートを「信頼アンカー」として持っており、そこに辿り着けば OK。
- **`ClientAuth`**: サーバーが「クライアント証明書をどう扱うか」を決める列挙体。本物の mTLS には `RequireAndVerifyClientCert` を使う。
- **`ClientCAs`**: サーバーが、繋いでくるクライアントの証明書を検証するための CA プール。
- **EKU (Extended Key Usage)**: 「この証明書は何用か」のフラグ。`serverAuth` (TLS サーバー), `clientAuth` (TLS クライアント) など。`gen-certs.sh` で server/client 用に分けている。
- **ハンドシェイク**: TLS セッション開始時に、双方が証明書を見せ合い、暗号アルゴリズムを合意し、共通鍵を確立するフェーズ。
- **mTLS (mutual TLS)**: 双方向の TLS。クライアントとサーバーの両方が証明書を提示する。
- **PEM**: `-----BEGIN CERTIFICATE-----` で始まるテキスト形式の証明書/鍵。Base64 + ヘッダー。中身は DER。
- **`PeerCertificates`**: ハンドシェイクで相手から受け取った証明書チェーン (`[]*x509.Certificate`)。`[0]` が葉。
- **`RootCAs`**: クライアントが、繋ぐ先のサーバー証明書を検証するための CA プール。
- **SAN (Subject Alternative Name)**: 「この証明書はどのホスト名/IP/URI を名乗ってよいか」のリスト。Go 1.15 以降、TLS のホスト名検証は SAN のみを見る。
- **`ServerName` (`tls.Config`)**: クライアントが「自分はこのホスト名に繋ぎたい」と宣言する値。SNI として送られ、サーバー証明書の SAN と照合される。
- **SNI (Server Name Indication)**: 1 IP で複数ホストを捌くため、ClientHello に乗ってサーバーへ伝わるホスト名。
- **X.509**: 証明書の標準フォーマット (RFC 5280)。PKI でやり取りする「証明書」と言ったらだいたいこれ。

## ここから先 (発展トピック)

本教材ではあえて触れていない、現実の mTLS 運用で出てくるテーマです。気になったら自分で素振りしてみてください。

1. **証明書ローテーション**
    `tls.Config.GetCertificate` (サーバー) / `GetClientCertificate` (クライアント) のコールバックを使って、プロセス再起動なしに証明書を入れ替える。
2. **失効 (revocation)**
    CRL / OCSP / OCSP Stapling。漏洩した証明書を「もう信頼するな」と伝える仕組み。
3. **SPIFFE / SPIFFE ID**
    SAN に `URI:spiffe://example.org/ns/foo/sa/bar` のような ID を入れて、ワークロード ID として認可に使う流派。
4. **Graceful shutdown / コンテキスト連携**
    `ln.Close()` だけでは進行中の接続が荒く切れる。`context.Context` を `Accept` ループと handler に通して、SIGTERM 等で綺麗に終わらせる作法。
5. **アプリ層プロトコルの設計**
    本教材は行ベースのテキストエコーで簡単にしたが、実運用では length-prefixed バイナリや gRPC over mTLS など、もう一段の枠組みが入る。
6. **Go の `crypto/x509` で証明書を生成する**
    今回 openssl で行った CA/証明書生成を、Go コードだけで再現すると `x509.Certificate` テンプレートの構造がより深く分かる。
7. **`tls.Config.VerifyPeerCertificate`**
    標準の検証に加えて、独自の追加検証 (例えば SAN の URI を見て認可) を差し込むフック。
