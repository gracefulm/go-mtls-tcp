// Step01: 素の TCP エコークライアント
//
// 標準入力から1行読んでサーバーへ送り、返ってきた行を標準出力に表示する。
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const addr = "localhost:9000"

func main() {
	// net.Dial は TCP 3way ハンドシェイクのみ行う。TLS は一切使われない。
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial failed: %v", err)
	}
	defer func() { _ = conn.Close() }()
	log.Printf("connected to %s (plaintext TCP)", addr)

	// 受信側を別 goroutine で回し続ける。
	go func() {
		s := bufio.NewScanner(conn)
		for s.Scan() {
			fmt.Println(s.Text())
		}
	}()

	// 標準入力からの行をそのまま送信。
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		if _, err := fmt.Fprintln(conn, in.Text()); err != nil {
			log.Fatalf("write failed: %v", err)
		}
	}
	if err := in.Err(); err != nil && err != io.EOF {
		log.Printf("stdin error: %v", err)
	}
}
