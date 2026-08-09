package main

import (
	"fmt"
	"github.com/bitofaphilistine/blog-aggregator/internal/config"
	"github.com/bitofaphilistine/blog-aggregator/internal/commands"
)

var user = "Philip"


func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println(*cfg)

	err = cfg.SetUser(user)
	if err != nil {
		fmt.Println("Error setting user:", err)
		return
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println(*cfg)
}