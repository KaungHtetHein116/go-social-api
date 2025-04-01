package main

import (
	"net/http"
	"social/internal/store"

	"github.com/go-playground/validator/v10"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "asc",
	}

	fq, err := fq.Parse(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if err := v.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	feed, err := app.store.Posts.GetUserFeed(r.Context(), 1, fq)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
