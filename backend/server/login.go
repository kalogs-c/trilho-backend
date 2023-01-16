package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kalogsc/ego/auth"
	"github.com/kalogsc/ego/models"
	"github.com/kalogsc/ego/responses"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err :=  io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user.Validate("login")

	token, err := server.SignIn(&user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(user *models.User) (string, error) {
	u := models.User{}

	err := server.DB.Debug().Model(models.User{}).Where("email = ?", user.Email).Take(&u).Error
	if err != nil {
		return "", err
	}

	err = models.VerifyPassword(u.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(u.ID)
}
