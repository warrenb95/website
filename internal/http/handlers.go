package http

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	nhtml "golang.org/x/net/html"
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
		logger.Fatal(err)
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
		logger.WithError(err).Error("Scanning blogs dynamodb table")
		http.Error(w, "failed to list blogs", http.StatusInternalServerError)
		return
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &retBlogs)
	if err != nil {
		logger.WithError(err).Error("Failed to UnmarshalListOfMaps return blogs")
		http.Error(w, "failed to unmarshal return blogs from DB", http.StatusInternalServerError)
		return
	}

	sort.Slice(retBlogs, func(i, j int) bool {
		timeA, err := time.Parse("2006-01-02T15:04:05-07:00", retBlogs[i].Uploaded)
		if err != nil {
			logger.WithError(err).Error("Failed to parse time for blog")
			return false
		}

		timeB, err := time.Parse("2006-01-02T15:04:05-07:00", retBlogs[j].Uploaded)
		if err != nil {
			logger.WithError(err).Error("Failed to parse time for blog")
			return false
		}

		return !timeA.Before(timeB)
	})

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "index.html", retBlogs); err != nil {
		logger.WithError(err).Error("Failed to execute index template")
		http.Error(w, "failed to execute index template", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) About(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithContext(r.Context())

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "about.html", nil); err != nil {
		logger.WithError(err).Error("Failed to execute about template")
		http.Error(w, "failed to execute about template", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Show(w http.ResponseWriter, r *http.Request) {
	logger := s.logger.WithContext(r.Context())

	title := mux.Vars(r)["title"]
	if title == "" {
		// TODO: redirect back to the index page.
		logger.Warn("Empty blog title in show request")
		http.Error(w, "empty blog title", http.StatusBadRequest)
		return
	}
	logger = logger.WithField("title", title)

	// Using the Config value, create the DynamoDB client
	item, err := s.dynamoDBClient.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String("blogs"),
		Key: map[string]types.AttributeValue{
			"title": &types.AttributeValueMemberS{Value: title},
		},
	})
	if err != nil {
		logger.WithError(err).Error("Failed to get blog data")
		http.Error(w, "failed to get blog data", http.StatusInternalServerError)
		return
	}

	var blog Blog
	err = attributevalue.UnmarshalMap(item.Item, &blog)
	if err != nil {
		logger.WithError(err).Error("Failed unmarshal blog")
		http.Error(w, "failed to unmarshal blog data", http.StatusInternalServerError)
		return
	}

	// Get the first page of results for ListObjectsV2 for a bucket
	object, err := s.s3Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String("warrenb95-blog"),
		Key:    aws.String(fmt.Sprintf("blogs/%s.md", title)),
	})
	if err != nil {
		logger.WithError(err).Error("Failed to get blog content")
		http.Error(w, "failed to get blog content", http.StatusInternalServerError)
		return
	}

	fbytes, err := io.ReadAll(object.Body)
	if err != nil {
		logger.WithError(err).Error("Failed to read all blog content")
		http.Error(w, "failed to read blog content", http.StatusInternalServerError)
		return
	}

	output := markdown.ToHTML(fbytes, nil, nil)
	htmlContent := template.HTML(string(output))

	doc, err := nhtml.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		logger.WithError(err).Error("failed to parse html")
		http.Error(w, "failed to parse html", http.StatusInternalServerError)
		return
	}

	var imgClassAdder func(n *nhtml.Node)
	imgClassAdder = func(n *nhtml.Node) {
		if n.Type == nhtml.ElementNode && n.Data == "img" {
			n.Attr = append(n.Attr, nhtml.Attribute{
				Namespace: doc.Namespace,
				Key:       "class",
				Val:       "img-fluid",
			})
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			imgClassAdder(c)
		}
	}
	imgClassAdder(doc)

	var b bytes.Buffer
	err = nhtml.Render(&b, doc)
	if err != nil {
		logger.WithError(err).Error("Failed to render html to bytes")
		http.Error(w, "failed to render html", http.StatusInternalServerError)
		return
	}

	blog.Content = template.HTML(b.String())

	tmpl := template.Must(template.ParseGlob("./views/*"))
	if err := tmpl.ExecuteTemplate(w, "show.html", blog); err != nil {
		logger.WithError(err).Error("Failed to execute show template")
		http.Error(w, "failed to exeute show templated", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
