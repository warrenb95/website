package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	logger := log.Default()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index := template.Must(template.ParseFiles("./public/views/index.html"))
		if err := index.Execute(w, nil); err != nil {
			logger.Fatalf("can't execute index template: %v", err)
		}
	})

	logger.Println("listing on 0.0.0.0:8080")
	logger.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
