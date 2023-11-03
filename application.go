package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
)

var (
	blogPath  = "public/blogs/"
	imagePath = "static/images/"
)

// Blog struct
type Blog struct {
	Title           string
	Content         string
	LastUpdated     string
	UpdatedDataTime time.Time
	ImagePath       string
}

func index(w http.ResponseWriter, r *http.Request) {
	blogDir, err := os.ReadDir(blogPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	var retBlogs []*Blog

	for _, de := range blogDir {
		if de.Type().IsRegular() {
			fname := strings.TrimSuffix(de.Name(), ".md")

			fbytes, err := os.ReadFile(filepath.Join(blogPath, de.Name()))
			if err != nil {
				log.Fatal(err.Error())
			}

			output := shrinkContent(fbytes, 400)

			finfo, err := de.Info()
			if err != nil {
				log.Fatal(err.Error())
			}

			lastUpdatedDuration := time.Now().UTC().Sub(finfo.ModTime().UTC())

			retBlogs = append(retBlogs, &Blog{
				Title:       fname,
				Content:     string(output),
				LastUpdated: durationToString(lastUpdatedDuration),
				ImagePath:   filepath.Join(imagePath, fmt.Sprintf("%s.png", fname)),
			})
		}
	}

	index := template.Must(template.ParseGlob("./views/*"))
	if err := index.ExecuteTemplate(w, "index.html", retBlogs); err != nil {
		log.Fatalf("can't execute index template: %v", err)
	}
}

func about(w http.ResponseWriter, r *http.Request) {
	index := template.Must(template.ParseGlob("./views/*"))
	if err := index.ExecuteTemplate(w, "about.html", nil); err != nil {
		log.Fatalf("can't execute about.html template: %v", err)
	}
}

// Show blog
// GET :id
// func (c *Controller) Show(ctx context.Context, id string) (blog *Blog, err error) {
// 	fbytes, err := os.ReadFile(filepath.Join(blogPath, id+".md"))
// 	if err != nil {
// 		return nil, err
// 	}

// 	output := markdown.ToHTML(fbytes, nil, nil)

// 	return &Blog{
// 		Title:   id,
// 		Content: string(output),
// 	}, nil
// }

func shrinkContent(content []byte, byteCount int) []byte {
	var shrunkContent []byte

	htmlContent := markdown.ToHTML(content, nil, nil)
	htmlReader := strings.NewReader(string(htmlContent))

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		return shrunkContent
	}

	var count int
	doc.Find("p").Each(
		func(i int, s *goquery.Selection) {
			if byteCount < count {
				return
			}

			shrunkContent = append(shrunkContent, []byte(s.Text())...)
			count += s.Length()
		},
	)

	return append(shrunkContent[:byteCount], []byte("...")...)
}

func durationToString(dur time.Duration) string {
	var ret string
	switch {
	case dur.Hours() > 24:
		days := int(dur.Round(time.Hour*24).Hours() / 24)
		ret = fmt.Sprintf("%dd", days)
	default:
		hours := int(dur.Round(time.Hour))
		minutes := int(dur.Minutes())
		ret = fmt.Sprintf("%dh %dm", hours, minutes)
	}

	return ret
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

	// Need to serve the static web content e.g. images at /static/assets/images/image.png.
	staticHandler := http.FileServer(http.Dir("assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Server handlers.
	http.HandleFunc("/", index)
	http.HandleFunc("/about", about)

	log.Printf("Listening on port %s\n\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
