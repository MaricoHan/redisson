package mw

import (
	"net/http"
)

// Middleware define a middleware
type Middleware func(http.Handler) http.Handler
