package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"log"
	"net/http"

	b64 "encoding/base64"
)

var userConns map[string]*websocket.Conn
var userMsgs map[string][]string

func main() {
	router := chi.NewRouter()
	userConns = make(map[string]*websocket.Conn)
	router.Get("/read", ReadMsg)
	httpPort := "8081"
	log.Fatal(http.ListenAndServe(":"+httpPort, router))
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func ReadMsg(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	userConns[r.URL.Query().Get("in_id")] = conn
	go func() {
		for {
			msgType, msg, err := userConns[r.URL.Query().Get("in_id")].ReadMessage()
			if err != nil {
				return
			}
			sEnc := b64.StdEncoding.EncodeToString(msg)

			fmt.Printf("%s Recieved: %s Endoded Message: %s \n", userConns[r.URL.Query().Get("in_id")].RemoteAddr(), string(msg), sEnc)
			out, ok := userConns[r.URL.Query().Get("out_id")]
			if ok {
				if err := out.WriteMessage(msgType, []byte(sEnc)); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				userMsgs[r.URL.Query().Get("out_id")] = append(userMsgs[r.URL.Query().Get("out_id")], sEnc)
			}
		}
	}()
}
