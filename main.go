package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const SERVER = "cloverse.duckdns.org:6667"

func handle_sending(message string, connection net.Conn) {
	connection.Write([]byte(message + "\r\n"))
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

func handle_reading(connection net.Conn) {
	for {
		buf := make([]byte, 4096)
		connection.Read(buf)
		msgs := strings.Split(string(buf), "\r\n")
		for _, msg := range msgs {
			if []byte(msg)[0] == 0x00 {
				continue
			}
			prefix, command, args := parseMessage(msg)
			textView.SetText(textView.GetText(false) + prefix + " " + command + " [" + strings.Join(args, " ") + "]\n")
			//fmt.Println(prefix, command, args)

		}
	}
}

var app = tview.NewApplication()

var textView = tview.NewTextView().
	SetDynamicColors(true).
	SetRegions(true).
	SetChangedFunc(func() {
		app.Draw()
	})

var flex = tview.NewFlex().SetDirection(tview.FlexRow)

func main() {
	connection, err := net.Dial("tcp", SERVER)

	if err != nil {
		fmt.Println(err)
		return
	}

	input := tview.NewInputField().SetFieldBackgroundColor(tcell.Color102)

	input.SetDoneFunc(func(key tcell.Key) {

		if key != tcell.KeyEnter {
			return
		}

		textView.SetText(textView.GetText(false) + input.GetText() + "\n")
		handle_sending(input.GetText(), connection)

		input.SetText("")

	})

	flex.AddItem(textView, 0, 10, false)
	flex.AddItem(input, 1, 1, true)

	go handle_reading(connection)

	connection.Write([]byte("NICK BotLol\r\n"))
	connection.Write([]byte("USER BotLol 0 test :Arpidanzo \r\n"))
	connection.Write([]byte("JOIN #test\r\n"))

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

}
