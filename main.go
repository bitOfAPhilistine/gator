package main

import (
	"os"
	"fmt"
	"github.com/bitofaphilistine/internal/config"
)

user := os.Getenv("USER")

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

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

	fmt.Println(cfg)
}