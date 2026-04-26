# tcpdump — 基本的な使い方

本教材では、各ステップで「実際にネットワーク上を流れているバイト列」を観察するために `tcpdump` を使うと理解が深まります。ここでは、本リポジトリの `tutorial/` を題材にした最低限の使い方をまとめます。

## 何ができるツールか

ネットワークインターフェースを通過するパケットをキャプチャ・表示するCLIツール。生のキャプチャには通常 `sudo` が必要です。

- 平文プロトコル (Step 01) の中身がそのまま見える
- TLSハンドシェイク (Step 02 / 03) の流れと、その後のApplication Dataが暗号化されて見えなくなる様子を比較できる
- `-w` で `.pcap` に保存して [Wireshark](https://www.wireshark.org/) で詳細解析できる

## 基本コマンド

```bash
# インターフェース一覧
sudo tcpdump -D

# 特定インターフェースをキャプチャ (lo0 = ループバック / en0 = 有線・Wi-Fi)
sudo tcpdump -i lo0

# ポート指定
sudo tcpdump -i lo0 port 8443

# ホスト指定
sudo tcpdump -i en0 host 192.168.1.10

# AND/OR の組み合わせ
sudo tcpdump -i lo0 'tcp and port 8443'
sudo tcpdump -i en0 'host 1.2.3.4 and (port 80 or port 443)'
```

## よく使うオプション

| オプション | 意味 |
|---|---|
| `-i <iface>` | キャプチャするインターフェース |
| `-n` | 名前解決しない (DNSで止まらないので速い) |
| `-nn` | ポートも数値で表示 |
| `-v` / `-vv` / `-vvv` | 表示の詳細度を上げる |
| `-X` | パケット内容を hex + ASCII で表示 |
| `-A` | ASCII のみ表示 (HTTPなど平文の確認に便利) |
| `-s 0` | 全長キャプチャ (古いtcpdumpで必要、新しい版はデフォルトで全長) |
| `-c <N>` | N個キャプチャしたら終了 |
| `-w file.pcap` | ファイルに保存 (後で Wireshark で開ける) |
| `-r file.pcap` | 保存ファイルから読み込み |

## このリポジトリでの使い方

`tutorial/` の各ステップは `localhost:8443` で動くため、ループバック (`lo0`) を見ます。
**3つ目のターミナル**を開いて、サーバー → tcpdump → クライアント の順に起動するのが見やすいです。

### Step 01 — 平文TCPの中身を覗く

```bash
sudo tcpdump -i lo0 -A 'tcp port 8443'
```

`Hello, mTLS Tutorial!` のような送受信文字列がそのままASCIIで見えます。
**「TLSなしだと中身が丸見え」**という感覚を掴むのに最適です。

### Step 02 / 03 — TLSハンドシェイク後に暗号化される様子を見る

```bash
sudo tcpdump -i lo0 -X 'tcp port 8443'
```

接続直後に ClientHello / ServerHello / Certificate / ... のハンドシェイクが流れたあと、Application Data 以降は暗号化されて意味のあるASCIIが出なくなります。

### pcap に保存して Wireshark で開く

```bash
sudo tcpdump -i lo0 -w /tmp/mtls.pcap 'tcp port 8443'
```

Wireshark でTLSフィルタ (`tls`) をかけると、ハンドシェイクの各メッセージが構造化されて読めます。Step 03 では `Certificate` メッセージがクライアント・サーバー両方から飛ぶことが確認でき、mTLSの「相互」性を視覚的に理解できます。

## 出力の読み方

```
12:34:56.789 IP 127.0.0.1.54321 > 127.0.0.1.8443: Flags [S], seq 0, win 65535
```

`時刻 IP 送信元.ポート > 宛先.ポート: Flags[...]` の形式。フラグ:

| 表記 | 意味 |
|---|---|
| `S` | SYN |
| `.` | ACK |
| `P` | PSH |
| `F` | FIN |
| `R` | RST |

3-wayハンドシェイクは `[S]` → `[S.]` → `[.]` の3行で確認できます (Step 01 README の3wayハンドシェイク解説と対応)。

## よく使うフィルタ式

```bash
'tcp port 443'                  # TCP 443番のみ
'udp port 53'                   # DNS
'host 1.2.3.4'                  # 送信元・宛先どちらかが一致
'src host 1.2.3.4'              # 送信元のみ
'dst port 8443'                 # 宛先ポートのみ
'tcp[tcpflags] & tcp-syn != 0'  # SYNパケットのみ
'not arp'                       # ARPを除外
```

## 注意点

- macOSでループバックを見るときは `-i lo0` (Linuxは `-i lo`)。
- `sudo` なしでは BPF デバイスが開けず Permission denied になります。
- ループバックはMTUが大きいため、1パケットに乗るデータ量がイーサネットより多く見えます。
- TLSの中身を「復号して」見たい場合は、`SSLKEYLOGFILE` などでキーをエクスポートしてWiresharkに食わせる必要があります (本教材の範囲外)。
