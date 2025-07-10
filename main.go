package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func handle_sending(connection net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("Sent", scanner.Text())
		connection.Write([]byte(scanner.Text() + "\r\n"))
	}
}

func main() {
	connection, err := net.Dial("tcp", "irc.dal.net:6667")

	if err != nil {
		fmt.Println(err)
		return
	}

	go handle_sending(connection)

	connection.Write([]byte("NICK Arpidanzo:\r\n"))
	connection.Write([]byte("USER Arpidanzo 0 * :twint\r\n"))
	connection.Write([]byte("JOIN #test\r\n"))

	for {
		buf := make([]byte, 4096)
		connection.Read(buf)

		fmt.Println(string(buf))
	}
}
