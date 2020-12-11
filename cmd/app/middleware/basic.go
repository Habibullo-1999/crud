package middleware

import (
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

)
//Basic 
func Basic(chAuth func(string, string) bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// extract username and password
			login, pass, err := gLogPass(r)
			if err != nil {
				log.Println(err)
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			
			if !chAuth(login, pass) {
				http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			handler.ServeHTTP(rw, r)
		})
	}
}
//func that extracts data from the request and returns the login and password
func gLogPass(r *http.Request) (string, string, error) {

	auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", "", errors.New("invalid auth method")
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	log.Println(pair)
	if len(pair) != 2 {

		return "", "", errors.New("invalid auth data")
	}

	return pair[0], pair[1], nil
}