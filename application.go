package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	// AWS Elastic Beanstalk runs off port 5000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Save the logs here for AWS Elastic Beanstalk.
	if os.Getenv("ENV") == "PRODUCTION" {
		f, _ := os.Create("/var/log/golang/golang-server.log")
		defer f.Close()
		log.SetOutput(f)
	}

	// Server handlers.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index := template.Must(template.ParseGlob("./views/*"))
		if err := index.ExecuteTemplate(w, "index.html", nil); err != nil {
			log.Fatalf("can't execute index template: %v", err)
		}
	})

	log.Printf("Listening on port %s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
