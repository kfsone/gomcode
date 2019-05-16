package main

func cmd_quit(ctx Context) error {
	ctx.User.Close()
	return nil
}

func cmd_help(ctx Context) error {
	ctx.User.WriteString("There's no help yet.")
	return nil
}
