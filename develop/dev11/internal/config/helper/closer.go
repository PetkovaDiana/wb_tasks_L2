package helper

import (
	"io"
	"log"
)

// Closer - close connection with err handling
func Closer(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
	}
}
