package main

import "net/http"

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s", r.Method, r.URL)

		_, err := app.authenticateToken(r)
		if err != nil {
			app.invalidCredentials(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
