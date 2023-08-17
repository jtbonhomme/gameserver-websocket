package manager_test

import (
	"encoding/json"
	"testing"

	"github.com/centrifugal/centrifuge"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/jtbonhomme/gameserver-websocket/internal/manager"
	"github.com/jtbonhomme/gameserver-websocket/internal/players"
	"github.com/jtbonhomme/gameserver-websocket/internal/storage/memory"
)

type Response struct {
	Status string `json:"status"`
	Result string `json:"result"`
}

func TestMemoryPlayer(t *testing.T) {
	var response Response

	logger := zerolog.Nop()
	// concrete memory test storage implementation
	s := memory.New(&logger)

	m := manager.New(&logger, s)

	replyChan := make(chan []byte)
	errChan := make(chan error)

	// first time player registration
	go func() {
		m.RegisterPlayer([]byte(`{"name":"name1"}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	r1 := <-replyChan
	e1 := <-errChan
	if e1 != nil {
		t.Errorf("error while registering name1: %s", e1.Error())
	}

	err := json.Unmarshal(r1, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(r1), err.Error())
	}

	var player1 players.Player
	err = json.Unmarshal([]byte(response.Result), &player1)
	if err != nil {
		t.Errorf("error while unmarshaling result %q: %s", response.Result, err.Error())
	}

	id := player1.ID.String()
	if id == uuid.Nil.String() {
		t.Errorf("expected non nil uuid for player name1")
	}

	go func() {
		m.RegisterPlayer([]byte(`{"name":"name2"}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	r2 := <-replyChan
	e2 := <-errChan
	if e2 != nil {
		t.Errorf("error while registering name2: %s", e2.Error())
	}

	err = json.Unmarshal(r2, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(r2), err.Error())
	}

	var player2 players.Player
	err = json.Unmarshal([]byte(response.Result), &player2)
	if err != nil {
		t.Errorf("error while unmarshalingresult  %q: %s", response.Result, err.Error())
	}

	if player2.ID.String() == "" {
		t.Errorf("error, expected new UUID created for player name2")
	}
	if player2.ID.String() == id {
		t.Errorf("expected player name2 uuid to be different from player name1")
	}

	go func() {
		m.RegisterPlayer([]byte(`{"id": "`+id+`", "name":"name1"}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	r3 := <-replyChan
	e3 := <-errChan
	if e3 != nil {
		t.Errorf("error while registering name1: %s", e3.Error())
	}

	err = json.Unmarshal(r3, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(r3), err.Error())
	}

	var player3 players.Player
	err = json.Unmarshal([]byte(response.Result), &player3)
	if err != nil {
		t.Errorf("error while unmarshaling result %q: %s", response.Result, err.Error())
	}

	if player3.ID.String() == "" {
		t.Errorf("error, expected UUID not nil for player name1")
	}

	if player3.ID.String() != id {
		t.Errorf("error, expected no new UUID was created for player name1")
	}

	go func() {
		m.ListPlayers([]byte(`{}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	rAll := <-replyChan
	eAll := <-errChan
	if eAll != nil {
		t.Errorf("error while registering name1: %s", eAll.Error())
	}

	err = json.Unmarshal(rAll, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(rAll), err.Error())
	}

	var players1 []players.Player
	err = json.Unmarshal([]byte(response.Result), &players1)
	if err != nil {
		t.Errorf("error while unmarshaling result %q: %s", response.Result, err.Error())
	}

	if len(players1) != 2 {
		t.Errorf("expected 2 registered players")
	}

	for _, p := range players1 {
		if p.ID != uuid.Nil {
			t.Errorf("error, expected UUID returned by listAll to nil")
		}
	}

	go func() {
		m.UnregisterPlayer([]byte(`{"id": "`+id+`"}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	r4 := <-replyChan
	e4 := <-errChan
	if e3 != nil {
		t.Errorf("error while registering name1: %s", e4.Error())
	}

	err = json.Unmarshal(r4, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(r4), err.Error())
	}

	go func() {
		m.ListPlayers([]byte(`{}`), func(r centrifuge.RPCReply, e error) {
			replyChan <- r.Data
			errChan <- e
		})
	}()

	rAll2 := <-replyChan
	eAll2 := <-errChan
	if eAll != nil {
		t.Errorf("error while registering name1: %s", eAll2.Error())
	}

	err = json.Unmarshal(rAll2, &response)
	if err != nil {
		t.Errorf("error while unmarshaling response %q: %s", string(rAll2), err.Error())
	}

	var players2 []players.Player
	err = json.Unmarshal([]byte(response.Result), &players2)
	if err != nil {
		t.Errorf("error while unmarshaling result %q: %s", response.Result, err.Error())
	}

	if len(players2) != 1 {
		t.Fatal("expected 1 registered players")
	}

	if players2[0].Name != "name2" {
		t.Errorf("expected last registered player to be name2 but got %s", players2[0].Name)
	}
}
