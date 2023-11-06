package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
)

// Blog struct
type Blog struct {
	ID            string
	Title         string
	ThumbnailPath string
	Uploaded      string
	Summary       string
	Content       template.HTML
}

func index(w http.ResponseWriter, r *http.Request) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cfg.Region = "eu-west-2"

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)
	var retBlogs []Blog
	out, err := svc.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String("blogs"),
	})
	if err != nil {
		log.Fatalf("scanning blogs dynamodb: %v", err)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &retBlogs)
	if err != nil {
		log.Fatalf("unmarshalling blogs: %v", err)
	}

	for _, b := range retBlogs {
		log.Println(b.ThumbnailPath)
	}

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "index.html", retBlogs); err != nil {
		log.Fatalf("can't execute index template: %v", err)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "about.html", nil); err != nil {
		log.Fatalf("can't execute about.html template: %v", err)
	}
}

// show, GET :id
func show(w http.ResponseWriter, r *http.Request) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cfg.Region = "eu-west-2"

	title := mux.Vars(r)["title"]
	if title == "" {
		// TODO: redirect back to the index page.
		log.Fatal("id is empty when showing blog")
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)
	item, err := svc.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String("blogs"),
		Key: map[string]types.AttributeValue{
			"title": &types.AttributeValueMemberS{Value: title},
		},
	})
	if err != nil {
		log.Fatalf("scanning blogs dynamodb: %v", err)
	}

	var blog Blog
	err = attributevalue.UnmarshalMap(item.Item, &blog)
	if err != nil {
		log.Fatalf("unmarshalling blog: %v", err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	object, err := client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("warrenb95-blog"),
		Key:    aws.String(fmt.Sprintf("blogs/%s.md", title)),
	})
	if err != nil {
		log.Fatal(err)
	}

	fbytes, err := io.ReadAll(object.Body)
	if err != nil {
		log.Fatalf("can't read blog file: %v", err)
	}

	output := markdown.ToHTML(fbytes, nil, nil)
	blog.Content = template.HTML(string(output))

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "show.html", blog); err != nil {
		log.Fatalf("can't execute show template: %v", err)
	}
}

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

	r := mux.NewRouter()

	// Need to serve the static web content e.g. images at /static/assets/images/image.png.
	staticHandler := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Server handlers.
	r.HandleFunc("/", index)
	r.HandleFunc("/about", about)
	r.HandleFunc("/blog/{title}", show)

	log.Printf("Listening on port %s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
