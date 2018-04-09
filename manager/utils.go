package manager

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/AAA-Intelligence/eve/db"
)

// function executes the given template and writes the result to a buffer
// if the execution was successfull the result is written to the ResponseWriter
// this avoids writing detailed error descriptions to the client
func saveExecute(w http.ResponseWriter, tpl *template.Template, data interface{}) error {
	var buffer bytes.Buffer
	err := tpl.Execute(&buffer, data)
	if err != nil {
		http.Error(w, db.ErrInternalServerError.Error(), http.StatusInternalServerError)
		return err
	}
	_, err = w.Write([]byte(buffer.String()))
	return err
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.RawURLEncoding.EncodeToString(b)[:s], err
}

// Formats the time struct to hh:mm (e.g. 19:45)
func formatTime(time *time.Time) string {
	return fmt.Sprintf("%d:%02d", time.Hour(), time.Minute())
}

// Returns the number of years that have passed since the given date
func yearsSince(date *time.Time) int {
	return int(time.Since(*date).Hours() / float64(24*356))
}
