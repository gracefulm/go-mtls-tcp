// Step03: mTLS クライアント
//
// Step02 との違い:
//   - tls.Config.Certificates に「自分のクライアント証明書」を渡す
//     (サーバー側の Certificates と同じフィールドだが、用途は逆)
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
	addr           = "localhost:9444"
	caFile         = "tutorial/step03-mtls/certs/ca.crt"
	clientCertFile = "tutorial/step03-mtls/certs/client.crt"
	clientKeyFile  = "tutorial/step03-mtls/certs/client.key"
	serverName     = "localhost"
)

func main() {
	// 自分のクライアント証明書/秘密鍵をロード。
	// サーバーがハンドシェイク中に CertificateRequest を投げてきたとき、
	// このペアの公開鍵証明書を提示し、秘密鍵で署名を返す。
	cert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("load client keypair failed: %v", err)
	}

	// サーバー証明書検証用の CA。Step02 と同じ役割。
	caPEM, err := os.ReadFile(caFile)
	if err != nil {
		log.Fatalf("read CA failed: %v", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		log.Fatalf("CA を CertPool に追加できませんでした")
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert}, // ← Step02 にはなかった
		RootCAs:      pool,
		ServerName:   serverName,
		MinVersion:   tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		log.Fatalf("tls dial failed: %v", err)
	}
	defer func() { _ = conn.Close() }()

	state := conn.ConnectionState()
	log.Printf("connected to %s over mTLS (server CN=%s)",
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
