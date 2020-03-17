package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func main() {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf(err.Error())

				return
			}

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				fmt.Printf(err.Error())

				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello\n")
	})

	http.ListenAndServe(":8080", nil)
}
