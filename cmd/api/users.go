package main

import (
	"context"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const userContextKey contextKey = "post"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.getUserFromContext(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	userToFollow := app.getUserFromContext(r)

	app.store.Followers.Follow(r.Context(), userToFollow.ID, 1)

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	userToUnfollow := app.getUserFromContext(r)

	app.store.Followers.Unfollow(r.Context(), userToUnfollow.ID, 1)

	if err := app.jsonResponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) usersContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "userID")
		userID, err := strconv.Atoi(idParam)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		user, err := app.store.Users.GetByUserID(r.Context(), userID)

		if err != nil {
			switch err {
			case store.ErrRecordNotFound:
				app.notFoundResponse(w, r, err)
				return

			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (app *application) getUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		return nil
	}
	return user
}
