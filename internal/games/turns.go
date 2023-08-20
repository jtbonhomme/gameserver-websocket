package games

import (
	"fmt"
	"time"
)

// wait for all players to initialize
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
	if !game.started {
		return fmt.Errorf("[%s] game not started", game.Name)
	}

	if !game.playerAnswerMap[pID] {
		game.playerAnswerMap[pID] = true
		game.wg.Done()
	}

	return nil
}
