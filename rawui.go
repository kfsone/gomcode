package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type RawUserInterface struct {
	UserInterface
}

func (u *RawUserInterface) Start() {
	u.UserInterface.Start()

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			if u.commands != nil {
				*u.commands <- strings.TrimRight(text, "\r\n 	")
			}
		}
	}()
	go func() {
		for text := range u.output {
			fmt.Println(text)
		}
	}()
}

func (u RawUserInterface) Write(text []byte) (count int, err error) {
	u.output <- string(text)
	return len(text), nil
}

func (u RawUserInterface) WriteString(text string) {
	if _, err := u.Write([]byte(text)); err != nil {
		log.Fatalf("Write Failed: %s", err)
	}
}

func NewRawUserInterface() *RawUserInterface {
	// read stdin until we need to close
	ui := &RawUserInterface{}
	ui.output = make(chan string, 200)
	return ui
}

func (u RawUserInterface) Error(text string) {
	u.WriteString("** Error: " + text)
}
