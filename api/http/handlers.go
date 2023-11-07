package http

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Blog struct
type Blog struct {
	ID            string
	Title         string
	ThumbnailPath string `dynamodbav:"thumbnail_path"`
	Uploaded      string
	Summary       string
	Content       template.HTML
}

type Server struct {
	config         aws.Config
	dynamoDBClient *dynamodb.Client
	s3Client       *s3.Client

	logger *logrus.Logger
}

func NewServer(region string, logger *logrus.Logger) *Server {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cfg.Region = region

	// Using the Config value, create the DynamoDB client
	dynamoDBClient := dynamodb.NewFromConfig(cfg)

	// Create an Amazon S3 service client
	s3Client := s3.NewFromConfig(cfg)

	return &Server{
		config:         cfg,
		dynamoDBClient: dynamoDBClient,
		s3Client:       s3Client,
		logger:         logger,
	}
}

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithContext(r.Context())

	var retBlogs []Blog
	out, err := s.dynamoDBClient.Scan(context.Background(), &dynamodb.ScanInput{
		TableName: aws.String("blogs"),
	})
	if err != nil {
		logger.Fatalf("scanning blogs dynamodb: %v", err)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &retBlogs)
	if err != nil {
		logger.Fatalf("unmarshalling blogs: %v", err)
	}

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "index.html", retBlogs); err != nil {
		logger.Fatalf("can't execute index template: %v", err)
	}
}

func (s *Server) About(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithContext(r.Context())

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "about.html", nil); err != nil {
		logger.Fatalf("can't execute about.html template: %v", err)
	}
}

func (s *Server) Show(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithContext(r.Context())

	title := mux.Vars(r)["title"]
	if title == "" {
		// TODO: redirect back to the index page.
		logger.Fatal("id is empty when showing blog")
	}

	// Using the Config value, create the DynamoDB client
	item, err := s.dynamoDBClient.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String("blogs"),
		Key: map[string]types.AttributeValue{
			"title": &types.AttributeValueMemberS{Value: title},
		},
	})
	if err != nil {
		logger.Fatalf("scanning blogs dynamodb: %v", err)
	}

	var blog Blog
	err = attributevalue.UnmarshalMap(item.Item, &blog)
	if err != nil {
		logger.Fatalf("unmarshalling blog: %v", err)
	}

	// Get the first page of results for ListObjectsV2 for a bucket
	object, err := s.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("warrenb95-blog"),
		Key:    aws.String(fmt.Sprintf("blogs/%s.md", title)),
	})
	if err != nil {
		logger.Fatal(err)
	}

	fbytes, err := io.ReadAll(object.Body)
	if err != nil {
		logger.Fatalf("can't read blog file: %v", err)
	}

	output := markdown.ToHTML(fbytes, nil, nil)
	blog.Content = template.HTML(string(output))

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "show.html", blog); err != nil {
		logger.Fatalf("can't execute show template: %v", err)
	}
}
