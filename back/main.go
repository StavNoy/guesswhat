package main

import (
	"encoding/json"
	"fmt"
	. "github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	users := map[int]User{}
	nextId := 0
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
		user := User{Conn: conn, Nickname: "", ID: nextId}
		nextId++
		users[user.ID] = user
		mutex.Unlock()

		for {
			_, input, err := conn.ReadMessage()
			if err != nil {
				if IsCloseError(err, CloseGoingAway) {
					mutex.Lock()
					delete(users, user.ID)
					mutex.Unlock()
				} else {
					fmt.Println(err)
				}

				return
			}

			inToOut(&user, input, users)
		}
	})

	http.ListenAndServe(":8080", nil)
}

type Message struct {
	Type string `json:"type"`
	Message string `json:"message"`
	Nickname string `json:"nickname"`
}

type OutMessage struct {
	Message string `json:"message"`
}

type User struct {
	Conn     *Conn
	Nickname string
	ID       int
}

func inToOut(from *User, input []byte, users map[int]User) {
	var msg Message
	json.NewDecoder(strings.NewReader(string(input))).Decode(&msg)

	if msg.Type == "message" {
		send(formatMessage(from, msg.Message), users)

		return
	}

	if msg.Type == "nickname" {
		output := formatMessage(from, "Changed nickname to \"" + msg.Nickname +"\"")
		from.Nickname = msg.Nickname
		send(output, users)
	}
}

func formatMessage(from *User, msg string) OutMessage {
	return OutMessage{ Message:from.Nickname + "(" + strconv.Itoa(from.ID) + "): " + strings.ToUpper(msg) }
}

func send(msg OutMessage, users map[int]User) {
	for _, user := range users {
		user.Conn.WriteJSON(msg)
	}
}