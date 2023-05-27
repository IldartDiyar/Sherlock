package app

import (
	"fmt"
	"log"
	"os"

	"passive/internal/ip"
	"passive/internal/username"
)

var usage = "USAGE: ./main -[OPTION] 192.168.1.1"

func ParseArg() {
	args := os.Args[1:]
	if help, _ := hasElement(args, "-h", "--help"); help {
		fmt.Print(
			`passive --help

			Welcome to passive v1.0.0
			
OPTIONS:
		-ip         Search with ip address
		-u          Search with username
			

flags:
        -h, --help            show this help message and exit
		-ip         		  Search with ip address
		-u          		  Search with username
`,
		)
		os.Exit(0)
	}
	if len(args) != 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	if flag, argIndex := hasElement(args, "-ip"); flag {
		// do with ip
		log.Println(args[argIndex+1])
		ip.SearchByIP(args[argIndex+1])
	}

	if flag, argIndex := hasElement(args, "-u"); flag {
		// do with username
		log.Println(args[argIndex+1])
		username.SearchByUsername(args[argIndex+1])
	}
}

func hasElement(array []string, targets ...string) (bool, int) {
	for index, item := range array {
		for _, target := range targets {
			if item == target {
				return true, index
			}
		}
	}
	return false, -1
}
