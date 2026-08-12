package commands

import (
	"fmt"
	"time"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/bitofaphilistine/gator/internal/rss"
	"github.com/bitofaphilistine/gator/internal/config"
	"github.com/bitofaphilistine/gator/internal/database"
)


type State struct {
	Config *config.Config
	Database *sql.DB
	Queries *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct{
	Registered map[string]func(state *State, c Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.Registered[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Registered[name] = f
}


func HandlerLogin(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}

	user, err := s.Queries.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.Config.SetUser(user.Username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Println("User set to:", user.Username)

	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("register command requires a username argument")
	}

	userParams := database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username: cmd.Args[0],
	}
	user, err := s.Queries.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	fmt.Println("User created:", user.Username)
	fmt.Println(user)

	err = s.Config.SetUser(user.Username)
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	return nil
}

func HandlerResetDb(s *State, cmd Command) error {
	err := s.Queries.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users: %w", err)
	}
	fmt.Println("Users table reset successfully.")
	return nil
}

func HandlerListUsers(s *State, cmd Command) error {
	users, err := s.Queries.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to list users: %w", err)
	}

	fmt.Println("Registered Users:")
	for _, user := range users {
		if user.Username == s.Config.CurrentUserName {
			fmt.Println(" *", user.Username, "(current)")
		} else {
			fmt.Println(" *", user.Username)
		}
	}

	return nil
}

func HandlerAggregate(s *State, cmd Command) error {
	testUrl := "https://www.wagslane.dev/index.xml"
	res, err := rss.FetchFeed(context.Background(), testUrl)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Println(res)
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 2 {
		return fmt.Errorf("addfeed command requires a name and url argument")
	}

	user, err := s.Queries.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	feed, err := s.Queries.CreateFeed(context.Background(), database.CreateFeedParams{
		Name: cmd.Args[0],
		Url: cmd.Args[1],
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	
	fmt.Println(feed)
	return nil
}