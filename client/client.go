package main

import (
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:4000")
	if err != nil {
		log.Fatalf("Ошибка подключения [%v]", err)
	}

	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)

}
