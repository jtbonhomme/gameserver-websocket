package games

import "fmt"

func (game *Game) revealCard(i, card int) error {
	if i >= game.decks[i].Len() || i < 0 {
		return fmt.Errorf("invalid index: %d", i)
	}

	ok, err := game.decks[i].IsVisible(card)
	if err != nil {
		return fmt.Errorf("error while checking card visibility: %w", err)
	}

	if ok {
		return fmt.Errorf("card is already visible")
	}

	err = game.decks[i].RevealCard(card)
	if err != nil {
		return fmt.Errorf("error while revealing card: %w", err)
	}

	return nil
}
