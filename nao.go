package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func add_card(deck string) {
	// read front and back from the user
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("\033[1mFront:\033[0m ")
	front, _ := reader.ReadString('\n')
	front = front[:len(front)-1]

	fmt.Printf("\033[1mBack:\033[0m  ")
	back, _ := reader.ReadString('\n')
	back = back[:len(back)-1]

	// get toadys date in unix time
	today := time.Now().Unix()
	today -= today % 86400

	// add the card to the deck, creating it if it does not exists yet
	deck_f, _ := os.OpenFile(deck, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer deck_f.Close()

	fmt.Fprintf(deck_f, "%s;%s;2.5;%d;0;0\n", front, back, today)
}

func parse_arguments() {
	args := os.Args[1:]

	// check if a command is given
	if len(args) == 0 {
		fmt.Printf("No command: expected a command\n")
		os.Exit(1)
	}

	// add command
	if args[0] == "add" {
		if len(args) < 2 {
			fmt.Printf("No argument: add need a deck as an argument\n")
			os.Exit(1)
		}

		if len(args) > 2 {
			fmt.Printf("\033[1;33mWarning:\033[0m extra arguments will be ignored\n")
		}

		add_card(args[1])
		os.Exit(0)
	}

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
