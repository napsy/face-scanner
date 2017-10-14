package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func serveHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("resources/app/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, nil)
	})
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			return err
		}
	})
	http.ListenAndServe(":4000", nil)
}
