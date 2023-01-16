package server

import "github.com/kalogs-c/piadocas/middlewares"

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", middlewares.FormatToJSON(server.HealthChecker)).Methods("GET")

	server.Router.HandleFunc("/user", middlewares.FormatToJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/user/{id}", middlewares.FormatToJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/user/{id}", middlewares.FormatToJSON(server.UpdateUser)).Methods("PATCH")

	server.Router.HandleFunc("/transaction/{user_id}", middlewares.FormatToJSON(server.GetUserTransactions)).Methods("GET")
	server.Router.HandleFunc("/transaction/{id}", middlewares.FormatToJSON(server.GetTransactionById)).Methods("GET")
	server.Router.HandleFunc("/transaction/{user_id}", middlewares.FormatToJSON(server.AddTransaction)).Methods("POST")
	server.Router.HandleFunc("/transaction/{id}", middlewares.FormatToJSON(server.UpdateTransaction)).Methods("PATCH")
}
