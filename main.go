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
	if os.Args == nil || len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	if os.Args[1] == "init" {
		if len(os.Args) < 2 {
			fmt.Println("init command requires connection url")
			os.Exit(1)
		}

		if err := config.Initialize(os.Args[2]); err != nil {
			fmt.Println("failed to initialize config: %w", err)
			os.Exit(1)
		}

		fmt.Println("Config initialized")
		os.Exit(0)
	} else if os.Args[1] == "help" {
		fmt.Println(`Commands:
help: print this message
init [connection_url]: initialize the config file with the given connection url
register [username]: add a user to the database
login [username]: login as the given user (saved)
users: list all users in the database
addfeed [name] [url]: add an rss feed at the given url, starts as followed by the logged in user
feeds: list all registered feeds
follow [url]: follow the given feed as the current user
unfollow [url]: unfollow the given feed as the current user
following: list all feeds followed by the current user
agg [interval]: downloads posts from each feed every interval
browse [amount(optional)]: lists the latest posts from feeds followed by the current user
reset [database(optional)]: reset the given database (users, feeds, posts, feed_follows, all), defaults to all, any databases dependant on the one reset will also be reset`)
		os.Exit(0)
	}

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

	cmd := commands.Command{
		Name: os.Args[1],
		Args: os.Args[2:],
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

	err = cmds.Run(&state, cmd)
	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}
}