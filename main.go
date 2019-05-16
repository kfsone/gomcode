package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

type Context struct {
	UseTUI  bool
	Timeout time.Duration
	User    UserInterfacer
	Remote  chan string
	Cmd     string
	Argv    []string
}

type CommandFn func(ctx Context) error

var commands = map[string]CommandFn{
	"q":    cmd_quit,
	"quit": cmd_quit,
	"exit": cmd_quit,
	"help": cmd_help,
}

func sendRaw(ctx Context, raw string) {
	raw = strings.TrimSpace(raw)
	timeout := time.After(ctx.Timeout)
	for {
		select {
		case ctx.Remote <- raw:
			{
				return
			}
		case <-timeout:
			{
				ctx.User.Error(fmt.Sprintf("Unable to send command within %v", ctx.Timeout))
				return
			}
		}
	}
}

func parse(ctx Context, cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}
	ctx.User.WriteString("> " + cmd)
	cmd = strings.Split(cmd, "#")[0]
	if cmd == "" {
		return
	}
	if cmd[0] == '"' || cmd[0] == '\'' {
		sendRaw(ctx, cmd[1:])
		return
	}
	argv := strings.Fields(cmd)
	cmdfn, ok := commands[argv[0]]
	if !ok {
		ctx.User.Error("No such command: " + argv[0])
		return
	}
	ctx.Cmd, ctx.Argv = argv[0], argv[1:]
	if err := cmdfn(ctx); err != nil {
		ctx.User.Error(fmt.Sprintf("'%s': %s", argv[0], err))
	}
}

func main() {
	var ctx Context

	flag.BoolVar(&ctx.UseTUI, "tui", false, "Use TUI for the user interface")
	flag.DurationVar(&ctx.Timeout, "timeout", 60*time.Second, "Timeout sending remote commands (seconds)")

	flag.Parse()

	var ui UserInterfacer
	if ctx.UseTUI {
		log.Println("Launching TUI")
		ui = NewTUIUserInterface()
	} else {
		log.Println("Launching RawUI")
		ui = NewRawUserInterface()
	}

	ui.Start()

	ctx.User = ui
	ctx.User.WriteString("-- Starting")

	ctx.Remote = make(chan string, 4)

	ui.WriteString("-- Ready")
	for cmd := range ui.Commands() {
		parse(ctx, cmd)
	}
}
