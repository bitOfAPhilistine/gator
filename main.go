package main

import _ "github.com/lib/pq"
import (
	"os"
	"fmt"
	"context"
	"database/sql"
	"github.com/bitofaphilistine/gator/internal/config"
	"github.com/bitofaphilistine/gator/internal/commands"
	"github.com/bitofaphilistine/gator/internal/database"
)

var user = "Philip"


func loginCheck(handler func(*commands.State, commands.Command, database.User) error) func(*commands.State, commands.Command) error {
	return func(s *commands.State, cmd commands.Command) error {
		user, err := s.Queries.GetUserByName(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		return handler(s, cmd, user)
	}
}

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
	cmds.Register("agg", commands.HandlerAggregate)
	cmds.Register("addfeed", loginCheck(commands.HandlerAddFeed))
	cmds.Register("feeds", commands.HandlerListFeeds)
	cmds.Register("follow", loginCheck(commands.HandlerFollowFeed))
	cmds.Register("unfollow", loginCheck(commands.HandlerUnfollowFeed))
	cmds.Register("following", loginCheck(commands.HandlerListFollowedFeeds))
	cmds.Register("browse", loginCheck(commands.HandlerListPosts))

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