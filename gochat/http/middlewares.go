package http

import (
	"context"
	"net/http"
)

// FlashMiddleware retrieves flash messages from request.
func FlashMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, _ := sessionStore.Get(r, "flash")
	messages := []FlashMessage{}
	for _, v := range session.Flashes("flash") {
		message := DeserializeFlashMessage(v.(string))
		messages = append(messages, message)
	}
	session.Save(r, w)
	r = r.WithContext(context.WithValue(r.Context(), "messages", messages))
	next(w, r)

}

// AuthMiddleware checks to see if a user is authenticated or not,
// blocking the user from proceeding to protected paths.
func (s *Server) AuthMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	sessionValues, ok := r.Context().Value("session").(map[interface{}]interface{})
	if !ok {
		next(w, r)
		return
	}
	is_auth, _ := sessionValues["is_auth"].(bool)
	if !is_auth {
		addFlash(w, r, "info", "You need to login to continue")
		http.Redirect(w, r, "/login?next="+r.URL.Path, http.StatusTemporaryRedirect)
		return
	}
	id, _ := sessionValues["user_id"].(int)
	user, err := s.services.user.FindUserByID(context.Background(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	r = r.WithContext(NewContextWithUser(r.Context(), user))
	next(w, r)
}

// TODO: Change session backend to a db.
// SessionMiddleware retrieves a user session from the session backend.
func SessionMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	session, _ := sessionStore.Get(r, "session.id")
	r = r.WithContext(context.WithValue(r.Context(), "session", session.Values))

	next(w, r)
}
