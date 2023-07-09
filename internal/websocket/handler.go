package websocket

import (
	"log"

	echo "github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"

	"github.com/jtbonhomme/pubsub/client"
)

func (s *Server) connect(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()

		client := client.New(
			s.nameGenerator.Generate(),
			ws,
		)

		// add this client into the list
		s.ps.AddClient(client)
		s.log.Printf("New Client %s - %s is connected, total clients %d and subscriptions %d", client.Name, client.ID, len(s.ps.Clients), len(s.ps.Subscriptions))

		for {
			var err error
			// Read
			var msg = make([]byte, 512)
			err = websocket.Message.Receive(ws, &msg)
			if err != nil {
				s.ps.RemoveClient(client)
				c.Logger().Error(err)
				log.Printf("removed client %s total clients %d and subscriptions %d", client.Name, len(s.ps.Clients), len(s.ps.Subscriptions))
				return
			}

			s.ps.HandleReceiveMessage(client, msg)
		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
