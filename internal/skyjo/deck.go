package skyjo

import (
	"errors"
	"fmt"
	"math"
)

const (
	DeckWidth  int = 4
	DeckHeight int = 3
)

type Deck struct {
	cards []Card
	width int
}

func NewDeck() *Deck {
	d := &Deck{}
	d.Reset()
	return d
}

func (d *Deck) Reset() {
	d.cards = []Card{}
	d.width = DeckWidth
}

func (d *Deck) cardAtIndex(i int) (Card, error) {
	if i >= d.width*DeckWidth || i < 0 {
		return Card{}, errors.New("index out of range")
	}

	return d.cards[i], nil
}

func (d *Deck) IsVisible(i int) (bool, error) {
	if i >= d.width*DeckHeight || i < 0 {
		return false, errors.New("index out of range")
	}

	return d.cards[i].Visible, nil
}

func (d *Deck) RevealCard(i int) error {
	if i >= d.width*DeckHeight || i < 0 {
		return errors.New("index out of range")
	}

	c, err := d.cardAtIndex(i)
	if err != nil {
		return fmt.Errorf("error getting card at index %d: %s", i, err.Error())
	}

	if c.Visible {
		return fmt.Errorf("card at index %d is already visible", i)
	}

	d.cards[i].Visible = true

	return nil
}

func (d *Deck) Set(i int, card Card) error {
	if i >= d.width*DeckHeight || i < 0 {
		return errors.New("index out of range")
	}

	d.cards[i] = card

	return nil
}

func (d *Deck) Hidden() []int {
	hidden := []int{}

	for i, c := range d.cards {
		if !c.Visible {
			hidden = append(hidden, i)
		}
	}

	return hidden
}

func (d *Deck) AddCard(card Card) error {
	if len(d.cards) >= d.width*DeckHeight {
		return errors.New("index out of range")
	}

	d.cards = append(d.cards, card)

	return nil
}

func (d *Deck) posToIndex(c, r int) int {
	return r*d.width + c
}

func (d *Deck) indexToPos(i int) (int, int) {
	return i % d.width, i / d.width
}

func (d *Deck) Display() {
	for r := 0; r < DeckHeight; r++ {
		for c := 0; c < d.width; c++ {
			i := d.posToIndex(c, r)
			card := d.cards[i]
			if card.Visible {
				fmt.Printf(" % 2d", card.Value)
			} else {
				fmt.Printf("  .")
			}
		}
		fmt.Println("")
	}
}

func (d *Deck) PublicDeck() []Card {
	publicDeck := make([]Card, d.width*DeckHeight)

	for r := 0; r < DeckHeight; r++ {
		for c := 0; c < d.width; c++ {
			i := d.posToIndex(c, r)

			card := d.cards[i]
			if card.Visible {
				publicDeck[i] = card
			}
		}
	}

	return publicDeck
}

func (d *Deck) GetCard(i int) (Card, error) {
	if i >= len(d.cards) || i < 0 {
		return Card{}, fmt.Errorf("invalid index: %d", i)
	}

	card := d.cards[i]
	return card, nil
}

func (d *Deck) Len() int {
	return len(d.cards)
}

func (d *Deck) Width() int {
	return d.width
}

func (d *Deck) RevealAll() {
	for r := 0; r < DeckHeight; r++ {
		for c := 0; c < d.width; c++ {
			i := d.posToIndex(c, r)
			d.cards[i].Visible = true
		}
	}
}

func (d *Deck) Value() int {
	var value int
	for r := 0; r < DeckHeight; r++ {
		for c := 0; c < d.width; c++ {
			i := d.posToIndex(c, r)
			card := d.cards[i]
			if card.Visible {
				value += card.Value
			}
		}
	}
	return value
}

func (d *Deck) CheckDeck() error {
	for col := 0; col < d.width; col++ {
		toRemove, err := d.checkColumn(col)
		if err != nil {
			return fmt.Errorf("error checking column %d: %w", col, err)
		}
		if toRemove {
			err := d.removeColumn(col)
			if err != nil {
				return fmt.Errorf("error removing column %d: %w", col, err)
			}
		}
	}
	return nil
}

func (d *Deck) checkColumn(col int) (bool, error) {
	if col >= d.width || col < 0 {
		return false, fmt.Errorf("index %d larger than deck width %d or negative", col, d.width)
	}

	cards, err := d.GetColumn(col)
	if err != nil {
		return false, fmt.Errorf("error getting column %d: %w", col, err)
	}

	lastValue := math.MaxInt
	for _, card := range cards {
		if lastValue == math.MaxInt {
			lastValue = card.Value
		}
		if card.Value != lastValue || card.Visible == false {
			return false, nil
		}
	}
	return true, nil
}

func (d *Deck) GetColumn(col int) ([]Card, error) {
	if col >= d.width || col < 0 {
		return []Card{}, fmt.Errorf("index %d larger than deck width %d or negative", col, d.width)
	}
	cards := []Card{}
	for r := 0; r < DeckHeight; r++ {
		i := r*d.width + col
		cards = append(cards, d.cards[i])
	}

	return cards, nil
}

func (d *Deck) removeColumn(col int) error {
	if col >= d.width || col < 0 {
		return fmt.Errorf("index %d larger than deck width %d or negative", col, d.width)
	}
	if d.width == 0 || len(d.cards) == 0 {
		return fmt.Errorf("can not remove column, deck is empty")
	}

	for r := DeckHeight - 1; r >= 0; r-- {
		i := r*d.width + col
		d.removeItemFromDeck(i)
	}

	d.width--

	return nil
}

func (d *Deck) removeItemFromDeck(i int) {
	// Remove the element at index i from a.
	copy(d.cards[i:], d.cards[i+1:])   // Shift a[i+1:] left one index.
	d.cards[len(d.cards)-1] = Card{}   // Erase last element (write zero value).
	d.cards = d.cards[:len(d.cards)-1] // Truncate slice.
}
