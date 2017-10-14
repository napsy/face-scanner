package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func serveHTTP(imgChan chan []interface{}) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, nil)
	})
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	fs := http.FileServer(http.Dir("static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("NEW CONN\n")
		for img := range imgChan {
			fmt.Printf("Sending image ...\n")
			d := struct {
				Type     string   `json:"type"`
				Filename string   `json:"filename"`
				Data     []string `json:"data"`
			}{"face", img[0].(string), img[1].([]string)}
			out, err := json.Marshal(d)
			if err != nil {
				log.Println(err)
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, out); err != nil {
				log.Println(err)
				return
			}
		}
	})
	http.ListenAndServe(":4000", nil)
}
