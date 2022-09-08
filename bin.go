package main

import (
	_ "embed"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

//go:embed index.html
var indexTemplate string

var deck = []string{
	"A♠", "2♠", "3♠", "4♠", "5♠", "6♠", "7♠", "8♠", "9♠", "10♠", "J♠", "Q♠", "K♠",
	"A♣", "2♣", "3♣", "4♣", "5♣", "6♣", "7♣", "8♣", "9♣", "10♣", "J♣", "Q♣", "K♣",
	"A♥", "2♥", "3♥", "4♥", "5♥", "6♥", "7♥", "8♥", "9♥", "10♥", "J♥", "Q♥", "K♥",
	"A♦", "2♦", "3♦", "4♦", "5♦", "6♦", "7♦", "8♦", "9♦", "10♦", "J♦", "Q♦", "K♦",
}

func generateSheet(seed int64) []string {
	// this makes a copy of a standard 52 card deck, then selects 24 unique random cards from it
	// it is important to note that if you provide the same seed twice it will generate the same deck

	// make a copy of the original deck
	thisDeck := []string{}
	thisDeck = append(thisDeck, deck...)

	// generate a source of randomness using the seed value, if you don't you will get the same deck every time
	sourceRandomness := rand.NewSource(seed)
	thisRand := rand.New(sourceRandomness)

	// select a random card 24 times, add it to the sheet
	bingoSheet := []string{}
	for i := 1; i <= 25; i++ {
		if i == 13 {
			bingoSheet = append(bingoSheet, "$")
			continue
		}
		// select a random number in the range of the deck
		randomNumber := thisRand.Intn(len(thisDeck))
		// select the card
		randomCard := thisDeck[randomNumber]
		// remove that card from the deck
		thisDeck = append(thisDeck[:randomNumber], thisDeck[randomNumber+1:]...)
		// add the card to the bingo sheet
		bingoSheet = append(bingoSheet, randomCard)
	}

	return bingoSheet
}

type PageData struct {
	Cells [5][5]string
}

func generatePageData(seed int64) PageData {
	// creates a 5x5 array of strings for the HTML template
	cells := [5][5]string{}
	for i, cell := range generateSheet(seed) {
		row := i / 5
		col := i % 5
		cells[row][col] = cell
	}
	return PageData{Cells: cells}
}

func getPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.New("index").Parse(indexTemplate)

	// we need a seed value for generating a random deck
	var seed int64
	// look for a seed value in the client's cookie
	cookieValue, err := r.Cookie("seed")
	// if you can't find the seed in the cookie, use the current timestamp
	// and send it back with the response
	if err != nil {
		seed = time.Now().UnixNano()
		// log.Println("request without a valid cookie, generating a new seed:", seed)

		cookie := &http.Cookie{
			Name:   "seed",
			Value:  strconv.FormatInt(seed, 10),
			MaxAge: 300,
		}
		http.SetCookie(w, cookie)
	} else {
		// otherwise, use existing seed value in the cookie
		seed, _ = strconv.ParseInt(cookieValue.Value, 10, 64)
		// log.Println("found seed:", seed)
	}

	// generate a bingo sheet using this seed
	data := generatePageData(seed)
	tmpl.Execute(w, data)
}

func main() {
	// server the getPage funtion, bind to port 8000
	http.HandleFunc("/", getPage)
	http.ListenAndServe(":8000", nil)
}
