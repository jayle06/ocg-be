package middlewares

import "net/http"

func CommonMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Content-Length, "+
			"Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, "+
			"Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, "+
			"Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
