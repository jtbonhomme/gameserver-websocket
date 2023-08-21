package skyjo

import (
	"math/rand"
	"time"
)

const (
	TotalCards         int = 150
	DefaultCardsNumber int = 10
	MinusTwoNumber     int = 5
	ZeroNumber         int = 15
	CardsPerPlayer     int = DeckWidth * DeckHeight
)

const (
	TakeFromDiscard int = iota
	TakeFromDraw
)

func generateCards() []Card {
	var index int
	cards := make([]Card, TotalCards)

	for i := 0; i < 5; i++ {
		cards[index].Value = MinusTwo
		index++
	}

	for i := 0; i < 15; i++ {
		cards[index].Value = 0
		index++
	}

	for j := -1; j < 13; j++ {
		if j == 0 {
			continue
		}

		for i := 0; i < 10; i++ {
			cards[index].Value = j
			index++
		}
	}

	return cards
}

func shuffle(slc []Card) {
	n := len(slc)

	for i := 0; i < n; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(n-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}

func GenerateCards() []Card {
	rand.Seed(time.Now().UTC().UnixNano())

	cards := generateCards()
	shuffle((cards))

	return cards
}
