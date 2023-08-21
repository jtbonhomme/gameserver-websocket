package games

import (
	"errors"

	"github.com/jtbonhomme/gameserver-websocket/internal/skyjo"
)

func (game *Game) drawCard() (skyjo.Card, error) {
	card, ok := game.drawPileCards.Pop()
	if !ok {
		return card, errors.New("too few cards in drawPile")
	}

	card.Visible = true
	game.discardPileCards.Push(card)

	return card, nil
}
