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
}
