package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/seed"
)

func TestSignIn(t *testing.T) {
	customUser := &[]*models.User{
		{
			Name:     "Alessia",
			LastName: "Cara",
			Username: "amusic",
			Password: "coxinha123",
		},
	}
	err := seed.LoadCustomUsers(serverInstance.DB, customUser)
	if err != nil {
		t.Fatal(err)
	}

	users := &[]struct {
		user                 models.User
		statusCode           int
		expectedErrorMessage string
	}{
		{
			user: models.User{
				Username: "amusic",
				Password: "coxinha123",
			},
			expectedErrorMessage: "",
		},
		{
			user: models.User{
				Username: "amusic",
				Password: "wrongpassword",
			},
			expectedErrorMessage: "crypto/bcrypt: hashedPassword is not the hash of the given password",
		},
		{
			user: models.User{
				Username: "notexist",
				Password: "kalogs",
			},
			expectedErrorMessage: "record not found",
		},
	}

	for _, v := range *users {
		token, err := serverInstance.SignIn(&v.user)
		if err == nil && token != "" {
			return
		}

		if err.Error() != v.expectedErrorMessage {
			t.Errorf("invalid error, expected to be %v but was %v", v.expectedErrorMessage, err.Error())
		} else if err != nil {
			t.Errorf("error while sign in: %v", err.Error())
		}

		if token == "" {
			t.Error("error: token empty")
		}
	}
}

func TestLogin(t *testing.T) {
	customUsers := &[]*models.User{
		{
			Name:     "Peter",
			LastName: "Parker",
			Username: "notspiderman",
			Password: "web123",
		},
	}
	err := seed.LoadCustomUsers(serverInstance.DB, customUsers)
	if err != nil {
		t.Fatal(err)
	}

	users := &[]struct {
		user         models.User
		statusCode   int
		errorMessage string
	}{
		{
			user: models.User{
				Username: "notspiderman",
				Password: "web123",
			},
			statusCode:   http.StatusOK,
			errorMessage: "",
		},
		{
			user: models.User{
				Username: "incorrect",
				Password: "notexist",
			},
			statusCode:   http.StatusNotFound,
			errorMessage: "record not found",
		},
		{
			user: models.User{
				Username: "notspiderman",
				Password: "itsnotcorrect",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "incorrect password",
		},
		{
			user: models.User{
				Username: "",
				Password: "web123",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Username' is required",
		},
		{
			user: models.User{
				Username: "notspiderman",
				Password: "",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Password' is required",
		},
	}

	for _, v := range *users {
		jsonUser, err := json.Marshal(v.user)
		if err != nil {
			t.Errorf("error mashalling json: %v", err)
		}

		req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonUser))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != v.statusCode {
			t.Errorf("error: status code expected to be equal %d but was %d\n Response Body: %v", v.statusCode, rr.Code, rr.Body.String())
		} else if v.statusCode == 200 && rr.Body.String() == "" {
			t.Errorf("error: body is empty")
		}

		if v.statusCode != http.StatusOK && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			if responseMap["error"] != v.errorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.errorMessage, responseMap["error"])
			}
		}
	}
}
