package main

import (
	"fmt"
	"net"
)

func main() {
	connection, err := net.Dial("tcp", "irc.dal.net:6667")

	if err != nil {
		fmt.Println(err)
		return
	}

	connection.Write([]byte("NICK Arpidanzo:\r\n"))
	connection.Write([]byte("USER Arpidanzo 0 * :twint\r\n"))
	connection.Write([]byte("JOIN #test\r\n"))

	for {
		buf := make([]byte, 4096)
		connection.Read(buf)

		fmt.Println(string(buf))
	}
}
