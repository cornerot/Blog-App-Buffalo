package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/pop"
	"github.com/mikaelm1/blog_app/models"
	"github.com/pkg/errors"
)

// PostsIndex default implementation.
func PostsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	posts := &models.Posts{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// Retrieve all Posts from the DB
	if err := q.All(posts); err != nil {
		return errors.WithStack(err)
	}
	// Make Users available inside the html template
	c.Set("posts", posts)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("posts/index.html"))
}

func PostsCreateGet(c buffalo.Context) error {
	c.Set("post", &models.Post{})
	return c.Render(200, r.HTML("posts/create"))
}

func PostsCreatePost(c buffalo.Context) error {
	// Allocate an empty Post
	post := &models.Post{}
	user := c.Value("current_user").(*models.User)
	// Bind post to the html form elements
	if err := c.Bind(post); err != nil {
		return errors.WithStack(err)
	}
	// Get the DB connection from the context
	tx := c.Value("tx").(*pop.Connection)
	// Validate the data from the html form
	post.AuthorID = user.ID
	verrs, err := tx.ValidateAndCreate(post)
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("post", post)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("posts/create"))
	}
	// If there are no errors set a success message
	c.Flash().Add("success", "New post added successfully.")
	// and redirect to the users index page
	return c.Redirect(302, "/")
}