package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Message string
}

func main() {
	dir, _ := os.Getwd()

	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)

		for {
			_, input, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf(err.Error())

				return
			}

			dec := json.NewDecoder(strings.NewReader(string(input)))

			var inputMessage Message
			decErr := dec.Decode(&inputMessage)
			if decErr != nil {
				fmt.Println(err.Error())
			}

			output := Message{Message:strings.ToUpper(inputMessage.Message)}

			if err = conn.WriteJSON(output); err != nil {
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
