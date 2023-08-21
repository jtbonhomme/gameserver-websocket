package skyjo_test

import (
	"testing"

	"github.com/jtbonhomme/gameserver-websocket/internal/skyjo"
)

func TestAllCards(t *testing.T) {
	cards := skyjo.GenerateCards()

	check := make(map[int]int)

	if len(cards) != skyjo.TotalCards {
		t.Errorf("expected cards to be %d long, and got %d", skyjo.TotalCards, len(cards))
	}

	for i := 0; i < len(cards); i++ {
		n, ok := check[cards[i].Value]
		if !ok {
			check[cards[i].Value] = 1
		} else {
			check[cards[i].Value] = n + 1
		}
	}

	for k, v := range check {
		switch k {
		case skyjo.MinusTwo:
			if v != skyjo.MinusTwoNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.MinusOne:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Zero:
			if v != skyjo.ZeroNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.One:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Two:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Three:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Four:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Five:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Six:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Seven:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Eight:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Nine:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Ten:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Eleven:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		case skyjo.Twelve:
			if v != skyjo.DefaultCardsNumber {
				t.Errorf("expected %d cards, got %d", skyjo.DefaultCardsNumber, v)
			}
		}
	}
}
