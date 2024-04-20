package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type ftpConn struct {
	c net.Conn
}

func (conn ftpConn) list(args []string) {
	var filename string
	switch len(args) {
	case 0:
		filename = "."
	case 1:
		filename = args[0]
	}
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(conn.c, "Ошибка при открытии файла [%v]\n", err)
	}
	stat, err := f.Stat()
	if err != nil {
		fmt.Fprintf(conn.c, "Stat error [%v]", err)
	}

	if stat.IsDir() {
		filenames, err := f.Readdirnames(0)
		if err != nil {
			fmt.Fprintf(conn.c, "Ошибка прочтения директории [%v]\n", err)
			return
		}
		for _, f := range filenames {
			fmt.Fprint(conn.c, f, "\n")
		}
	} else {
		fmt.Fprint(conn.c, filename, "\n")
	}

}

func (conn ftpConn) get(args []string) {
	if len(args) != 1 {
		fmt.Fprint(conn.c, "get [имя файла]")
	}
	filename := args[0]
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprint(conn.c, "Файл не найден\n")
	}
	io.Copy(conn.c, f)
}

func handleConn(conn ftpConn) {
	defer conn.c.Close()
	input := bufio.NewScanner(conn.c)
	var (
		cmd  string
		args []string
	)

	for input.Scan() {

		fields := strings.Fields(input.Text())
		if len(fields) == 0 {
			continue
		}
		cmd = strings.ToLower(fields[0])
		if len(fields) > 1 {
			args = fields[1:]
		}

		switch cmd {
		case "ls":
			conn.list(args)
		case "get":
			conn.get(args)
		case "close":
			return
		default:
			fmt.Fprint(conn.c, "Команда не определена")
		}

	}

}

func main() {
	listener, err := net.Listen("tcp", "localhost:4000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		var conn ftpConn
		conn.c, err = listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
