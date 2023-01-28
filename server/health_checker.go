package server

import (
	"net/http"

	"github.com/kalogsc/trilho/responses"
)

func (server *Server) HealthChecker(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "trilho Backend HC")
}
