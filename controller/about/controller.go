package about

import (
	context "context"
)

// Controller for abouts
type Controller struct {
}

// About struct
type About struct {
	ID int `json:"id"`
}

// Index of abouts
// GET /about
func (c *Controller) Index(ctx context.Context) (abouts []*About, err error) {
	return []*About{}, nil
}
