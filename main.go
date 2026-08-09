package main

import (
	"fmt"
	"os"
	"github.com/bitofaphilistine/blog-aggregator/internal/config"
	"github.com/bitofaphilistine/blog-aggregator/internal/commands"
)

var user = "Philip"


func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	state := commands.State{
		Config: cfg,
	}

	cmds := commands.Commands{
		Registered: make(map[string]func(state *commands.State, c commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)

	if os.Args == nil || len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
	}

	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}