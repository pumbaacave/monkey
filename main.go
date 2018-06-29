package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hi %s, MONKEY Interpreter starts\n", user.Username)
	fmt.Printf("Pls type some commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
