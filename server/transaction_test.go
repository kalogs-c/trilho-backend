package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/seed"
)

func TestCreateTransaction(t *testing.T) {
	user := &models.User{
		Name:     "Zakk",
		LastName: "Wylde",
		Username:    "zakk",
		Password: "guitar123",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{user})
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
		return
	}

	user.Password = "guitar123" // Reset password
	token, err := serverInstance.SignIn(user)
	if err != nil {
		t.Errorf("error: failed while sign in %e", err)
		return
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	exampleTransaction := &models.Transaction{
		Name:    "CD Player",
		Amount:  250,
		OwnerId: user.ID,
	}
	transactions := []struct {
		transaction          *models.Transaction
		expectedStatusCode   int
		tokenGiven           string
		expectedErrorMessage string
	}{
		{
			transaction:          exampleTransaction,
			expectedStatusCode:   http.StatusCreated,
			tokenGiven:           tokenString,
			expectedErrorMessage: "",
		},
		{
			transaction: &models.Transaction{
				Name:    "CD Player",
				Amount:  250,
				OwnerId: 100,
			},
			expectedStatusCode:   http.StatusUnauthorized,
			tokenGiven:           tokenString,
			expectedErrorMessage: "Unauthorized",
		},
		{
			transaction:          exampleTransaction,
			expectedStatusCode:   http.StatusUnauthorized,
			tokenGiven:           "", // empty
			expectedErrorMessage: "Unauthorized",
		},
		{
			transaction:          exampleTransaction,
			expectedStatusCode:   http.StatusUnauthorized,
			tokenGiven:           "incorrect token",
			expectedErrorMessage: "Unauthorized",
		},
		{
			transaction: &models.Transaction{
				Name:    "",
				Amount:  250,
				OwnerId: user.ID,
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			tokenGiven:           tokenString,
			expectedErrorMessage: "field 'Name' is required",
		},
	}
	for _, v := range transactions {
		transactionJson, err := json.Marshal(v.transaction)
		if err != nil {
			t.Errorf("error marshalling json: %v", err)
			return
		}

		req, err := http.NewRequest("POST", "/transaction", bytes.NewBuffer(transactionJson))
		if err != nil {
			t.Errorf("error creating request: %v", err)
			return
		}
		req = mux.SetURLVars(req, map[string]string{
			"user_id": strconv.Itoa(int(v.transaction.OwnerId)),
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.CreateTransaction)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
			return
		}
		if rr.Code != v.expectedStatusCode {
			t.Errorf("error: expected status code be %v but was %v", v.expectedStatusCode, rr.Code)
			return
		}
		if v.expectedStatusCode == http.StatusCreated {
			if responseMap["name"] != v.transaction.Name {
				t.Errorf("error: expected field 'name' be %v but was %v", responseMap["name"], v.transaction.Name)
				return
			}
			if responseMap["amount"] != v.transaction.Amount {
				t.Errorf("error: expected field 'amount' be %v but was %v", responseMap["amount"], v.transaction.Amount)
				return
			}
			if responseMap["owner_id"] != float64(v.transaction.OwnerId) {
				t.Errorf("error: expected field 'owner_id' be %v but was %v", responseMap["owner_id"], v.transaction.OwnerId)
				return
			}
		} else if v.expectedErrorMessage != "" {
			if responseMap["error"] != v.expectedErrorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.expectedErrorMessage, responseMap["error"])
				return
			}
		}
	}
}

func TestListUsersTransactions(t *testing.T) {
	user := &models.User{
		Name:     "Zakk",
		LastName: "Wylde",
		Username:    "zakk",
		Password: "guitar123",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{user})
	if err != nil {
		t.Errorf("error seeding custom users: %v", err)
		return
	}

	user.Password = "guitar123" // Reset password
	token, err := serverInstance.SignIn(user)
	if err != nil {
		t.Errorf("error: failed while sign in %e", err)
		return
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	req = mux.SetURLVars(req, map[string]string{
		"user_id": strconv.Itoa(int(user.ID)),
	})
	req.Header.Add("Authorization", tokenString)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serverInstance.ListUserTransactions)
	handler.ServeHTTP(rr, req)

	var transactions []models.Transaction
	err = json.Unmarshal(rr.Body.Bytes(), &transactions)
	if err != nil {
		t.Errorf("Cannot convert to json: %v", err)
		return
	}

	if rr.Code != http.StatusOK {
		t.Errorf("error: expected status code be %v but was %v", http.StatusOK, rr.Code)
		return
	}
	if len(transactions) < 1 {
		t.Errorf("error: expected at least one transaction. %v", transactions)
		return
	}
}

