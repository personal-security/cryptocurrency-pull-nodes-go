package controllers

import (
	"net/http"

	easysdk "github.com/personal-security/easy-sdk-go"
)

var StatusGetNow = func(w http.ResponseWriter, r *http.Request) {
	resp := easysdk.Message(true, "Success")
	easysdk.Respond(w, resp)
}
