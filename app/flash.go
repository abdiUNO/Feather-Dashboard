package app

import (
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("SESSION_KEY"))

// var Flash = func(next http.Handler) http.Handler {
//
// }