func TestUpdateTransaction(t *testing.T) {
	user := &models.User{
		Name:     "Zakk",
		LastName: "Wylde",
		Username:    "zakk",
		Password: "guitar123",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{user})
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
		return
	}

	user.Password = "guitar123" // Reset password
	token, err := serverInstance.SignIn(user)
	if err != nil {
		t.Errorf("error: failed while sign in %e", err)
		return
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	oldTransaction := models.Transaction{
		OwnerId: user.ID,
	}
	transactionsList, err := oldTransaction.CollectUserTransactions(serverInstance.DB)
	if err != nil {
		t.Errorf("error: failed while collecting transactions")
		return
	}
	oldTransaction = (*transactionsList)[0]

	updateTransaction := &models.Transaction{
		ID:      oldTransaction.ID,
		Name:    "CD Player",
		Amount:  250,
		OwnerId: user.ID,
	}
	transactions := []struct {
		transaction          *models.Transaction
		expectedStatusCode   int
		tokenGiven           string
		expectedErrorMessage string
	}{
		{
			transaction:          updateTransaction,
			expectedStatusCode:   http.StatusOK,
			tokenGiven:           tokenString,
			expectedErrorMessage: "",
		},
		{
			transaction:          updateTransaction,
			tokenGiven:           "",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
		{
			transaction:          updateTransaction,
			tokenGiven:           "this is an incorrect token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
		{
			transaction: &models.Transaction{
				ID:      oldTransaction.ID,
				Name:    "",
				Amount:  250,
				OwnerId: user.ID,
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			tokenGiven:           tokenString,
			expectedErrorMessage: "field 'Name' is required",
		},
		{
			transaction: &models.Transaction{
				ID:      oldTransaction.ID,
				Name:    "",
				Amount:  250,
				OwnerId: 100,
			},
			expectedStatusCode:   http.StatusUnauthorized,
			tokenGiven:           tokenString,
			expectedErrorMessage: "Unauthorized",
		},
	}
	for _, v := range transactions {
		transactionJson, err := json.Marshal(v.transaction)
		if err != nil {
			t.Errorf("error marshalling json: %v", err)
			return
		}

		req, err := http.NewRequest("POST", "/transaction", bytes.NewBuffer(transactionJson))
		if err != nil {
			t.Errorf("error creating request: %v", err)
			return
		}
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(int(v.transaction.ID)),
		})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.UpdateTransaction)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
			return
		}
		if rr.Code != v.expectedStatusCode {
			t.Errorf("error: expected status code be %v but was %v", v.expectedStatusCode, rr.Code)
			return
		}
		if v.expectedStatusCode == http.StatusOK {
			if responseMap["name"] == oldTransaction.Name {
				t.Errorf("error: expected field 'name' be different %v but was %v", responseMap["name"], v.transaction.Name)
				return
			}
			if responseMap["amount"] == oldTransaction.Amount {
				t.Errorf("error: expected field 'amount' be different %v but was %v", responseMap["amount"], v.transaction.Amount)
				return
			}
			if responseMap["owner_id"] != float64(v.transaction.OwnerId) {
				t.Errorf("error: expected field 'owner_id' be %v but was %v", responseMap["owner_id"], v.transaction.OwnerId)
				return
			}
		} else if v.expectedErrorMessage != "" {
			if responseMap["error"] != v.expectedErrorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.expectedErrorMessage, responseMap["error"])
				return
			}
		}
	}
}

func TestDeleteTransaction(t *testing.T) {
	user := &models.User{
		Name:     "Zakk",
		LastName: "Wylde",
		Username:    "zakk",
		Password: "guitar123",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{user})
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
		return
	}

	user.Password = "guitar123" // Reset password
	token, err := serverInstance.SignIn(user)
	if err != nil {
		t.Errorf("error: failed while sign in %e", err)
		return
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	transaction := models.Transaction{
		OwnerId: user.ID,
	}
	transactionsList, err := transaction.CollectUserTransactions(serverInstance.DB)
	if err != nil {
		t.Errorf("error: failed while collecting transactions")
		return
	}
	transaction = (*transactionsList)[0]

	transactions := []struct {
		transactionID        string
		tokenGiven           string
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		{
			transactionID:        strconv.Itoa(int(transaction.ID)),
			tokenGiven:           tokenString,
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
		{
			transactionID:        strconv.Itoa(int(transaction.ID)),
			tokenGiven:           "",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
		{
			transactionID:        strconv.Itoa(int(transaction.ID)),
			tokenGiven:           "This is an incorrect token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
	}
	for _, v := range transactions {
		req, _ := http.NewRequest("DELETE", "/transaction", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.transactionID})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.DeleteTransaction)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		if rr.Code != v.expectedStatusCode {
			t.Errorf("error: expected status code be %v but was %v", v.expectedStatusCode, rr.Code)
			return
		}

		if v.expectedStatusCode != 204 && v.expectedErrorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			if responseMap["error"] != v.expectedErrorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.expectedErrorMessage, responseMap["error"])
				return
			}
		}
	}
}
