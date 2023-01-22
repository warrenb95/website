package controller

import (
	context "context"
	"os"
	"path/filepath"

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
			retBlogs = append(retBlogs, &Blog{
				Title: de.Name(),
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
