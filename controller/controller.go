package controller

import (
	context "context"
)

// Controller for blogs
type Controller struct {
}

// Blog struct
type Blog struct {
	ID int `json:"id"`
}

// Index of blogs
// GET
func (c *Controller) Index(ctx context.Context) (blogs []*Blog, err error) {
	return []*Blog{}, nil
}

// Show blog
// GET :id
func (c *Controller) Show(ctx context.Context, id int) (blog *Blog, err error) {
	return &Blog{
		ID: id,
	}, nil
}
