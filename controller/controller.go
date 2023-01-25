package controller

import (
	context "context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
)

const (
	blogPath = "./public/blogs/"
)

// Controller for blogs
type Controller struct {
}

// Blog struct
type Blog struct {
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	LastUpdated     string    `json:"last_updated"`
	UpdatedDataTime time.Time `json:"updated_date_time"`
	ImagePath       string    `json:"image_path"`
}

// Index of blogs
// GET
func (c *Controller) Index(ctx context.Context) (blogs []*Blog, err error) {
	blogDir, err := os.ReadDir(blogPath)
	if err != nil {
		return nil, err
	}

	var retBlogs []*Blog

	for _, de := range blogDir {
		if de.Type().IsRegular() {
			fname := strings.TrimSuffix(de.Name(), ".md")

			fbytes, err := os.ReadFile(filepath.Join(blogPath, de.Name()))
			if err != nil {
				return nil, err
			}

			output := shrinkContent(fbytes, 400)

			finfo, err := de.Info()
			if err != nil {
				return nil, err
			}

			lastUpdatedDuration := time.Now().UTC().Sub(finfo.ModTime().UTC())

			retBlogs = append(retBlogs, &Blog{
				Title:       fname,
				Content:     string(output),
				LastUpdated: durationToString(lastUpdatedDuration),
				ImagePath:   fmt.Sprintf("/images/%s.png", fname),
			})
		}
	}

	return retBlogs, nil
}

// Show blog
// GET :id
func (c *Controller) Show(ctx context.Context, id string) (blog *Blog, err error) {
	fbytes, err := os.ReadFile(filepath.Join(blogPath, id+".md"))
	if err != nil {
		return nil, err
	}

	output := markdown.ToHTML(fbytes, nil, nil)

	return &Blog{
		Title:   id,
		Content: string(output),
	}, nil
}

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
		ret = dur.String()
	}

	return ret
}
