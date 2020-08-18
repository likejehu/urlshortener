package models

import (
	"encoding/json"
)

// MineURL is structure to store URL
type MineURL struct {
	LongURL  string `json:"longUrl" form:"longUrl" query:"longUrl" `
	ShortURL string `json:"shortUrl" form:"shortUrl" query:"shortUrl" `
}

//MarshalBinary is for impelmentation
func (m MineURL) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}
