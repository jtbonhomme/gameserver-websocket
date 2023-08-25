package games

import (
	"errors"
	"fmt"
	"time"

	"github.com/jtbonhomme/gameserver-websocket/internal/skyjo"
	"github.com/jtbonhomme/gameserver-websocket/internal/utils"
)

const (
	maxScore int = 100
)

// turnsLoop runs turns in a loop.
func (game *Game) turnsLoop() error {
	var gameIsOver bool

	for !gameIsOver {
		var err error

		_, err = game.newTurn()
		if err != nil {
			return fmt.Errorf("error reseting game: %s", err.Error())
		}
		/*
			var turnIsOver bool
			for !turnIsOver {
			}
		*/
	}

	return nil
}

// newTurn initializes a new turn.
func (game *Game) newTurn() (int, error) {
	var err error
	var firstPlayer int = -1
	var maxScore int = -1

	game.log.Debug().Msgf("[%s] new turn ...", game.Name)
	game.drawPileCards = utils.Stack[skyjo.Card]{}
	for _, c := range skyjo.GenerateCards() {
		game.drawPileCards.Push(c)
	}
	game.discardPileCards = utils.Stack[skyjo.Card]{}

	// reset player's decks
	game.decks = []*skyjo.Deck{}
	for i := 0; i < len(game.players); i++ {
		game.decks = append(game.decks, skyjo.NewDeck())
	}

	// deal 12 cards to each player
	for i := 0; i < skyjo.CardsPerPlayer; i++ {
		for j, p := range game.players {
			card, ok := game.drawPileCards.Pop()
			if !ok {
				return firstPlayer, errors.New("too few cards in drawPile")
			}

			err = game.decks[j].AddCard(card)
			if err != nil {
				return firstPlayer, fmt.Errorf("error adding card %v to deck of %s: %s", card, p.Name, err.Error())
			}
		}
	}

	// ask to each player to reveal 2 cards
	// publication to all clients who subscribed to a channel
	game.log.Debug().Msgf("[%s] publish a call to rpc: revealTwoCards", game.Name)
	_, err = game.publish(`{"type": "rpc", "emitter": "game", "id": "` + game.ID.String() + `", "data": "revealTwoCards"}`)
	if err != nil {
		game.log.Error().Msgf("[%s] publication error: %s", game.Name, err.Error())
	}

	err = game.waitPlayersRevealCards(2)
	if err != nil {
		return firstPlayer, fmt.Errorf("error initializing players: %s", err.Error())
	}

	// decide first player
	for i, _ := range game.players {
		if game.decks[i].Value() > maxScore {
			maxScore = game.decks[i].Value()
			firstPlayer = i
		}
	}

	// return first card from drawPile
	_, err = game.drawCard()
	if err != nil {
		return firstPlayer, fmt.Errorf("error drawing first card: %s", err.Error())
	}

	return firstPlayer, nil
}

// waitPlayersRevealCards waits for all players to reveal n cards.
func (game *Game) waitPlayersRevealCards(n int) error {
	var done = make(chan struct{})

	game.playerAnswerMap = make(map[string]int)
	for _, player := range game.players {
		game.wg.Add(n)
		game.playerAnswerMap[player.ID.String()] = 0
	}

	// Will automatically close done channel and end select after all RPC have been called.
	go func() {
		defer close(done)
		game.wg.Wait()
		game.log.Debug().Msgf("[%s] ... done !", game.Name)
	}()

	// Blocks until done channel is closed, or timeout occurs.
	game.log.Debug().Msgf("[%s] wait for all players to send rpc reveal two cards twice ...", game.Name)
	select {
	case <-done:
		return nil
	case <-time.After(game.waitForRPCTimeout):
		return fmt.Errorf("wait ended after timeout")
	}
}

// RevealCard turns a card visible from a player's deck.
func (game *Game) RevealCard(pID string, card int) error {
	game.log.Debug().Msgf("[%s] player %s reveal card %d", game.Name, pID, card)
	if !game.Started {
		return fmt.Errorf("[%s] game not started", game.Name)
	}

	i, err := game.playerByID(pID)
	if err != nil {
		return fmt.Errorf("unknown player id: %s", err.Error())
	}

	err = game.revealCard(i, card)
	if err != nil {
		return fmt.Errorf("error revealing card %d on deck %d: %s", card, i, err.Error())
	}

	count, ok := game.playerAnswerMap[pID]
	if !ok {
		return fmt.Errorf("unknown player id")
	}

	game.playerAnswerMap[pID] = count + 1
	game.wg.Done()

	return nil
}
