package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kalogsc/trilho/auth"
	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/responses"
)

func (server *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Save(server.DB)
	if err != nil {
		switch err.Error() {
		case "username already taken":
			responses.ERROR(w, http.StatusConflict, err)
		default:
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
		}
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, user.ID))
	responses.JSON(w, http.StatusCreated, user)
}

func (server *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	user := models.User{ID: uint32(userId)}
	err = user.CollectUserData(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, user)
}

func (server *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, users)
}

func (server *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	jwtUserId, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if jwtUserId != uint32(userId) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	body, err := io.ReadAll(r.Body)
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
	user.ID = uint32(userId)

	err = user.UpdateUser(server.DB)
	if err != nil {
		switch err.Error() {
		case "username already taken":
			responses.ERROR(w, http.StatusConflict, err)
		default:
			responses.ERROR(w, http.StatusUnprocessableEntity, err)
		}
		return
	}
	responses.JSON(w, http.StatusOK, user)
}

func (server *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		if strings.Contains(err.Error(), "Cannot convert this id to an integer") {
			responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("cannot convert this id to an integer"))
		} else {
			responses.ERROR(w, http.StatusBadRequest, err)
		}
		return
	}

	jwtUserId, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if jwtUserId != uint32(userId) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	user := models.User{ID: uint32(userId)}
	err = user.Delete(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", userId))

	responseMap := make(map[string]string)
	responseMap["message"] = "Deleted sucessfully"

	responses.JSON(w, http.StatusOK, responseMap)
}
