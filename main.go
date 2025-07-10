package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const SERVER = "cloverse.duckdns.org:6667"
const NICK = "AdriBot"
const REAL_NAME = "Adriano"

var channels []string = []string{"#test"}
var current_channel int = 0

func send(message string, connection net.Conn) {
	connection.Write([]byte(message + "\r\n"))
}

func parseMessage(sender string, args []string) {
}

func log(message string) {
	textView.SetText(textView.GetText(false) + message)
	textView.ScrollToEnd()
}

func logCommand(_ string, command string, args []string) {
	log(command + " [" + strings.Join(args, " ") + "]\n")
}

func logMessage(sender string, message string, self bool) {
	var color string
	if self {
		color = "[blue]"
	} else {
		color = "[red]"
	}
	log(color + sender + ": [white]" + message + "\n")
}

func parseCommand(message string) (string, string, []string) {

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
			prefix, command, args := parseCommand(msg)
			if command == "PRIVMSG" {
				logMessage(prefix, args[1], false)
			} else {
				logCommand(prefix, command, args)
			}
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

		logMessage("You", input.GetText(), true)
		send("PRIVMSG "+channels[current_channel]+" :"+input.GetText(), connection)

		input.SetText("")

	})

	flex.AddItem(textView, 0, 10, false)
	flex.AddItem(input, 1, 1, true)

	go handle_reading(connection)

	send("NICK "+NICK, connection)
	send("USER "+NICK+" 0 * :"+REAL_NAME, connection)
	send("JOIN "+channels[0], connection)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

}
