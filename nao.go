package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
	"github.com/mattn/go-runewidth"
)

// global variables
var INTERVAL0 int = 1
var INTERVAL1 int = 2
var MAXINTERVAL int = 0
var LINELENGTH int = 79
var DECKDIR string = ""

type Flashcard struct {
	front       string
	back        string
	eFactor     float64
	dueDate     int
	repetitions int
	interval    int
}

// get then n-th flashcard of the given deck
func getCard(deck string, n int) Flashcard {
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
	eFactor_s := line[:i]
	line = line[i+1:]
	eFactor, _ := strconv.ParseFloat(eFactor_s, 64)

	// extract due date
	for i = 0; line[i] != ';'; i++ {
	}
	dueDate_s := line[:i]
	line = line[i+1:]
	dueDate, _ := strconv.Atoi(dueDate_s)

	// extract repetitions
	for i = 0; line[i] != ';'; i++ {
	}
	repetitions_s := line[:i]
	line = line[i+1:]
	repetitions, _ := strconv.Atoi(repetitions_s)

	// extract interval
	interval_s := line
	interval, _ := strconv.Atoi(interval_s)

	return Flashcard{front, back, eFactor, dueDate, repetitions, interval}
}

// get the number of cards in the given deck
func getDeckn(deck string) int {
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

// read a single character from stin
func getKey(prompt string) byte {
	// prompt
	fmt.Printf("%s ", prompt)

	// read
	reader := bufio.NewReader(os.Stdin)
	char, err := reader.ReadByte()
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	return char
}

// print function designed for handling long flashcards in a nice way
func prettyPrint(s1, s2 string) int {
	lines := 0
	s1len := len(s1)
	s2 = strings.Replace(s2, "<br>", "\n", -1)

	fmt.Printf("\033[1m%s\033[0m", s1)

	col := len(s1)+1
	for _, rune := range s2 {
		fmt.Printf("%c", rune)
		col += runewidth.RuneWidth(rune)
		if rune == '\n' {
			fmt.Printf("%s",strings.Repeat(" ",s1len))
			col = s1len
			lines++
		}

		if col > LINELENGTH {
			fmt.Printf("\n")
			fmt.Printf("%s",strings.Repeat(" ",s1len))
			col = s1len
			lines++
		}
	}
	return lines
}

// get today's date in unix time
func getToday() int {
	t := int(time.Now().Unix())
	t -= t % 86400
	return t
}

// study card and return the updated (or not, if cramming) card
func studyCard(card Flashcard, cram bool) Flashcard {
	// show the card and gather answer quality
	prettyPrint("Front: ", card.front)
	getKey("")
	prettyPrint("Back:  ", card.back)
	fmt.Printf("\n")

	// if cramming wait for enter and return
	if cram {
		getKey("")
		return card
	}

	// if not cramming get the user to evaluate his answer and then
	// update the flashcard data
	quality := getKey("\033[1mEvaluate your answer [0-5]:\033[0m") - '0'
	for quality > 5 || quality < 0 {
		quality = getKey("\033[1mEvaluate your answer [0-5]:\033[0m") - '0'
	}

	// update e-factor
	card.eFactor += 0.1 - (5-float64(quality))*(0.08+(5-float64(quality))*0.02)
	if card.eFactor > 2.5 {
		card.eFactor = 2.5
	} else if card.eFactor < 1.3 {
		card.eFactor = 1.3
	}

	if quality >= 3 {
		if card.repetitions == 0 {
			card.interval = INTERVAL0
		} else if card.repetitions == 1 {
			card.interval = INTERVAL1
		} else {
			card.interval = int(math.Floor(float64(card.interval) * card.eFactor))
			if card.interval >= MAXINTERVAL && MAXINTERVAL != 0 {
				card.interval = MAXINTERVAL
			}
		}

		card.repetitions += 1
		card.dueDate = getToday() + card.interval*86400
	} else {
		card.repetitions = 0
		card.interval = 0
		card.dueDate = getToday()
	}

	return card
}

// write the card to the file in the right position
func writeCard(deck string, card Flashcard, n int) {
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
	defer os.Remove("tmpdeck")

	// copy lines to tmpdeck until the line of interest
	for i := 0; i < n; i++ {
		deck_s.Scan()
		fmt.Fprintf(tmpdeck_f, "%s\n", deck_s.Text())
	}

	// insert the updated line to the tmpdeck
	deck_s.Scan()
	fmt.Fprintf(tmpdeck_f, "%s;", card.front)
	fmt.Fprintf(tmpdeck_f, "%s;", card.back)
	fmt.Fprintf(tmpdeck_f, "%f;", card.eFactor)
	fmt.Fprintf(tmpdeck_f, "%d;", card.dueDate)
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
func studyDeck(deck string, cram bool) {
	// make the support array
	rand.Seed(time.Now().Unix())
	deckn := getDeckn(deck)
	decka := rand.Perm(deckn)

	today := getToday()

	// set failed variable
	fail := false

	for i := 0; i < deckn; i++ {
		card := getCard(deck, decka[i])

		if card.dueDate <= today || cram {
			// review/cram deck
			card = studyCard(card, cram)

			if !cram {
				writeCard(deck, card, decka[i])
			}
		}

		if card.dueDate <= today && !cram {
			fail = true
		}
	}

	if fail {
		studyDeck(deck, cram)
	} else if cram {
		fmt.Printf("You have finished cramming \033[1m%s\033[0m!\n", deck)
	} else {
		fmt.Printf("You have finished studying \033[1m%s\033[0m for today!\n", deck)
	}
}

func addCard(deck string, front string, back string) {
	// add the card to the deck, creating it if it does not exists yet
	deck_f, _ := os.OpenFile(deck, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer deck_f.Close()

	fmt.Fprintf(deck_f, "%s;%s;2.5;%d;0;0\n", front, back, getToday())
}

func infoDeck(deck string) {
	deckn := getDeckn(deck)
	dueToday := 0
	dueTomorrow := 0
	averageEfactor := 0.0

	today := getToday()

	// get data
	for i := 0; i < deckn; i++ {
		card := getCard(deck, i)

		averageEfactor += card.eFactor
		if card.dueDate <= today {
			dueToday += 1
		} else if card.dueDate == today+86400 {
			dueTomorrow += 1
		}
	}

	averageEfactor /= float64(deckn)
	overallEase := int(math.Floor(((averageEfactor-2.5)/1.3 + 1) * 100))

	// print data
	fmt.Printf("\033[1m%s\033[0m's infos\n", deck)
	prettyPrint("Card total:   ", strconv.Itoa(deckn))
	prettyPrint("Overall ease: ", strconv.Itoa(overallEase)+"/100")
	prettyPrint("Due today:    ", strconv.Itoa(dueToday))
	prettyPrint("Due tomorrow: ", strconv.Itoa(dueTomorrow))
}

func parseArguments() {
	args := os.Args[1:]

	// check if a command is given
	if len(args) == 0 {
		fmt.Printf("\033[1;31mError:\033[0m no command, expected a command\n")
		os.Exit(1)
	}

	switch args[0] {
	case "add", "a":
		if len(args) < 4 {
			fmt.Printf("\033[1;31mError:\033[0m the add command need 3 arguments\n")
			os.Exit(1)
		}

		if len(args) > 4 {
			fmt.Printf("\033[1;33mWarning:\033[0m extra arguments will be ignored\n")
		}

		addCard(args[1], args[2], args[3])
	case "review", "r":
		if len(args) < 2 {
			fmt.Printf("\033[1;31mError:\033[0m review need an argument\n")
			os.Exit(1)
		}

		for i := 1; i < len(args); i++ {
			studyDeck(args[i], false)
		}
	case "cram", "c":
		if len(args) < 2 {
			fmt.Printf("\033[1;31mError:\033[0m cram need an argument\n")
			os.Exit(1)
		}

		for i := 1; i < len(args); i++ {
			studyDeck(args[i], true)
		}
	case "info", "i":
		if len(args) < 2 {
			fmt.Printf("\033[1;31mError:\033[0m info need an argument\n")
			os.Exit(1)
		}

		for i := 1; i < len(args); i++ {
			infoDeck(args[i])
			if i != len(args)-1 {
				fmt.Printf("\n")
			}
		}
	default:
		fmt.Printf("\033[1;31mError:\033[0m unrecognized command \"%s\"\n", args[0])
		os.Exit(1)
	}

	os.Exit(0)
}

// read the config file and set up variables accordingly, returing error when
// unxpected things happen
func parseConfig(configfile_f *os.File) {
	configfile_s := bufio.NewScanner(configfile_f)
	i := 0

	for configfile_s.Scan() {
		i++
		words := strings.Fields(configfile_s.Text())

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "interval0":
			if len(words) == 1 {
				fmt.Printf("\033[1;31mError:\033[0m no argument provided on line %d of naorc\n", i)
				os.Exit(1)
			}

			INTERVAL0, _ = strconv.Atoi(words[1])
		case "interval1":
			if len(words) == 1 {
				fmt.Printf("\033[1;31mError:\033[0m no argument provided on line %d of naorc\n", i)
				os.Exit(1)
			}

			INTERVAL1, _ = strconv.Atoi(words[1])
		case "maxinterval":
			if len(words) == 1 {
				fmt.Printf("\033[1;31mError:\033[0m no argument provided on line %d of naorc\n", i)
				os.Exit(1)
			}

			MAXINTERVAL, _ = strconv.Atoi(words[1])
		case "linelength":
			if len(words) == 1 {
				fmt.Printf("\033[1;31mError:\033[0m no argument provided on line %d of naorc\n", i)
				os.Exit(1)
			}

			// set LINELENGTH, then make it odd if even since odd linelenghts
			// don't break wide unicode characters
			LINELENGTH, _ = strconv.Atoi(words[1])
			if LINELENGTH%2 == 0 {
				LINELENGTH += 1
			}
		case "deckdir":
			if len(words) == 1 {
				fmt.Printf("\033[1;31mError:\033[0m no argument provided on line %d of naorc\n", i)
				os.Exit(1)
			}

			// get everything after the deckdir option and set that as the DECKDIR
			j := 0
			for j = 0; j < len(configfile_s.Text()); j++ {
				if configfile_s.Text()[j] == ' ' {
					break
				}
			}

			DECKDIR = configfile_s.Text()[j+1:]
		default:
			fmt.Printf("\033[1;31mError:\033[0m unrecognized option \"%s\" on line %d of naorc\n", words[0], i)
			os.Exit(1)
		}
	}
}

func init() {
	// set default DECKDIR
	user, _ := user.Current()
	DECKDIR = user.HomeDir + "/flashcards"

	// open config file, if any the parse it and set the variables accordingly
	configfile_f, err := os.Open(user.HomeDir + "/.config/nao/naorc")
	if err == nil {
		defer configfile_f.Close()
		parseConfig(configfile_f)
	}

	// expand '~', check if the path is absolute then check if DECKDIR is
	// a directory and if everything is ok change current directory
	// to DECKDIR
	if DECKDIR[0] == '~' {
		DECKDIR = user.HomeDir + DECKDIR[1:]
	} else if DECKDIR[0] != '/' {
		fmt.Printf("\033[1;31mError:\033[0m \"%s\" is not an absolute path; check your naorc\n", DECKDIR)
		os.Exit(1)
	}

	deckdirStats, err := os.Stat(DECKDIR)
	if err != nil {
		fmt.Printf("\033[1;31mError:\033[0m \"%s\" no such directory\n", DECKDIR)
		os.Exit(1)
	} else if deckdirStats.IsDir() == false {
		fmt.Printf("\033[1;31mError:\033[0m \"%s\" is not a directory\n", DECKDIR)
		os.Exit(1)
	}

	os.Chdir(DECKDIR)
}

func main() {
	parseArguments()

	os.Exit(0)
}
