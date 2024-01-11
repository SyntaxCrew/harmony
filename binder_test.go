package harmony

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type (
	testBindUser struct {
		Username string `path:"username"`
		IsActive string `query:"is_active"`
		Title    string `json:"title"`
	}
)

func TestBinder_Bind(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/users/:username/post?is_active=true", strings.NewReader(`{"title":"Hello, Harmony!"}`))
	rec := httptest.NewRecorder()

	ctx := newContext(rec, r)
	ctx.SetPathParam("username", "sujamess")
	var u testBindUser
	if err := ctx.Bind(&u); assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "sujamess", u.Username)
		assert.Equal(t, "Hello, Harmony!", u.Title)
		assert.Equal(t, "true", u.IsActive)
	}
}
