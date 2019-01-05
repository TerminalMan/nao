package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Flashcard struct {
	front       string
	back        string
	efactor     float64
	duedate     int
	repetitions int
	interval    int
}

// get then n-th flashcard of the given deck
func get_card(deck string, n int) Flashcard {
	// set up the file for reading through a scanner
	deck_f, err := os.Open(deck)
	if err != nil {
		fmt.Printf("\033[1;31mError:\033[0m no \"%s\" deck found\n", deck)
		os.Exit(1)
	}
	defer deck_f.Close()
	deck_s := bufio.NewScanner(deck_f)

	// skips to the line that has the wanted flashcard's data
	for i := 0; i <= n; i++ {
		deck_s.Scan()
	}
	line := deck_s.Text()

	// extract front
	var i int
	var front string
	for i = 0; line[i] != ';'; i++ {
		front += string(line[i])
	}
	line = line[i+1:]

	// extract back
	var back string
	for i = 0; line[i] != ';'; i++ {
		back += string(line[i])
	}
	line = line[i+1:]

	// extract e-factor
	var efactor_s string
	for i = 0; line[i] != ';'; i++ {
		efactor_s += string(line[i])
	}
	line = line[i+1:]
	efactor, _ := strconv.ParseFloat(efactor_s, 64)

	// extract due date
	var duedate_s string
	for i = 0; line[i] != ';'; i++ {
		duedate_s += string(line[i])
	}
	line = line[i+1:]
	duedate, _ := strconv.Atoi(duedate_s)

	// extract repetitions
	var repetitions_s string
	for i = 0; line[i] != ';'; i++ {
		repetitions_s += string(line[i])
	}
	line = line[i+1:]
	repetitions, _ := strconv.Atoi(repetitions_s)

	// extract interval
	interval_s := line
	interval, _ := strconv.Atoi(interval_s)

	return Flashcard{front, back, efactor, duedate, repetitions, interval}
}

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
