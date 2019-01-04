package main

import "fmt"
import "os"

func parse_arguments() {
	args := os.Args[1:]

	// check if a command is given
	if len(args) == 0 {
		fmt.Printf("No command: expected a command\n")
		os.Exit(1)
	}

	// HERE PUT CONDITIONS FOR COMMANDS

	// if nothing works
	fmt.Printf("Unrecognized command \"%s\"\n", args[0])
	os.Exit(1)
}

func main() {
	// get into the decks directory
	os.Chdir("/home/grastello/flashcards")

	parse_arguments()

	os.Exit(0)
}
