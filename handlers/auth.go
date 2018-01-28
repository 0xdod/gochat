package handlers

import (
	"context"
	"net/http"
)

func (uh *UserHandler) MustAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session.id")
		id, ok := session.Values["user_id"].(uint)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if rID, ok := session.Values["room_id"].(uint); ok {
			room := uh.RoomService.FindByID(rID)
			ctx := context.WithValue(r.Context(), "room", room)
			r = r.WithContext(ctx)
		}
		user := uh.UserService.FindByID(id)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func Get(r *http.Request, key string) interface{} {
	return r.Context().Value(key)
}

func Set(r *http.Request, key string, value interface{}) {
	ctx := context.WithValue(r.Context(), key, value)
	r = r.WithContext(ctx)
}
