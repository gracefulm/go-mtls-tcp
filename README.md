# Go Template Project

> [!IMPORTANT]
>
> - このプロジェクトは、[project-layout](https://github.com/golang-standards/project-layout/blob/master/README_ja.md)を参考にしています。
>   - project-layoutはGoの公式のプロジェクトのレイアウトではないことに注意してください。
>   - Goの公式のプロジェクト構成は[こちら](https://go.dev/doc/modules/layout)を参考にしてください。
> - 実際に使用する時は必要なディレクトリのみ残してスモールスタートで開発してください。
>   - 特に、簡単なコマンドのみで構成されるプロジェクトではフラットにして開発するのが良いです。

## プロジェクトの構成

```txt
.
├── api/                 # APIの定義(Swagger/Protocol Buffers etc...)
├── assets/              # アセット(画像/音声 etc...)
├── build/               # ビルド/テスト/デプロイのためのツール
├── cmd/                 # ソースコードのエントリポイント
├── configs/             # 設定ファイル
├── deployments/         # IaaS、PaaS、システム、コンテナオーケストレーションのデプロイメント設定とテンプレート
├── docs/                # デザインドキュメントとユーザードキュメント
├── examples/            # サンプルコード
├── githooks/            # Gitのフックスクリプト
├── go.mod               # Goのモジュール定義
├── init/                # システムinit(systemd, upstart, sysv)とプロセスマネージャ/スーパーバイザ(runit, supervisord)の設定
├── internal/            # 内部のパッケージ
├── LICENSE              # ライセンス
├── Makefile             # Makefile
├── pkg/                 # 公開パッケージ
├── README.md            # プロジェクトの説明
├── scripts/             # スクリプト類
├── test/                # 追加の外部テストアプリとテストデータ
├── third_party/         # 外部ヘルパーツール
├── tools/               # このプロジェクトをサポートするツール類
├── vender/              # アプリケーションの依存関係(go mod vendorコマンドで生成される、基本的に使わない)
├── web/                 # ウェブアプリケーション固有のコンポーネント(SPAなど)
└── website/             # プロジェクトのウェブサイト置き場
```

## Setup

```bash
make init
```
