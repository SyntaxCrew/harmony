package harmony

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContext(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec := httptest.NewRecorder()

	ctx := newContext(rec, r)
	assert.IsType(t, &context{}, ctx)
}

func TestContext_JSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec := httptest.NewRecorder()

	if err := newContext(rec, r).JSON(http.StatusOK, testUser); assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, MIMEApplicationJSONCharsetUTF8, rec.Header().Get(HeaderContentType))
		assert.Equal(t, userJSON+"\n", rec.Body.String())
	}
}

func TestContext_SetAndGetPathParams(t *testing.T) {
	ctx := setPathParam()
	assert.Equal(t, map[string]string{"id": "1"}, ctx.PathParams())
	assert.NotPanics(t, func() {
		ctx.reset()
	})
}

func TestContext_PathParam(t *testing.T) {
	ctx := setPathParam()
	assert.Equal(t, "1", ctx.PathParam("id"))
}

func TestContext_PathParamInt(t *testing.T) {
	ctx := setPathParam()
	id, err := ctx.PathParamInt("id")
	assert.NoError(t, err)
	assert.Equal(t, 1, id)
}

func TestContext_QueryString(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?name=sujamess", nil)
	rec := httptest.NewRecorder()

	ctx := newContext(rec, r)
	assert.Equal(t, "sujamess", ctx.QueryString("name"))
}

func TestContext_QueryInt(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?id=1", nil)
	rec := httptest.NewRecorder()

	ctx := newContext(rec, r)
	id := ctx.QueryInt("id")
	assert.Equal(t, 1, id)
}

func TestContext_QueryFloat64(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?latitude=1.2345&longitude=1.1111", nil)
	rec := httptest.NewRecorder()

	ctx := newContext(rec, r)
	latitude := ctx.QueryFloat64("latitude")
	longitude := ctx.QueryFloat64("longitude")
	assert.Equal(t, 1.2345, latitude)
	assert.Equal(t, 1.1111, longitude)
}

func TestContext_QueryBool(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/?is_admin=true", nil)
	rec := httptest.NewRecorder()

	ctx := newContext(rec, r)
	isAdmin := ctx.QueryBool("is_admin")
	assert.True(t, isAdmin)
}

func TestContext_SendStatus(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	if err := newContext(rec, r).SendStatus(http.StatusOK); assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestContext_String(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	if err := newContext(rec, r).String(http.StatusOK, "OK"); assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, MIMETextPlainCharsetUTF8, rec.Header().Get(HeaderContentType))
		assert.Equal(t, "OK", rec.Body.String())
	}
}

func TestContext_Bind(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userJSON))
	rec := httptest.NewRecorder()

	var u user
	if err := newContext(rec, r).Bind(&u); assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, testUser, u)
	}
}

func newContext(w http.ResponseWriter, r *http.Request) Context {
	return NewContext(w, r, newBinder())
}

func setPathParam() Context {
	r := httptest.NewRequest(http.MethodGet, "/:id", nil)
	ctx := newContext(nil, r)
	ctx.SetPathParam("id", "1")
	return ctx
}
