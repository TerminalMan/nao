package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var INTERVAL_0 int = 1
var INTERVAL_1 int = 2
var LINELENGTH int = 79

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

	var i int

	// extract front
	for i = 0; line[i] != ';'; i++ {
	}
	front := line[:i]
	line = line[i+1:]

	// extract back
	for i = 0; line[i] != ';'; i++ {
	}
	back := line[:i]
	line = line[i+1:]

	// extract e-factor
	for i = 0; line[i] != ';'; i++ {
	}
	efactor_s := line[:i]
	line = line[i+1:]
	efactor, _ := strconv.ParseFloat(efactor_s, 64)

	// extract due date
	for i = 0; line[i] != ';'; i++ {
	}
	duedate_s := line[:i]
	line = line[i+1:]
	duedate, _ := strconv.Atoi(duedate_s)

	// extract repetitions
	for i = 0; line[i] != ';'; i++ {
	}
	repetitions_s := line[:i]
	line = line[i+1:]
	repetitions, _ := strconv.Atoi(repetitions_s)

	// extract interval
	interval_s := line
	interval, _ := strconv.Atoi(interval_s)

	return Flashcard{front, back, efactor, duedate, repetitions, interval}
}

// get the number of cards in the given deck
func get_deckn(deck string) int {
	// set up the file for reading through a scanner
	deck_f, err := os.Open(deck)
	if err != nil {
		fmt.Printf("\033[1;31mError:\033[0m no \"%s\" deck found\n", deck)
		os.Exit(1)
	}
	defer deck_f.Close()
	deck_s := bufio.NewScanner(deck_f)

	i := 0
	for deck_s.Scan() {
		i++
	}

	return i
}

// read a single character from stin without a need for the enter key
func getkey() byte {
	// block terminal buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	// get the character
	var c []byte = make([]byte, 1)
	os.Stdin.Read(c)

	// reset terminal properties and return
	exec.Command("stty", "-F", "/dev/tty", "sane").Run()
	return c[0]
}

// print function designed for handling long flashcards in a nice way
// returns the number of lines printed
func pretty_print(s1, s2 string) int {
	lines := 0
	s1n := len(s1)

	fmt.Printf("\033[1m%s\033[0m", s1)

	if len(s2) <= LINELENGTH-s1n {
		fmt.Printf("%s\n", s2)
		s2 = ""
	} else {
		ss := s2[:LINELENGTH-s1n]
		s2 = s2[LINELENGTH-s1n:]
		fmt.Printf("%s\n", ss)
	}
	lines += 1

	for len(s2) > 0 {
		for i := 0; i < s1n; i++ {
			fmt.Printf(" ")
		}

		if len(s2) <= LINELENGTH-s1n {
			fmt.Printf("%s\n", s2)
			s2 = ""
		} else {
			ss := s2[:LINELENGTH-s1n]
			s2 = s2[LINELENGTH-s1n:]
			fmt.Printf("%s\n", ss)
		}

		lines += 1
	}

	return lines
}

// clear n lines of output
func clear_lines(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("\033[1A\r")

		for j := 0; j < LINELENGTH; j++ {
			fmt.Printf(" ")
		}
	}

	fmt.Printf("\r")
}

// study card and return the updated (or not, if cramming) card
func study_card(card Flashcard, cram bool) Flashcard {
	// show the card and gather answer quality
	lines := pretty_print("Front: ", card.front)
	getkey()
	lines += pretty_print("Back:  ", card.back)
	fmt.Printf("\033[1mEvaluate your answer:\033[0m \033[0;31m0 1 \033[0;33m2 3 \033[0;32m4 5\033[0m\n")
	lines += 1

	quality := getkey() - '0'
	for quality > 5 {
		quality = getkey() - '0'
	}

	clear_lines(lines)

	// update e-factor
	card.efactor += 0.1 - (5-float64(quality))*(0.08+(5-float64(quality))*0.02)
	if card.efactor > 2.5 {
		card.efactor = 2.5
	} else if card.efactor < 1.3 {
		card.efactor = 1.3
	}

	// get today's date. Update due date, interval and repetition number
	// according to the quality obtained
	today := int(time.Now().Unix())
	today -= today % 86400

	if quality >= 3 {
		if card.repetitions == 0 {
			card.interval = INTERVAL_0
		} else if card.repetitions == 1 {
			card.interval = INTERVAL_1
		} else {
			card.interval = int(math.Floor(float64(card.interval) * card.efactor))
		}

		card.repetitions += 1
		card.duedate = today + card.interval*86400
	} else {
		card.repetitions = 0
		card.interval = 0
		card.duedate = today
	}

	return card
}

