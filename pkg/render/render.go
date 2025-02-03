package render

import (
	"context"
	"net/http"

	"github.com/a-h/templ"
)

func WrapHandler(c templ.Component) *templ.ComponentHandler {
	return templ.Handler(c)
}

func Render(w http.ResponseWriter, c templ.Component) error {
	return c.Render(context.Background(), w)
}
