package web

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func LoggingHandler(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}
