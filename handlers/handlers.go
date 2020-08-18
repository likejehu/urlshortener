package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/likejehu/urlshortener/db"
	"github.com/likejehu/urlshortener/models"
)

//Storer is  interface for  basic Key/Value (real and mock) datastorage for links
type Storer interface {
	Get(key string) (string, error)
	Set(key string, value *models.MineURL) error
	Exists(key string) (bool, error)
}

//Ider is  interface for idgenerator
type Ider interface {
	NewID() string
}

// Handler is struct for handlers
type Handler struct {
	Store Storer
	Cache Storer
	ID    Ider
}

// Handlers
//
var homeurl = "http://localhost:8080/"

// CreateLink is  function for creating short link from long URL
func (h *Handler) CreateLink(c echo.Context) (err error) {
	mineurl := new(models.MineURL)
	if err := c.Bind(mineurl); err != nil {
		return err
	}
	id := h.ID.NewID()
	mineurl.ShortURL = homeurl + id
	h.Cache.Set(id, mineurl)
	if err := h.Store.Set(id, mineurl); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, mineurl.ShortURL)
}

// Redirect is function for redirecting from short link to original URL
func (h *Handler) Redirect(c echo.Context) (err error) {
	key := c.Param("shortUrl")
	ok, _ := h.Cache.Exists(key)
	if ok {
		url, _ := h.Cache.Get(key)
		return c.Redirect(http.StatusFound, url)
	}
	url, err := h.Store.Get(key)
	if err != nil {
		if err == db.Error404 {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.Redirect(http.StatusFound, url)
}
