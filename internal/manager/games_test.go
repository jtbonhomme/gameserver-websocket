package manager_test

import (
	"encoding/json"
	"testing"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"

	"github.com/jtbonhomme/gameserver-websocket/internal/games"
)

func TestGames(t *testing.T) {
	var response Response

	replyChan := make(chan []byte)
	errChan := make(chan error)

	// first time player registration
	go func() {
		mgr.CreateGame([]byte(`{"minPlayers":1, "maxPlayers": 4}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	r1 := <-replyChan
	e1 := <-errChan
	if e1 != nil {
		t.Errorf("error while creating game: %s", e1.Error())
	}

	err := json.Unmarshal(r1, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(r1), err.Error())
	}

	var game games.Game
	err = json.Unmarshal([]byte(response.Result), &game)
	if err != nil {
		t.Errorf("error while unmarshaling result %q: %s", response.Result, err.Error())
	}

	id := game.ID.String()
	if id == uuid.Nil.String() {
		t.Errorf("expected non nil uuid for player name1")
	}
}
