package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/utils"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) InstanciateDB(DbName string, mode utils.DbModeEnum) error {
	var err error
	server.DB, err = gorm.Open("postgres", os.Getenv("DB_CONN"))
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Initialize(DbName string, mode utils.DbModeEnum) {
	err := server.InstanciateDB(DbName, mode)
	if err != nil {
		fmt.Println("Cannot connect to postgres database.")
		log.Fatal("Failed to connect to db: ", err)
	}
	fmt.Println("Connected to postgres database.")

	err = server.DB.Debug().AutoMigrate(&models.User{}, &models.Transaction{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = server.DB.Debug().Model(&models.Transaction{}).AddForeignKey("owner_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
