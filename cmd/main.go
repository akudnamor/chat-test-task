package main

import (
	"chat-backend-final/internals/http-server/handlers"
	"chat-backend-final/internals/storage"
	"github.com/go-chi/chi/v5"
	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
	"log"
	"net/http"
)

// TODO: -config

func main() {
	r := chi.NewRouter()
	st, err := storage.New("storage/storage.db")
	if err != nil {
		log.Println("failed to open/create storage", err)
		return
	}
	socket := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&websocket.Transport{
				CheckOrigin: func(r *http.Request) bool {
					// CORS settings.
					// TODO: add allow-origin for your host
					return true
				},
			},
		},
	})

	go socket.Serve()
	defer socket.Close()

	r.HandleFunc("/socket.io/", handlers.HandlerWS(st, socket))
	server := http.Server{
		Addr:    ":8000",
		Handler: r,
	}
	log.Println("Server start on..", server.Addr)
	log.Fatal(server.ListenAndServe())
}
