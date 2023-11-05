package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gorilla/mux"
)

var (
	s3ImageAddr = "https://warrenb95-blog.s3.eu-west-2.amazonaws.com/blogs/images/"
)

// Blog struct
type Blog struct {
	ID            string
	Title         string
	ThumbnailPath string
	Uploaded      string
	Summary       string
}

func index(w http.ResponseWriter, r *http.Request) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	// Get the first page of results for ListObjectsV2 for a bucket
	output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String("warrenb95-blog"),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("first page results:")
	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}

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
		b.ThumbnailPath = filepath.Join(s3ImageAddr, b.Title+".png")
		log.Println(b.ThumbnailPath)
		log.Println(b.Title)
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
// func show(w http.ResponseWriter, r *http.Request) {
// 	// Load the Shared AWS Configuration (~/.aws/config)
// 	// cfg, err := config.LoadDefaultConfig(context.TODO())
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// // Create an Amazon S3 service client
// 	// client := s3.NewFromConfig(cfg)

// 	// // Get the first page of results for ListObjectsV2 for a bucket
// 	// output, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
// 	// 	Bucket: aws.String("warrenb95-blog"),
// 	// })
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	title := mux.Vars(r)["title"]
// 	if title == "" {
// 		// TODO: redirect back to the index page.
// 		log.Fatal("id is empty when showing blog")
// 	}

// 	fbytes, err := os.ReadFile(filepath.Join(blogPath, title+".md"))
// 	if err != nil {
// 		log.Fatalf("can't read blog file: %v", err)
// 	}

// 	output := markdown.ToHTML(fbytes, nil, nil)
// 	blog := &Blog{
// 		Title:   title,
// 		Content: template.HTML(string(output)),
// 	}

// 	tmpl := template.Must(template.ParseGlob("./views/*"))
// 	if err := tmpl.ExecuteTemplate(w, "show.html", blog); err != nil {
// 		log.Fatalf("can't execute show template: %v", err)
// 	}
// }

// func shrinkContent(content []byte, byteCount int) []byte {
// 	var shrunkContent []byte

// 	htmlContent := markdown.ToHTML(content, nil, nil)
// 	htmlReader := strings.NewReader(string(htmlContent))

// 	doc, err := goquery.NewDocumentFromReader(htmlReader)
// 	if err != nil {
// 		return shrunkContent
// 	}

// 	var count int
// 	doc.Find("p").Each(
// 		func(i int, s *goquery.Selection) {
// 			if byteCount < count {
// 				return
// 			}

// 			shrunkContent = append(shrunkContent, []byte(s.Text())...)
// 			count += s.Length()
// 		},
// 	)

// 	return append(shrunkContent[:byteCount], []byte("...")...)
// }

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
	r.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Server handlers.
	r.HandleFunc("/", index)
	r.HandleFunc("/about", about)
	// r.HandleFunc("/blog/{id}", show)

	log.Printf("Listening on port %s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
