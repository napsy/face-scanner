package main

import (
	"html/template"
	"net/http"
)

func serveHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("resources/app/index.html")
		if err != nil {
			panic(err)
		}
		t.Execute(w, nil)
	})
	http.ListenAndServe(":4000", nil)
}
