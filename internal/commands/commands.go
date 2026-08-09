package commands

import (
	"fmt"
	"github.com/bitofaphilistine/blog-aggregator/internal/config"
)


type State struct {
	Config *Config
}

type Command struct {
	Name string
	Args []string
}

var Commands := struct{
	Commands map[string]func(state *State, c Command) error
}{
	Commands: map[string]func(state *State, c Command) error{
		"login": handlerLogin,
	},
}

func (c *Commands) Run(s *State, cmd Command) error {
	if handler, ok := c.Commands[cmd.Name]; ok {
		return handler(s, cmd)
	}
	return fmt.Errorf("unknown command: %s", cmd.Name)
}

func (c *Commands) register(name string, f func(*state, command) error) {
	c.Commands[name] = f
}


func handlerLogin(s *State, cmd Command) error {
	if cmd.Args == nil || len(cmd.Args) < 1 {
		return fmt.Errorf("login command requires a username argument")
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}
	fmt.Println("User set to:", cmd.Args[0])

	return nil
}