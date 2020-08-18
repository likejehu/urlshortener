package idgen

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/likejehu/urlshortener/db"
)

// ShortID is an implementation of the ider Interface
type ShortID struct {
	SI string
}

// NewID  retuns new short id
func (s *ShortID) NewID() string {
	s.SI = NewShortID()
	return s.SI
}

// NewShortID  return a string from 4 random bytes
func NewShortID() string {
	rb := make([]byte, 4)
	rand.Read(rb)
	sid := hex.EncodeToString(rb)
	if t, _ := db.LinksStore.Exists(sid); t == true {
		return NewShortID()
	}
	return sid
}

// IDG is instance of ShortID stucture
var IDG = &ShortID{}
