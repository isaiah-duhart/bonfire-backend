package handlers

import (
	"fmt"
	"net/http"

)

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, `
            <!DOCTYPE html>
            <html>
            <head>
                <title>Welcome</title>
            </head>
            <body>
                <h1>Hello from Bonfire!</h1>
            </body>
            </html>
        `)
}