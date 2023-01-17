package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kalogsc/ego/models"
	"github.com/kalogsc/ego/utils"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) InstanciateDB(DbName string, mode utils.DbModeEnum) error {
	var err error
	var DbUrl string
	if mode == utils.DB_MODE_PROD {
		DbUrl = fmt.Sprintf("%s@unix(%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_HOST"), DbName)
	} else if mode == utils.DB_MODE_TEST {
		DbUrl = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("DB_PORT"), DbName)
	}
	server.DB, err = gorm.Open("mysql", DbUrl)
	if err != nil {
		return err
	}
	return nil
}

func (server *Server) Initialize(DbName string, mode utils.DbModeEnum) {
	err := server.InstanciateDB(DbName, mode)
	if err != nil {
		fmt.Println("Cannot connect to mysql database.")
		log.Fatal("Failed to connect to db: ", err)
	}
	fmt.Println("Connected to mysql database.")

	server.DB.Debug().AutoMigrate(&models.User{}, &models.Transaction{})

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listening to port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}
