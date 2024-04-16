package handlers

import (
	"chat-backend-final/internals/storage"
	"encoding/json"
	"errors"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func HandlerWS(st *storage.Storage, socket *socketio.Server) func(http.ResponseWriter, *http.Request) {
	socket.OnEvent("/", "message", func(s socketio.Conn, msg storage.Message) error {
		if strings.Contains(msg.Text, "testWordForValidation") {
			return errors.New("incorrect word")
		}

		// TODO: add some other checks for input validation

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

func HandlerAPI(st *storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS for api
		w.Header().Set("Access-Control-Allow-Origin", "*")
		skip, err := strconv.Atoi(r.URL.Query().Get("skip"))
		if err != nil {
			log.Println("failed to Atoi `skip`", err)
			return
		}

		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			log.Println("failed to Atoi `limit`", err)
			return
		}

		msgs, err := st.GetMessages(limit, skip)
		if err != nil {
			log.Println("failed to get messages from DB", err)
			return
		}

		// TODO: use fast JSON
		msgsJson, err := json.Marshal(msgs)
		if err != nil {
			log.Println("failed to marshal", err)
			return
		}

		_, err = w.Write(msgsJson)
		if err != nil {
			log.Println("failed to send msgs")
			return
		}
		return
	}
}
