package commands

import (
	"fmt"
	"github.com/bitofaphilistine/blog-aggregator/internal/config"
)


type State struct {
	Config *config.Config
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

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Println("User set to:", cmd.Args[0])

	return nil
}