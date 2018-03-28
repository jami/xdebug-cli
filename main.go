package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func handleXDebugConnection(c net.Conn) {
	defer c.Close()

	message := ""
	for {
		tmp := make([]byte, 256)
		bytesRead, err := c.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		fmt.Println("got", bytesRead, "bytes.")
		message = message + string(tmp[:bytesRead])

		if bytesRead < len(tmp) {
			fmt.Println("message \n", message)
			fmt.Print("\n Enter text: ")
			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = text + "\n"
			c.Write([]byte(text))
			c.Write([]byte{0})
		}
	}

	fmt.Println("connection closed")
}

func main() {
	fmt.Println("xdebug-cli")
	ln, err := net.Listen("tcp", "127.0.0.1:9998")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		fmt.Println("accept connection")
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleXDebugConnection(conn)
	}
}
