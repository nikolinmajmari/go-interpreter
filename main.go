package main

import (
	"fmt"
	"interpreter/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %, this is Monkey Language Terminal\n", usr.Name)
	fmt.Printf("Feel free to type any command! \n")
	repl.Start(os.Stdin, os.Stdout)
}
