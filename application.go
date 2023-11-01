package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	logger := log.Default()

	// AWS Elastic Beanstalk runs off port 5000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Save the logs here for AWS Elastic Beanstalk.
	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer f.Close()
	log.SetOutput(f)

	// Server handlers.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index := template.Must(template.ParseFiles("./public/views/index.html"))
		if err := index.Execute(w, nil); err != nil {
			logger.Fatalf("can't execute index template: %v", err)
		}
	})

	logger.Printf("Listening on port %s\n\n", port)
	logger.Fatal(http.ListenAndServe(":"+port, nil))
}
