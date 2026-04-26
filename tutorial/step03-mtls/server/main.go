// Step03: mTLS (相互認証) のエコーサーバー
//
// Step02 との違い:
//   - ClientCAs にクライアント証明書検証用の CA を積む
//   - ClientAuth = RequireAndVerifyClientCert で「証明書必須+検証必須」にする
//   - 接続後、PeerCertificates から相手の CN を取り出してログに出す
package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	addr     = ":9444"
	certFile = "tutorial/step03-mtls/certs/server.crt"
	keyFile  = "tutorial/step03-mtls/certs/server.key"
	caFile   = "tutorial/step03-mtls/certs/ca.crt"
)

func main() {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("load server keypair failed: %v", err)
	}

	// クライアント証明書を検証するための CA をロードする。
	// (Step02 ではクライアントだけが CA を持っていた。mTLS では双方が持つ)
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("read CA failed: %v", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		log.Fatalf("CA を CertPool に追加できませんでした")
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		// クライアント証明書検証用の CA。
		ClientCAs: pool,
		// クライアントは証明書を「必ず」提示し、ClientCAs で検証可能でなければ拒否。
		// 他の選択肢:
		//   tls.NoClientCert              ... 要求しない (= Step02 と同じ)
		//   tls.RequestClientCert         ... 出してもいいが検証はしない
		//   tls.RequireAnyClientCert      ... 必ず出させるが署名検証はしない
		//   tls.VerifyClientCertIfGiven   ... 出してきたら検証する (任意)
		//   tls.RequireAndVerifyClientCert ... 出させて かつ 検証する  ← これ
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}

	ln, err := tls.Listen("tcp", addr, cfg)
	if err != nil {
		log.Fatalf("tls listen failed: %v", err)
	}
	defer func() { _ = ln.Close() }()

	log.Printf("step03 mTLS server listening on %s (mutual auth)", addr)

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

	// tls.Listen 経由なので *tls.Conn が返ってくる。
	// 明示的に Handshake() を呼ぶと、その時点で失敗を捕まえられる
	// (デフォルトでは最初の Read/Write でハンドシェイクが行われる)。
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		log.Printf("not a tls connection")
		return
	}
	if err := tlsConn.Handshake(); err != nil {
		log.Printf("handshake failed from %s: %v", conn.RemoteAddr(), err)
		return
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		// ClientAuth=RequireAndVerifyClientCert を設定していれば
		// ハンドシェイクの段階で弾かれているため、本来ここには来ない。
		log.Printf("no peer cert (unexpected)")
		return
	}
	peerCN := state.PeerCertificates[0].Subject.CommonName
	log.Printf("connected: %s (peer CN=%s)", conn.RemoteAddr(), peerCN)

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("recv from CN=%s: %q", peerCN, line)
		// 認証された相手の名前をエコーに混ぜる = サーバーが「誰か」を知っている証拠。
		if _, err := fmt.Fprintf(conn, "hello %s, you said: %s\n", peerCN, line); err != nil {
			return
		}
	}
	log.Printf("disconnected: %s (CN=%s)", conn.RemoteAddr(), peerCN)
}
