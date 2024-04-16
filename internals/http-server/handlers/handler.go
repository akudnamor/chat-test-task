package handlers

import (
	"chat-backend-final/internals/storage"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"time"
)

func HandlerWS(st *storage.Storage, socket *socketio.Server) func(http.ResponseWriter, *http.Request) {
	socket.OnEvent("/", "message", func(s socketio.Conn, msg storage.Message) error {

		msg.CreatedAt = time.Now()
		id, err := st.AddMessage(msg)
		msg.ID = id
		if err != nil {
			log.Println("failed to add message in DB", err)
			return err
		}
		socket.BroadcastToNamespace("/", "message", msg)

		return nil
	})
	socket.OnConnect("/", func(s socketio.Conn) error {
		log.Println("connected id:", s.ID())
		return nil
	})
	return func(w http.ResponseWriter, r *http.Request) {
		socket.ServeHTTP(w, r)
	}
}
