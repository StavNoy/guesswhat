package main

import (
	"encoding/json"
	"fmt"
	. "github.com/gorilla/websocket"
	"net/http"
	"strings"
)

type Message struct {
	Message string `json:"message"`
}


func main() {
	var conns []*Conn

	upgrader := Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	http.HandleFunc("/", handleConnect(upgrader, conns))

	http.ListenAndServe(":8080", nil)
}

func handleConnect(upgrader Upgrader, conns []*Conn) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println(err)

			return
		}
		conns = append(conns, conn)
		connI := len(conns) - 1

		for {
			_, input, err := conn.ReadMessage()
			if err != nil {
				if IsCloseError(err, CloseGoingAway) {
					conns = append(conns[:connI], conns[connI+1:]...)
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
	}
}
