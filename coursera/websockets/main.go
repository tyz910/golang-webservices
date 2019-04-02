package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/icrowley/fake"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func sendNewMsgNotifications(client *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	for {
		w, err := client.NextWriter(websocket.TextMessage)
		if err != nil {
			ticker.Stop()
			break
		}

		msg := newMessage()
		w.Write(msg)
		w.Close()

		<-ticker.C
	}
}

func main() {
	tmpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}
		go sendNewMsgNotifications(ws)
	})

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

func newMessage() []byte {
	data, _ := json.Marshal(map[string]string{
		"email":   fake.EmailAddress(),
		"name":    fake.FirstName() + " " + fake.LastName(),
		"subject": fake.Product() + " " + fake.Model(),
	})
	return data
}
