package controller

import (
	context "context"
	"os"
	"path/filepath"
	"strings"

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
	Title   string `json:"title"`
	Content string `json:"content"`
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

			output := shrinkContent(fbytes, 200)

			retBlogs = append(retBlogs, &Blog{
				Title:   fname,
				Content: string(output),
			})
		}
	}

	return retBlogs, nil
}

// Show blog
// GET :id
func (c *Controller) Show(ctx context.Context, title string) (blog *Blog, err error) {
	fbytes, err := os.ReadFile(filepath.Join(blogPath, title))
	if err != nil {
		return nil, err
	}

	output := markdown.ToHTML(fbytes, nil, nil)

	return &Blog{
		Title:   title,
		Content: string(output),
	}, nil
}

func shrinkContent(content []byte, byteCount int) []byte {
	var shrunkContent []byte

	htmlContent := markdown.ToHTML(content, nil, nil)
	htmlReader := strings.NewReader(string(htmlContent))

	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		// TODO: handle error
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
