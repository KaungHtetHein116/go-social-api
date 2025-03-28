package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s -> %s: %v", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s -> %s: %v", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s -> %s: %v", r.Method, r.URL.Path, err)

	writeJSONError(w, http.StatusNotFound, err.Error())
}
