package server

import "github.com/kalogsc/ego/middlewares"

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", middlewares.FormatToJSON(server.HealthChecker)).Methods("GET")

	server.Router.HandleFunc("/login", middlewares.FormatToJSON(server.Login)).Methods("POST")

	server.Router.HandleFunc("/user", middlewares.FormatToJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/user", middlewares.FormatToJSON(server.ListUsers)).Methods("GET")
	server.Router.HandleFunc("/user/{id}", middlewares.FormatToJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/user/{id}", middlewares.FormatToJSON(middlewares.Authentication(server.DeleteUser))).Methods("DELETE")
	server.Router.HandleFunc("/user/{id}", middlewares.FormatToJSON(middlewares.Authentication(server.UpdateUser))).Methods("PATCH")

	server.Router.HandleFunc("/transaction/{user_id}", middlewares.FormatToJSON(server.CreateTransaction)).Methods("POST")
	server.Router.HandleFunc("/transaction/{user_id}", middlewares.FormatToJSON(server.GetUserTransactions)).Methods("GET")
	server.Router.HandleFunc("/transaction/{id}", middlewares.FormatToJSON(middlewares.Authentication(server.DeleteTransaction))).Methods("DELETE")
	server.Router.HandleFunc("/transaction/{id}", middlewares.FormatToJSON(middlewares.Authentication(server.UpdateTransaction))).Methods("PATCH")
}
