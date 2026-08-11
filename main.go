package main

import _ "github.com/lib/pq"
import (
	"fmt"
	"os"
	"database/sql"
	"github.com/bitofaphilistine/gator/internal/config"
	"github.com/bitofaphilistine/gator/internal/commands"
	"github.com/bitofaphilistine/gator/internal/database"
)

var user = "Philip"


func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	state := commands.State{
		Config: cfg,
		Database: db,
		Queries: database.New(db),
	}

	cmds := commands.Commands{
		Registered: make(map[string]func(state *commands.State, c commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)
	cmds.Register("register", commands.HandlerRegister)
	cmds.Register("reset", commands.HandlerResetDb)
	cmds.Register("users", commands.HandlerListUsers)

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