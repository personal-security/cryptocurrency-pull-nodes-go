package app

import (
	"net/http"

	easysdk "github.com/personal-security/easy-sdk-go"
)

var NotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	easysdk.Respond(w, easysdk.Message(false, "This resources was not found on our server"))
}
