package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kalogsc/ego/models"
)

func TestCreateUser(t *testing.T) {
	users := &[]struct {
		user         *models.User
		statusCode   int
		errorMessage string
	}{
		{
			user: &models.User{
				Name:     "Carlos",
				LastName: "Henrique",
				Email:    "carlos@gmail.com",
				Password: "carlinhos!",
			},
			statusCode:   http.StatusCreated,
			errorMessage: "",
		},
		{
			user: &models.User{
				Name:     "Other",
				LastName: "Carlos",
				Email:    "carlos@gmail.com",
				Password: "notcarlinhos!",
			},
			statusCode:   http.StatusConflict,
			errorMessage: "email already taken",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "Santana",
				Email:    "mariagmail.com",
				Password: "maria!",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "invalid email",
		},
		{
			user: &models.User{
				Name:     "",
				LastName: "Santana",
				Email:    "maria@gmail.com",
				Password: "maria!",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Name' is required",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "",
				Email:    "maria@gmail.com",
				Password: "maria!",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'LastName' is required",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "Santana",
				Email:    "",
				Password: "maria!",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "field 'Email' is required",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "Santana",
				Email:    "maria@gmail.com",
				Password: "",
			},
			statusCode:   http.StatusUnprocessableEntity,
			errorMessage: "Required Password",
		},
	}

	for _, v := range *users {
		userJsonFormat, err := json.Marshal(v.user)
		if err != nil {
			t.Errorf("error marshalling json: %v", err)
		}

		req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(string(userJsonFormat)))
		if err != nil {
			t.Errorf("error creating request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.Bytes()), &responseMap)
		if err != nil {
			t.Errorf("failed to convert to json: %v", err)
		}
		if rr.Code != v.statusCode {
			t.Errorf("error: expected status code be %v but was %v", v.statusCode, rr.Code)
		}
		if v.statusCode == http.StatusCreated {
			if responseMap["name"] != v.user.Name {
				t.Errorf("error: expected field 'name' be %v but was %v", responseMap["name"], v.user.Name)
			}
			if responseMap["last_name"] != v.user.LastName {
				t.Errorf("error: expected field 'last_name' be %v but was %v", responseMap["last_name"], v.user.LastName)
			}
			if responseMap["email"] != v.user.Email {
				t.Errorf("error: expected field 'email' be %v but was %v", responseMap["email"], v.user.Email)
			}
			if responseMap["password"] == v.user.Password {
				t.Error("error: expected field 'password' be different")
			}
		}
		if v.statusCode == http.StatusUnprocessableEntity || v.statusCode == http.StatusInternalServerError && v.errorMessage != "" {
			if responseMap["error"] != v.errorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.errorMessage, responseMap["error"])
			}
		}
	}
}
