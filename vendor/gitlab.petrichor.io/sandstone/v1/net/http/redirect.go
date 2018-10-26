package http

import (
	"fmt"
	"net/http"
)

var RedirectToSecure = http.HandlerFunc(
	func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, fmt.Sprintf("https://%s%s", req.Host, req.URL.EscapedPath()), 301)
	},
)
