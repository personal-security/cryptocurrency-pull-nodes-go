package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	easysdk "github.com/personal-security/easy-sdk-go"
)

var AuthGetUserToken = func(r *http.Request) string {
	tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

	if len(tokenHeader) == 0 {
		params := r.URL.Query()
		tokenHeader = params.Get("token")
	}

	return tokenHeader
}

var KeyAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/v1/api/status"} //List of endpoints that doesn't require auth
		notAuthPrefix := []string{
			"/captcha/",
			"/v1/api/prices",
		}
		notAuthRegExp := []string{
			`^/v1/api/node/[a-z0-9]+/info$`,
			`^/v1/api/node/[a-z0-9]+/wallet/[a-zA-Z0-9]+$`,
			`^/v1/api/node/[a-z0-9]+/tx/[a-z0-9]+$`,
			`^/v1/api/node/[a-z0-9]+/helper/transaction$`,
			`^/v1/api/node/[a-z0-9]+/broadcast/raw$`,
		}
		requestPath := r.URL.Path //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		for _, value := range notAuthPrefix {

			if strings.HasPrefix(requestPath, value) {
				next.ServeHTTP(w, r)
				return
			}
		}

		for _, value := range notAuthRegExp {
			reg, err := regexp.Compile(value)
			if err != nil {
				//fmt.Print(err.Error())
			} else {
				if reg.MatchString(r.URL.EscapedPath()) {
					next.ServeHTTP(w, r)
					return
				}
			}
		}

		//response := make(map[string]interface{})
		apiKeyReceived := AuthGetUserToken(r)

		if apiKeyReceived == "" {
			response := easysdk.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			easysdk.Respond(w, response)
			return
		}

		type ApiKey struct {
			Key string `json:"api_key"`
		}

		var apiKey ApiKey
		file, _ := ioutil.ReadFile("config/keys.json")
		json.Unmarshal([]byte(file), &apiKey)

		if apiKey.Key != apiKeyReceived {
			response := easysdk.Message(false, "Bad auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			easysdk.Respond(w, response)
			return
		}

		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}
