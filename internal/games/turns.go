package games

import (
	"fmt"
	"time"
)

// turnsLoop runs turns in a loop.
func (game *Game) turnsLoop() error {
	var gameIsOver bool

	for !gameIsOver {
		var err error

		_, err = game.newTurn()
		if err != nil {
			return fmt.Errorf("error reseting game: %w", err)
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
	return 0, nil
}

// waitAllPlayersInitialized waits for all players to initialize.
func (game *Game) waitAllPlayersInitialized() {
	var done = make(chan struct{})

	game.log.Info().Msgf("[%s] wait all players to initialize", game.Name)

	game.playerAnswerMap = make(map[string]bool)
	for _, pID := range game.Players() {
		game.wg.Add(1)
		game.playerAnswerMap[pID] = false
	}

	// Will automatically close done channel and end select after all RPC have been called.
	go func() {
		game.log.Debug().Msgf("[%s] waiting...", game.Name)
		defer close(done)
		game.wg.Wait()
		game.log.Debug().Msgf("[%s] WaitGroup done !", game.Name)
	}()

	// Blocks until done channel is closed, or timeout occurs.
	game.log.Debug().Msgf("wait for InitRPC")
	select {
	case <-done:
		game.log.Debug().Msgf("channel closed")
		game.log.Debug().Msgf("exit waitForInitRPC")
		go game.startTurnLoop()
		return
	case <-time.After(game.waitForRPCTimeout):
		game.log.Debug().Msgf("** timeout **")
		game.log.Debug().Msgf("exit waitForInitRPC")
		// TODO handle game termination
	}

}

func (game *Game) startTurnLoop() {
	game.log.Info().Msgf("[%s] enter turn loop", game.Name)
	// TODO: here, initialize turn
}

// PlayerInit -.
func (game *Game) PlayerInit(pID string) error {
	if !game.Started {
		return fmt.Errorf("[%s] game not Started", game.Name)
	}

	if !game.playerAnswerMap[pID] {
		game.playerAnswerMap[pID] = true
		game.wg.Done()
	}

	return nil
}

func (game *Game) resetTurn() (int, error) {
	var firstPlayer int
	/*var maxScore int

	g.drawPileCards = utils.Stack[skyjo.Card]{}
	for _, c := range skyjo.GenerateCards() {
		g.drawPileCards.Push(c)
	}
	g.discardPileCards = utils.Stack[skyjo.Card]{}

	// reset player's decks
	for _, p := range g.players {
		p.ResetDeck()
	}

	// deal 12 cards to each player
	for i := 0; i < skyjo.CardsPerPlayer; i++ {
		for _, p := range g.players {
			card, ok := g.drawPileCards.Pop()
			if !ok {
				return firstPlayer, errors.New("too few cards in drawPile")
			}

			err := p.AddCard(card)
			if err != nil {
				return firstPlayer, fmt.Errorf("error adding card %v to deck of %s: %s", card, p.Name(), err.Error())
			}
		}
	}

	// each player reveal 2 cards
	for i, p := range g.players {
		err := p.StartTurn()
		if err != nil {
			return firstPlayer, fmt.Errorf("error initializing player %d: %w", i, err)
		}
	}

	// decide first player
	maxScore = -1
	firstPlayer = -1

	for i, p := range g.players {
		if p.DeckValue() > maxScore {
			maxScore = p.DeckValue()
			firstPlayer = i
		}
	}

	// return first card from drawPile
	_, err := g.drawCard()
	if err != nil {
		return firstPlayer, fmt.Errorf("error drawing first card: %s", err.Error())
	}
	*/
	return firstPlayer, nil
}
