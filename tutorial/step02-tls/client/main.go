// Step02: TLS クライアント (サーバー認証のみ)
//
// Step01 との違い:
//   - net.Dial ではなく tls.Dial を使う
//   - サーバーが提示する証明書を検証するための CA を RootCAs に積む
//   - 接続先のホスト名が証明書の SAN と一致するか ServerName で照合する
package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
)

const (
	addr   = "localhost:9443"
	caFile = "tutorial/step02-tls/certs/ca.crt"
	// SAN と照合される値。証明書の subjectAltName=DNS:localhost と一致させる。
	serverName = "localhost"
)

func main() {
	// 「このサーバー証明書を信頼する」ための CA をロードする。
	// プロセスが信頼する CA の集合 = CertPool。
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("read CA failed: %v", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		log.Fatalf("CA を CertPool に追加できませんでした (PEM が壊れている?)")
	}

	cfg := &tls.Config{
		// サーバー証明書は「この CA で署名されている」必要がある。
		RootCAs: pool,
		// 接続先ホスト名 (= 証明書の SAN と照合される)。
		// SNI として TLS ハンドシェイクでも送られる。
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		log.Fatalf("tls dial failed: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// ハンドシェイク完了後、相手が提示してきた証明書チェーンを覗き見る。
	state := conn.ConnectionState()
	log.Printf("connected to %s over TLS (server CN=%s)",
		addr, state.PeerCertificates[0].Subject.CommonName)

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
