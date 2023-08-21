package skyjo

import (
	"testing"
)

func TestCardAtIndex(t *testing.T) {
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i})
	}

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		card, err := d.cardAtIndex(i)
		if err != nil {
			t.Errorf("expected nil error at index %d, got %s", i, err.Error())
		}

		if card.Value != i {
			t.Errorf("expected card value %d at index %d, got %d", i, i, card.Value)
		}
	}
}

type PosToIndexTest struct {
	c, r int
	i    int
}

/*
   | 0  1  2  3
---+------------
 0 | 0  1  2  3
 1 | 4  5  6  7
 2 | 8  9 10 11

*/

var posToIndexTests = []PosToIndexTest{
	{
		c: 0,
		r: 0,
		i: 0,
	},
	{
		c: 1,
		r: 0,
		i: 1,
	},
	{
		c: 3,
		r: 0,
		i: 3,
	},
	{
		c: 0,
		r: 1,
		i: 4,
	},
	{
		c: 1,
		r: 1,
		i: 5,
	},
	{
		c: 0,
		r: 2,
		i: 8,
	},
	{
		c: 2,
		r: 1,
		i: 6,
	},
	{
		c: 3,
		r: 2,
		i: 11,
	},
	{
		c: 3,
		r: 1,
		i: 7,
	},
}

func TestPosToIndex(t *testing.T) {
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i})
	}

	for _, test := range posToIndexTests {
		res := d.posToIndex(test.c, test.r)
		if res != test.i {
			t.Errorf("posToIndex expected to return %d for c = %d and r = %d, but got %d", test.i, test.c, test.r, res)
		}
	}
}

func TestIndexToPos(t *testing.T) {
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i})
	}

	for _, test := range posToIndexTests {
		c, r := d.indexToPos(test.i)
		if c != test.c || r != test.r {
			t.Errorf("indexToPos expected to return %d for c and %d for r at index %d, but got %d for c and %d for r", test.c, test.r, test.i, c, r)
		}
	}
}

func TestCheckDeck(t *testing.T) {
	var err error
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i, Visible: true})
	}

	d.Set(8, Card{Value: 8, Visible: false})

	d.Set(5, Card{Value: 1, Visible: true})
	d.Set(9, Card{Value: 1, Visible: false})

	d.Set(7, Card{Value: 3, Visible: true})
	d.Set(11, Card{Value: 3, Visible: true})

	if d.width != DeckWidth {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth, d.width)
	}
	if len(d.cards) != 12 {
		t.Errorf("deck expected to have %d cards, but got %d", 12, len(d.cards))
	}

	d.Display()

	// row 0
	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 1, 1, t)
	checkDeckValue(d, 2, 2, t)
	checkDeckValue(d, 3, 3, t)

	// row 1
	checkDeckValue(d, 4, 4, t)
	checkDeckValue(d, 5, 1, t)
	checkDeckValue(d, 6, 6, t)
	checkDeckValue(d, 7, 3, t)

	// row 2
	checkDeckValue(d, 8, 8, t)
	checkDeckValue(d, 9, 1, t)
	checkDeckValue(d, 10, 10, t)
	checkDeckValue(d, 11, 3, t)

	err = d.CheckDeck()
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if d.width != DeckWidth-1 {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth-1, d.width)
	}
	if len(d.cards) != 9 {
		t.Errorf("deck expected to have %d cards, but got %d", 9, len(d.cards))
	}

	// row 0
	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 1, 1, t)
	checkDeckValue(d, 2, 2, t)

	// row 1
	checkDeckValue(d, 3, 4, t)
	checkDeckValue(d, 4, 1, t)
	checkDeckValue(d, 5, 6, t)

	// row 2
	checkDeckValue(d, 6, 8, t)
	checkDeckValue(d, 7, 1, t)
	checkDeckValue(d, 8, 10, t)

	d.Display()
}

