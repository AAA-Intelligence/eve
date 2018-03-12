package main

import (
	"bytes"
	"html/template"
	"net/http"
)

// function executes the given template and writes the result to a buffer
// if the execution was successfull the result is written to the ResponseWriter
// this avoids writing detailed error descriptions to the client
func saveExecute(w http.ResponseWriter, tpl *template.Template, data interface{}) error {
	var buffer bytes.Buffer
	err := tpl.Execute(&buffer, data)
	if err != nil {
		http.Error(w, ErrInternalServerError, http.StatusInternalServerError)
		return err
	}
	_, err = w.Write([]byte(buffer.String()))
	return err
)

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
