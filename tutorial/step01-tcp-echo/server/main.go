// Step01: 素の TCP エコーサーバー
//
// このステップでは TLS を使わず、平文の TCP で行ベースのエコーを返す。
// Step02 (TLS) / Step03 (mTLS) との差分を体感するためのベースライン。
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const addr = ":9000"

func main() {
	// net.Listen は IPv4/IPv6 両方の "any" にバインドする TCP リスナーを返す。
	// この時点では暗号化も認証もない。誰でも接続できる。
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen failed: %v", err)
	}
	defer func() { _ = ln.Close() }()

	log.Printf("step01 echo server listening on %s (plaintext TCP)", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept failed: %v", err)
			continue
		}
		// 1接続 1 goroutine の素朴なモデル。教材なので意図的にシンプル。
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	log.Printf("connected: %s", conn.RemoteAddr())

	// 行単位でテキストを受け取り、そのまま返す。
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("recv from %s: %q", conn.RemoteAddr(), line)
		if _, err := fmt.Fprintf(conn, "echo: %s\n", line); err != nil {
			log.Printf("write failed: %v", err)
			return
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("scan error: %v", err)
	}
	log.Printf("disconnected: %s", conn.RemoteAddr())
}
