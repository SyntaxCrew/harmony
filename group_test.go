package harmony

import (
	"net/http"
	"testing"
)

func TestGroup_Group(t *testing.T) {
	app := New()
	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/users", writeStringOKHandler())
	testMethod(t, http.MethodGet, "/api/v1/users", app)
}
