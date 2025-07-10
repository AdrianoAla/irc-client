package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const SERVER = "cloverse.duckdns.org:6667"

func handle_sending(connection net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		fmt.Println("Sent", scanner.Text())
		connection.Write([]byte(scanner.Text() + "\r\n"))
	}
}

func parseMessage(message string) (string, string, []string) {

	var prefix string
	var command string

	if message[0] == ':' {
		split := strings.Split(message[1:], " ")
		prefix = split[0]
		message = strings.Join(split[1:], " ")
	}

	split := strings.Split(message, " :")

	args := strings.Split(split[0], " ")
	command, args = args[0], args[1:]

	if len(split) > 1 {
		args = append(args, split[1])
	}

	return prefix, command, args

}

func main() {
	connection, err := net.Dial("tcp", SERVER)

	if err != nil {
		fmt.Println(err)
		return
	}

	go handle_sending(connection)

	connection.Write([]byte("NICK BotLol\r\n"))
	connection.Write([]byte("USER BotLol 0 test :Arpidanzo \r\n"))
	connection.Write([]byte("JOIN #test\r\n"))

	for {
		buf := make([]byte, 4096)
		connection.Read(buf)
		msgs := strings.Split(string(buf), "\r\n")
		for _, msg := range msgs {
			if []byte(msg)[0] == 0x00 {
				continue
			}
			prefix, command, args := parseMessage(msg)
			fmt.Println(prefix, command, args)
		}
	}
}