func TestCheckColumn(t *testing.T) {
	var err error
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i, Visible: true})
	}

	d.Set(8, Card{Visible: false})

	d.Set(5, Card{Value: 1, Visible: true})
	d.Set(9, Card{Value: 1, Visible: false})

	d.Set(7, Card{Value: 3, Visible: true})
	d.Set(11, Card{Value: 3, Visible: true})

	if d.width != DeckWidth {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth, d.width)
	}
	if len(d.cards) != 12 {
		t.Errorf("deck expected to have %d cards, but got %d", 12, len(d.cards))
	}

	check0, err := d.checkColumn(0)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if check0 != false {
		t.Errorf("expected check to return %T, but got %T", false, check0)
	}

	check1, err := d.checkColumn(1)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if check1 != false {
		t.Errorf("expected check to return %T, but got %T", false, check1)
	}

	check2, err := d.checkColumn(2)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if check2 != false {
		t.Errorf("expected check to return %T, but got %T", false, check2)
	}

	check3, err := d.checkColumn(3)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if check3 != true {
		t.Errorf("expected check to return %T, but got %T", true, check3)
	}

	check4, err := d.checkColumn(4)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if check4 != false {
		t.Errorf("expected check to return %T, but got %T", false, check4)
	}

	checkNegative, err := d.checkColumn(-1)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if checkNegative != false {
		t.Errorf("expected check to return %T, but got %T", false, checkNegative)
	}
}

func TestGetColumn(t *testing.T) {
	var err error
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i, Visible: true})
	}

	if d.width != DeckWidth {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth, d.width)
	}
	if len(d.cards) != 12 {
		t.Errorf("deck expected to have %d cards, but got %d", 12, len(d.cards))
	}

	cardsCol1, err := d.GetColumn(1)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if len(cardsCol1) != 3 {
		t.Errorf("expected to retrieve %d cards, but got %d", 3, len(cardsCol1))
	}
	checkCardValue(cardsCol1[0], 1, true, t)
	checkCardValue(cardsCol1[1], 5, true, t)
	checkCardValue(cardsCol1[2], 9, true, t)

	cardsCol4, err := d.GetColumn(4)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if len(cardsCol4) != 0 {
		t.Errorf("expected to retrieve %d cards, but got %d", 0, len(cardsCol4))
	}

	cardsColNegative, err := d.GetColumn(-1)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if len(cardsColNegative) != 0 {
		t.Errorf("expected to retrieve %d cards, but got %d", 0, len(cardsColNegative))
	}
}

