package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/kalogsc/ego/auth"
	"github.com/kalogsc/ego/models"
	"github.com/kalogsc/ego/responses"
)

func (server *Server) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	transaction := models.Transaction{}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	vars := mux.Vars(r)
	ownerId, err := strconv.ParseInt(vars["user_id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	transaction.OwnerId = uint32(ownerId)

	transactionCreated, err := transaction.Save(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, transactionCreated.ID))
	responses.JSON(w, http.StatusCreated, transactionCreated)
}

func (server *Server) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionId, err := strconv.ParseInt(vars["user_id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	transaction := models.Transaction{OwnerId: uint32(transactionId)}
	transactionsData, err := transaction.CollectUserTransactions(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	responses.JSON(w, http.StatusOK, transactionsData)
}

func (server *Server) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionId, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	transaction := models.Transaction{}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	transaction.ID = uint32(transactionId)
	tokenID, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(transactionId) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}

	err = transaction.UpdateTransaction(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, transaction)
}

func (server *Server) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionId, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}

	transaction := models.Transaction{ID: uint32(transactionId)}
	err = transaction.Delete(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Entity", fmt.Sprintf("%d", transactionId))
	responses.JSON(w, http.StatusOK, "Deleted sucessfully")
}
