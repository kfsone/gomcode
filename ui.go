package main

type UserInterfacer interface {
	Close()
	Commands() chan string
	Start()
	Write(text []byte) (count int, err error)
	WriteString(text string)
	Error(text string)
}

type UserInterface struct {
	output   chan string  // application -> user
	commands *chan string // application <- user
}

func (u UserInterface) Close() {
	close(*u.commands)
}

func (u UserInterface) Commands() chan string {
	return *u.commands
}

func (u *UserInterface) Start() {
	if u.commands != nil {
		panic("Interface command channel already created")
	}
	commands := make(chan string, 200)
	u.commands = &commands
}
