package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	apihttp "github.com/warrenb95/website/api/http"
)

func main() {
	// AWS Elastic Beanstalk runs off port 5000.
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	// Save the logs here for AWS Elastic Beanstalk.
	if os.Getenv("ENV") == "PRODUCTION" {
		f, _ := os.Create("/var/log/blog.log")
		defer f.Close()
		log.SetOutput(f)
	}

	log := logrus.New()

	s := apihttp.NewServer("eu-west-2", log)

	r := mux.NewRouter()
	// Middleware.
	r.Use(s.Logger)

	// Server handlers.
	r.HandleFunc("/", s.Index)
	r.HandleFunc("/about", s.About)
	r.HandleFunc("/blog/{title}", s.Show)

	log.Printf("Listening on port %s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
