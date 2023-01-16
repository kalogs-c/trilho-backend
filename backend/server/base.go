package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kalogsc/ego/models"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) InstanciateDB(DbName string) error {
	var err error
	DbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), DbName)
	server.DB, err = gorm.Open("mysql", DbUrl)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Initialize(DbName string) {
	err := server.InstanciateDB(DbName)
	if err != nil {
		fmt.Println("Cannot connect to mysql database.")
		log.Fatal("Failed to connect to db: ", err)
	}
	fmt.Println("Connected to mysql database.")

	server.DB.Debug().AutoMigrate(&models.User{})
	server.DB.Debug().AutoMigrate(&models.Transaction{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