func TestRemoveColumn(t *testing.T) {
	var err error
	d := NewDeck()

	for i := 0; i < DeckWidth*DeckHeight; i++ {
		_ = d.AddCard(Card{Value: i, Visible: true})
	}

	if d.width != DeckWidth {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth, d.width)
	}
	if len(d.cards) != 12 {
		t.Errorf("deck expected to have %d cards, but got %d", 12, len(d.cards))
	}

	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 3, 3, t)
	checkDeckValue(d, 4, 4, t)
	checkDeckValue(d, 7, 7, t)
	checkDeckValue(d, 8, 8, t)
	checkDeckValue(d, 11, 11, t)

	err = d.removeColumn(4)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if d.width != DeckWidth {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth, d.width)
	}
	if len(d.cards) != 12 {
		t.Errorf("deck expected to have %d cards, but got %d", 12, len(d.cards))
	}

	// row 0
	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 1, 1, t)
	checkDeckValue(d, 2, 2, t)
	checkDeckValue(d, 3, 3, t)

	// row 1
	checkDeckValue(d, 4, 4, t)
	checkDeckValue(d, 5, 5, t)
	checkDeckValue(d, 6, 6, t)
	checkDeckValue(d, 7, 7, t)

	// row 2
	checkDeckValue(d, 8, 8, t)
	checkDeckValue(d, 9, 9, t)
	checkDeckValue(d, 10, 10, t)
	checkDeckValue(d, 11, 11, t)

	err = d.removeColumn(-1)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if d.width != DeckWidth {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth, d.width)
	}
	if len(d.cards) != 12 {
		t.Errorf("deck expected to have %d cards, but got %d", 12, len(d.cards))
	}

	// row 0
	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 1, 1, t)
	checkDeckValue(d, 2, 2, t)
	checkDeckValue(d, 3, 3, t)

	// row 1
	checkDeckValue(d, 4, 4, t)
	checkDeckValue(d, 5, 5, t)
	checkDeckValue(d, 6, 6, t)
	checkDeckValue(d, 7, 7, t)

	// row 2
	checkDeckValue(d, 8, 8, t)
	checkDeckValue(d, 9, 9, t)
	checkDeckValue(d, 10, 10, t)
	checkDeckValue(d, 11, 11, t)

	err = d.removeColumn(1)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if d.width != DeckWidth-1 {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth-1, d.width)
	}
	if len(d.cards) != 9 {
		t.Errorf("deck expected to have %d cards, but got %d", 9, len(d.cards))
	}

	// row 0
	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 1, 2, t)
	checkDeckValue(d, 2, 3, t)

	// row 1
	checkDeckValue(d, 3, 4, t)
	checkDeckValue(d, 4, 6, t)
	checkDeckValue(d, 5, 7, t)

	// row 2
	checkDeckValue(d, 6, 8, t)
	checkDeckValue(d, 7, 10, t)
	checkDeckValue(d, 8, 11, t)

	err = d.removeColumn(2)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if d.width != DeckWidth-2 {
		t.Errorf("deck width expected to be %d, but got %d", DeckWidth-2, d.width)
	}
	if len(d.cards) != 6 {
		t.Errorf("deck expected to have %d cards, but got %d", 6, len(d.cards))
	}

	// row 0
	checkDeckValue(d, 0, 0, t)
	checkDeckValue(d, 1, 2, t)

	// row 1
	checkDeckValue(d, 2, 4, t)
	checkDeckValue(d, 3, 6, t)

	// row 2
	checkDeckValue(d, 4, 8, t)
	checkDeckValue(d, 5, 10, t)

	err = d.removeColumn(0)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if d.width != 1 {
		t.Errorf("deck width expected to be %d, but got %d", 1, d.width)
	}
	if len(d.cards) != 3 {
		t.Errorf("deck expected to have %d cards, but got %d", 3, len(d.cards))
	}

	// row 0
	checkDeckValue(d, 0, 2, t)

	// row 1
	checkDeckValue(d, 1, 6, t)

	// row 2
	checkDeckValue(d, 2, 10, t)

	err = d.removeColumn(0)
	if err != nil {
		t.Errorf("expected nil error, got %s", err.Error())
	}
	if d.width != 0 {
		t.Errorf("deck width expected to be %d, but got %d", 0, d.width)
	}
	if len(d.cards) != 0 {
		t.Errorf("deck expected to have no cards, but got %d", len(d.cards))
	}

	err = d.removeColumn(0)
	if err == nil {
		t.Errorf("expected non nil error, got nil error")
	}
	if d.width != 0 {
		t.Errorf("deck width expected to be %d, but got %d", 0, d.width)
	}
	if len(d.cards) != 0 {
		t.Errorf("deck expected to have no cards, but got %d", len(d.cards))
	}

}

func checkDeckValue(d *Deck, index, val int, t *testing.T) {
	card := d.cards[index].Value
	if card != val {
		t.Errorf("deck card at pos %d expected to be %d, but got %d", index, val, card)
	}
}

func checkCardValue(card Card, val int, visible bool, t *testing.T) {
	if card.Value != val || card.Visible != visible {
		t.Errorf("deck card expected to be {%d, %T}, but got {%d, %T}", val, visible, card.Value, card.Visible)
	}
}
