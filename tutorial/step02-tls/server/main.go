// Step02: TLS (サーバー認証のみ) のエコーサーバー
//
// Step01 との違い:
//   - net.Listen ではなく tls.Listen を使う
//   - 自分の証明書/秘密鍵を tls.Config.Certificates に渡す
//   - クライアントは何も提示しない (= サーバー認証のみ)
package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

const (
	addr     = ":9443"
	certFile = "tutorial/step02-tls/certs/server.crt"
	keyFile  = "tutorial/step02-tls/certs/server.key"
)

func main() {
	// 自サーバーの証明書 (公開鍵) と秘密鍵をペアでロードする。
	// 証明書は scripts/gen-certs.sh が生成したもの。
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("load server keypair failed: %v", err)
	}

	cfg := &tls.Config{
		// クライアントとのハンドシェイク時にこの証明書を提示する。
		Certificates: []tls.Certificate{cert},
		// TLS 1.2 未満は脆弱性が多いので明示的に下限を引き上げる。
		MinVersion: tls.VersionTLS12,
	}

	// tls.Listen は内部で net.Listen("tcp", ...) を呼び、
	// 各 Accept で TLS ハンドシェイク済みの接続を返してくれる。
	ln, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		log.Fatalf("tls listen failed: %v", err)
	}
	defer func() { _ = ln.Close() }()

	log.Printf("step02 TLS server listening on %s (server auth only)", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept failed: %v", err)
			continue
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	log.Printf("connected: %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("recv from %s: %q", conn.RemoteAddr(), line)
		if _, err := fmt.Fprintf(conn, "echo: %s\n", line); err != nil {
			log.Printf("write failed: %v", err)
			return
		}
	}
	log.Printf("disconnected: %s", conn.RemoteAddr())
}
