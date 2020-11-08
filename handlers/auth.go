package handlers

import (
	"context"
	"net/http"
)

func (uh *UserHandler) MustAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session.id")
		var uid uint
		switch id := session.Values["id"].(type) {
		case int:
			uid = uint(id)
		case uint:
			uid = id
		default:
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user := uh.UserService.FindByID(uid)
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
