// Step04: Envoy 経由で mTLS サーバーに繋ぐクライアント
//
// Step03 との違い:
//   - クライアントから TLS コードが消えた (crypto/tls / crypto/x509 不要)
//   - 接続先は Envoy のローカルリスナー :9445 で、平文 TCP で喋る
//   - mTLS のハンドシェイク (証明書提示 + 検証) は Envoy が代行する
//
// これが service mesh の sidecar 構成の最小形。アプリは TLS のことを
// 知らなくてよくなる代わりに、Envoy をプロセス境界で守る運用責任が乗る。
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const addr = "localhost:9445" // ← Envoy のリスナー (ここからは平文で OK)

func main() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer func() { _ = conn.Close() }()

	log.Printf("connected to %s (plaintext) — Envoy will originate mTLS to upstream", addr)

	go func() {
		s := bufio.NewScanner(conn)
		for s.Scan() {
			fmt.Println(s.Text())
		}
	}()

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		if _, err := fmt.Fprintln(conn, in.Text()); err != nil {
			log.Fatalf("write failed: %v", err)
		}
	}
}
