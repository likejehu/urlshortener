package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/likejehu/urlshortener/db"
	"github.com/likejehu/urlshortener/handlers/mocks"
	"github.com/likejehu/urlshortener/models"
	"github.com/stretchr/testify/assert"
)

var (
	err         = db.Error404
	longURLJSON = `{"longUrl":"https://golang.org/doc/effective_go.html"}`
	id          = "f202ad3f"
	url         = "https://golang.org/doc/effective_go.html"
	testmineURL = &models.MineURL{LongURL: "https://golang.org/doc/effective_go.html", ShortURL: "http://localhost:8080/f202ad3f"}
	idt         = &TestID{}
	link        = "\"http://localhost:8080/f202ad3f\"\n"
	err500      = errors.New("something bad happened")
)

// TestID is an implementation of the ider Interface
type TestID struct {
	SI string
}

// NewID  retuns new short id
func (s *TestID) NewID() string {
	s.SI = id
	return s.SI
}

func TestCreateLink(t *testing.T) {

	t.Run("succes case", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURLJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/createlink")
		mockStore := &mocks.Storer{}
		handler := &Handler{mockStore, mockStore, idt}
		mockStore.On("Set", "f202ad3f", testmineURL).Return(nil)
		handler.CreateLink(c)
		mockStore.AssertExpectations(t)
		// Assertions
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, link, rec.Body.String())
	})

	t.Run("error case", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURLJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/createlink")
		mockStore := &mocks.Storer{}
		handler := &Handler{mockStore, mockStore, idt}
		mockStore.On("Set", "f202ad3f", testmineURL).Return(err)
		handler.CreateLink(c)
		mockStore.AssertExpectations(t)
		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestRedirect(t *testing.T) {

	t.Run("succes case with cache", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:shortUrl")
		c.SetParamNames("shortUrl")
		c.SetParamValues("f202ad3f")
		mockStore := &mocks.Storer{}
		mockCache := &mocks.Storer{}
		mockCache.On("Exists", id).Return(true, nil)
		mockCache.On("Get", id).Return(url, nil)
		handler := &Handler{mockStore, mockCache, idt}
		handler.Redirect(c)
		mockCache.AssertExpectations(t)
		// Assertions
		assert.Equal(t, http.StatusFound, rec.Code)
	})

	t.Run("succes case with datastore", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:shortUrl")
		c.SetParamNames("shortUrl")
		c.SetParamValues("f202ad3f")
		mockStore := &mocks.Storer{}
		mockCache := &mocks.Storer{}
		mockCache.On("Exists", id).Return(false, err)
		mockStore.On("Get", id).Return(url, nil)
		handler := &Handler{mockStore, mockCache, idt}
		handler.Redirect(c)
		mockStore.AssertExpectations(t)
		mockCache.AssertExpectations(t)
		// Assertions
		assert.Equal(t, http.StatusFound, rec.Code)
	})

	t.Run("error case", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:shortUrl")
		c.SetParamNames("shortUrl")
		c.SetParamValues("f202ad3f")
		mockStore := &mocks.Storer{}
		mockCache := &mocks.Storer{}
		mockCache.On("Exists", id).Return(false, err)
		mockStore.On("Get", id).Return("", err500)
		handler := &Handler{mockStore, mockCache, idt}
		handler.Redirect(c)
		mockStore.AssertExpectations(t)
		mockCache.AssertExpectations(t)
		// Assertions
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Setup
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/:shortUrl")
		c.SetParamNames("shortUrl")
		c.SetParamValues("f202a")
		mockStore := &mocks.Storer{}
		mockCache := &mocks.Storer{}
		mockCache.On("Exists", c.Param("shortUrl")).Return(false, err)
		mockStore.On("Get", c.Param("shortUrl")).Return("", err)
		handler := &Handler{mockStore, mockCache, idt}
		handler.Redirect(c)
		mockStore.AssertExpectations(t)
		mockCache.AssertExpectations(t)
		// Assertions
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
