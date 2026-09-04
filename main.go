package main

import (
	"fmt"
	"os"

	"github.com/JiveshL-KDK/go-gator/internal/commands"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		fmt.Println("please input a command")
		return
	}

	command := args[1]

	extraArgs := make([]string, 0)

	if len(args) >= 3 {
		extraArgs = append(extraArgs, args[2:]...)
	}

	cm := commands.NewCommand(command, extraArgs)

	if err := cm.Run(); err != nil {
		fmt.Println(err)
		return
	}

}
