package main

import (
	"encoding/json"
	"fmt"
	. "github.com/gorilla/websocket"
	"net/http"
	"strings"
	"sync"
)

type Message struct {
	Message string `json:"message"`
}

func main() {
	conns := map[int]*Conn{}
	index := 0
	mutex := &sync.Mutex{}

	upgrader := Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)

			return
		}

		mutex.Lock()
		connI := index
		index++
		conns[connI] = conn
		mutex.Unlock()

		for {
			_, input, err := conn.ReadMessage()
			if err != nil {
				if IsCloseError(err, CloseGoingAway) {
					mutex.Lock()
					delete(conns, connI)
					mutex.Unlock()
				} else {
					fmt.Println(err)
				}

				return
			}

			dec := json.NewDecoder(strings.NewReader(string(input)))

			var inputMessage Message
			decErr := dec.Decode(&inputMessage)
			if decErr != nil {
				fmt.Println(err.Error())
			}

			output := Message{Message: strings.ToUpper(inputMessage.Message)}

			for _, aConn := range conns {
				if err = aConn.WriteJSON(output); err != nil {
					fmt.Println(err)

					return
				}
			}
		}
	})

	http.ListenAndServe(":8080", nil)
}
