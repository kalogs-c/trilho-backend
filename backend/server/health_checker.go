package server

import (
	"net/http"

	"github.com/kalogsc/ego/responses"
)

func (server *Server) HealthChecker(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Ego Backend HC")
}
