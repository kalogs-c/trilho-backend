package server

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/kalogsc/trilho/auth"
	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/responses"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(&user)
	if err != nil {
		switch err.Error() {
		case bcrypt.ErrMismatchedHashAndPassword.Error():
			responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("incorrect password"))
		case "record not found":
			responses.ERROR(w, http.StatusNotFound, err)
		default:
			responses.ERROR(w, http.StatusInternalServerError, err)
		}
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
	if err != nil || err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(u.ID)
}