// write the card to the file in the right position
func write_card(deck string, card Flashcard, n int) {
	// set up the file for reading through a scanner
	deck_f, err := os.Open(deck)
	if err != nil {
		fmt.Printf("\033[1;31mError:\033[0m no \"%s\" deck found\n", deck)
		os.Exit(1)
	}
	defer deck_f.Close()
	deck_s := bufio.NewScanner(deck_f)

	// set up a temporary file to write to
	tmpdeck_f, _ := os.Create("tmpdeck")
	defer tmpdeck_f.Close()

	// copy lines to tmpdeck until the line of interest
	for i := 0; i < n; i++ {
		deck_s.Scan()
		fmt.Fprintf(tmpdeck_f, "%s\n", deck_s.Text())
	}

	// insert the updated line to the tmpdeck
	deck_s.Scan()
	fmt.Fprintf(tmpdeck_f, "%s;", card.front)
	fmt.Fprintf(tmpdeck_f, "%s;", card.back)
	fmt.Fprintf(tmpdeck_f, "%f;", card.efactor)
	fmt.Fprintf(tmpdeck_f, "%d;", card.duedate)
	fmt.Fprintf(tmpdeck_f, "%d;", card.repetitions)
	fmt.Fprintf(tmpdeck_f, "%d\n", card.interval)

	// copy the remaining lines to the tmpdeck
	for deck_s.Scan() {
		fmt.Fprintf(tmpdeck_f, "%s\n", deck_s.Text())
	}

	// copy the tmpdeck to the actual deck
	tmpdeck_f, _ = os.Open("tmpdeck")
	defer tmpdeck_f.Close()
	tmpdeck_s := bufio.NewScanner(tmpdeck_f)

	deck_f, _ = os.Create(deck)
	defer deck_f.Close()

	for tmpdeck_s.Scan() {
		fmt.Fprintf(deck_f, "%s\n", tmpdeck_s.Text())
	}
}

// cram true: cram the given deck (just study every card)
// cram false: review the given deck (study due cards updating their
// local data, also repeat until every card has received a passing score)
func study_deck(deck string, cram bool) {
	// make the support array
	rand.Seed(time.Now().Unix())
	deckn := get_deckn(deck)
	decka := rand.Perm(deckn)

	// get today's date
	today := int(time.Now().Unix())
	today -= today % 86400

	// set failed variable
	fail := false

	for i := 0; i < deckn; i++ {
		card := get_card(deck, decka[i])

		if card.duedate <= today || cram {
			// review/cram deck
			card = study_card(card, cram)
			write_card(deck, card, decka[i])
		}

		if card.duedate <= today {
			fail = true
		}
	}

	if fail {
		study_deck(deck, cram)
	} else {
		fmt.Printf("You have finished studying \033[1m%s\033[0m for today!\n", deck)
	}
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
			fmt.Printf("\033[1;31mError:\033[0m the add command need an argument\n")
			os.Exit(1)
		}

		if len(args) > 2 {
			fmt.Printf("\033[1;33mWarning:\033[0m extra arguments will be ignored\n")
		}

		add_card(args[1])
		os.Exit(0)
	}

	// review commands
	if args[0] == "review" {
		if len(args) < 2 {
			fmt.Printf("\033[1;31mError:\033[0m review need an argument\n")
			os.Exit(1)
		}

		for i := 1; i < len(args); i++ {
			study_deck(args[i], false)
		}

		os.Exit(0)
	}

	// if nothing works
	fmt.Printf("\033[1;31mError:\033[0m unrecognized command \"%s\"\n", args[0])
	os.Exit(1)
}

func main() {
	// get into the decks directory
	os.Chdir("/home/grastello/flashcards")

	parse_arguments()

	os.Exit(0)
}
