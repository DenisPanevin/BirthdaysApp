package main

import (
	"context"
	"net/http"
)

func (a *Application) authUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := a.sessionStore.Get(r, SessionName)
		if err != nil {
			a.customErr(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["UserId"]
		if !ok {
			a.customErr(w, r, http.StatusUnauthorized, notAuth)
			return
		}

		u, err := a.storage.FindById(id.(int))
		if !ok {
			a.customErr(w, r, http.StatusInternalServerError, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyUser, u)))
	})
}
