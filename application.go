package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	handler "github.com/warrenb95/website/internal/http"
)

func main() {
	// AWS Elastic Beanstalk runs off port 5000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	// Save the logs here for AWS Elastic Beanstalk.
	if os.Getenv("ENV") == "PRODUCTION" {
		f, _ := os.Create("/var/log/blog.log")
		defer f.Close()
		log.SetOutput(f)
	}

	log := logrus.New()

	s := handler.NewServer("eu-west-2", log)

	r := mux.NewRouter()
	// Middleware.
	r.Use(s.Logger)

	fs := http.FileServer(http.Dir("assets"))
	// r.Handle("/static/", http.StripPrefix("/static/", fs))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Server handlers.
	r.HandleFunc("/", s.Index)
	r.HandleFunc("/about", s.About)
	r.HandleFunc("/blog/{title}", s.Show)

	log.Printf("Listening on port %s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
