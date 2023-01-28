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
			Email:    "amusic@gmail.com",
			Password: "coxinha123",
		},
	}
	err := seed.LoadCustomUsers(serverInstance.DB, customUser)
	if err != nil {
		t.Fatal(err)
	}

	users := &[]models.User{
		{
			Email:    "amusic@gmail.com",
			Password: "coxinha123",
		},
		{
			Email:    "amusic@gmail.com",
			Password: "wrongpassword",
		},
		{
			Email:    "notexist@gmail.com",
			Password: "kalogs",
		},
	}

	expectedErrors := &[]string{
		"",
		"crypto/bcrypt: hashedPassword is not the hash of the given password",
		"record not found",
	}

	for i := range *users {
		u := &(*users)[i]
		expectedErr := &(*expectedErrors)[i]
		token, err := serverInstance.SignIn(u)
		if err == nil && token != "" {
			return
		}

		if err.Error() != *expectedErr {
			t.Errorf("invalid error, expected to be %v but was %v", *expectedErr, err.Error())
		} else if err != nil {
			t.Errorf("error while sign in: %e", err)
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
			Email:    "notspiderman@gmail.com",
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
				Email:    "notspiderman@gmail.com",
				Password: "web123",
			},
			statusCode:   http.StatusOK,
			errorMessage: "",
		},
		{
			user: models.User{
				Email:    "incorrect@gmail.com",
				Password: "notexist",
			},
			statusCode:   http.StatusNotFound,
			errorMessage: "record not found",
		},
		{
			user: models.User{
				Email:    "notspiderman@gmail.com",
				Password: "itsnotcorrect",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "incorrect password",
		},
		{
			user: models.User{
				Email:    "agmail.com",
				Password: "coxinha123",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "invalid email",
		},
		{
			user: models.User{
				Email:    "",
				Password: "web123",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Email' is required",
		},
		{
			user: models.User{
				Email:    "notspiderman@gmail.com",
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
