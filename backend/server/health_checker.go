package server

import (
	"net/http"

	"github.com/kalogs-c/piadocas/responses"
)

func (server *Server) HealthChecker(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Piadocas HC")
}
