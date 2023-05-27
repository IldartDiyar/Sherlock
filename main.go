package main

import (
	"fmt"
	"os"

	"passive/internal/app"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Println("Welcome to passive v1.0.0")
		fmt.Println()
		fmt.Println("OPTIONS:")
		fmt.Println("-ip         Search with ip address")
		fmt.Println("-u          Search with username")
		os.Exit(0)
	}
	app.ParseArg()
}
