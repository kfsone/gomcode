package main

import (
	"log"

	tui "github.com/marcusolsson/tui-go"
)

type TUIUserInterface struct {
	UserInterface
	Tui        tui.UI
	output     chan string
	entry      *tui.Entry
	scrollback *tui.Box
}

func (u *TUIUserInterface) Close() {
	u.UserInterface.Close()
	u.Tui.Quit()
}

func (u *TUIUserInterface) Start() {
	u.UserInterface.Start()

	go func() {
		err := u.Tui.Run()
		if err != nil {
			log.Fatal(err)
		}
		u.Close()
	}()
}

func (u TUIUserInterface) Write(output []byte) (written int, err error) {
	u.scrollback.Append(tui.NewHBox(tui.NewLabel(string(output))))
	return len(output), nil
}

func (u TUIUserInterface) WriteString(text string) {
	if _, err := u.Write([]byte(text)); err != nil {
		log.Fatalf("Write Failed: %s", err)
	}
}

func NewTUIUserInterface() *TUIUserInterface {
	ui := &TUIUserInterface{}
	ui.output = make(chan string, 200)

	scrollback := tui.NewVBox()
	scrollarea := tui.NewScrollArea(scrollback)
	scrollarea.SetAutoscrollToBottom(true)
	scrollBackBox := tui.NewVBox(scrollarea)
	scrollBackBox.SetBorder(true)

	entry := tui.NewEntry()
	entry.SetSizePolicy(tui.Expanding, tui.Maximum)

	entryBox := tui.NewHBox(entry)
	entryBox.SetBorder(true)
	entryBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	window := tui.NewVBox(scrollBackBox, entryBox)
	window.SetBorder(false)

	screen := tui.NewHBox()
	screen.Append(window)

	entry.SetFocused(true)

	tuiui, err := tui.New(screen)
	if err != nil {
		log.Fatal(err)
	}

	ui.Tui = tuiui
	ui.Tui.SetKeybinding("Esc", func() { *ui.commands <- "quit" })
	ui.entry = entry
	ui.scrollback = scrollback

	entry.OnSubmit(func(e *tui.Entry) {
		if ui.commands != nil {
			*ui.commands <- e.Text()
		}
		e.SetText("")
	})

	return ui
}

func (u *TUIUserInterface) Error(text string) {
	u.WriteString("** ERR: " + text)
}
