package main

import (
	"fmt"
	"net/http"
)

func (app *application) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s", r.Method, r.URL)

		_, err := app.authenticateToken(r)
		if err != nil {
			fmt.Println("-------------- errrrrroooorrr auth")
			app.invalidCredentials(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}
