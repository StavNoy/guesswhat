package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strings"
)

func main() {
	dir, _ := os.Getwd()

	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		for {
			msgType, input, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf(err.Error())

				return
			}

			output := []byte(strings.ToUpper(string(input)))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, output); err != nil {
				fmt.Printf(err.Error())

				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, dir + "/front/simple.html")
	})

	http.ListenAndServe(":8080", nil)
}
