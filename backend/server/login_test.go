package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kalogsc/ego/models"
)

func TestSignIn(t *testing.T) {
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
			fmt.Println(err)
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
	users := &[]struct {
		inputJSON    string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			inputJSON:    `{"email": "notspiderman@gmail.com", "password": "web123"}`,
			statusCode:   http.StatusOK,
			errorMessage: "",
		},
		{
			inputJSON:    `{"email": "incorrect@gmail.com", "password": "notexist"}`,
			statusCode:   http.StatusNotFound,
			errorMessage: "record not found",
		},
		{
			inputJSON:    `{"email": "amusic@gmail.com", "password": "notcoxinha123"}`,
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "incorrect password",
		},
		{
			inputJSON:    `{"email": "agmail.com", "password": "coxinha123"}`,
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "invalid email",
		},
		{
			inputJSON:    `{"email": "", "password": "web123"}`,
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Email' is required",
		},
		{
			inputJSON:    `{"email": "notspiderman@gmail.com", "password": ""}`,
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Password' is required",
		},
	}

	for _, v := range *users {
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.Login)
		handler.ServeHTTP(rr, req)

		if rr.Code != v.statusCode {
			t.Errorf("error: status code expected to be equal %d but was %d", v.statusCode, rr.Code)
		} else if v.statusCode == 200 && rr.Body.String() == "" {
			t.Errorf("error: body is empty")
		}

		if v.statusCode != http.StatusOK && v.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
			if err != nil {
				fmt.Println(rr.Body.String())
				t.Errorf("Cannot convert to json: %v", err)
			}
			if responseMap["error"] != v.errorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.errorMessage, responseMap["error"])
			}
		}
	}

}
