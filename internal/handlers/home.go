package handlers

import (
	"net/http"

	"github.com/justestif/go-web-template/components"
)

func Home(w http.ResponseWriter, r *http.Request) {
	components.Home().Render(r.Context(), w)
}
