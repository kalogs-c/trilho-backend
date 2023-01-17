package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
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
			errorMessage: "field 'Password' is required",
		},
	}

	for _, v := range *users {
		userJsonFormat, err := json.Marshal(v.user)
		if err != nil {
			t.Errorf("error marshalling json: %v", err)
			return
		}

		req, err := http.NewRequest("POST", "/user", bytes.NewBufferString(string(userJsonFormat)))
		if err != nil {
			t.Errorf("error creating request: %v", err)
			return
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.CreateUser)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("failed to convert to json: %v", err)
			return
		}
		if rr.Code != v.statusCode {
			t.Errorf("error: expected status code be %v but was %v", v.statusCode, rr.Code)
			return
		}
		if v.statusCode == http.StatusCreated {
			if responseMap["name"] != v.user.Name {
				t.Errorf("error: expected field 'name' be %v but was %v", responseMap["name"], v.user.Name)
				return
			}
			if responseMap["last_name"] != v.user.LastName {
				t.Errorf("error: expected field 'last_name' be %v but was %v", responseMap["last_name"], v.user.LastName)
				return
			}
			if responseMap["email"] != v.user.Email {
				t.Errorf("error: expected field 'email' be %v but was %v", responseMap["email"], v.user.Email)
				return
			}
			if responseMap["password"] == v.user.Password {
				t.Error("error: expected field 'password' be different")
				return
			}
		}
		if v.statusCode == http.StatusUnprocessableEntity || v.statusCode == http.StatusInternalServerError && v.errorMessage != "" {
			if responseMap["error"] != v.errorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.errorMessage, responseMap["error"])
				return
			}
		}
	}
}

func TestListUsers(t *testing.T) {
	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
		return
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serverInstance.ListUsers)
	handler.ServeHTTP(rr, req)

	var users *[]models.User
	err = json.Unmarshal(rr.Body.Bytes(), &users)
	if err != nil {
		t.Errorf("failed to convert to json: %v", err)
		return
	}

	if rr.Code != http.StatusOK {
		t.Errorf("error: expected status code be %v but was %v", http.StatusOK, rr.Code)
		return
	}
	if len(*users) < 1 {
		t.Errorf("no users was returned: %v", users)
		return
	}
}

func TestGetUser(t *testing.T) {
	user := &models.User{
		Name:     "Test",
		LastName: "Serious Test",
		Email:    "really@serious.com",
		Password: "trust_me157",
	}
	err := user.Save(serverInstance.DB)
	if err != nil {
		t.Errorf("failed to seed one user, %e", err)
		return
	}

	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
		return
	}
	req = mux.SetURLVars(req, map[string]string{
		"id": strconv.Itoa(int(user.ID)),
	})
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(serverInstance.GetUser)
	handler.ServeHTTP(rr, req)

	var userFromResponse *models.User
	err = json.Unmarshal(rr.Body.Bytes(), &userFromResponse)
	if err != nil {
		t.Errorf("failed to convert to json: %v", err)
		return
	}

	if rr.Code != http.StatusOK {
		t.Errorf("error: expected status code be %v but was %v", http.StatusOK, rr.Code)
	}

	if userFromResponse.Name != user.Name {
		t.Errorf("Expected name field be equal %v but was %v\n", user.Name, userFromResponse.Name)
	}
	if userFromResponse.LastName != user.LastName {
		t.Errorf("Expected lastname field be equal %v but was %v\n", user.LastName, userFromResponse.LastName)
	}
	if userFromResponse.Email != user.Email {
		t.Errorf("Expected email field be equal %v but was %v\n", user.Email, userFromResponse.Email)
	}
	if userFromResponse.Password != user.Password {
		t.Errorf("Expected password field be equal %v but was %v\n", user.Password, userFromResponse.Password)
	}
}

func TestUpdateUser(t *testing.T) {}

// func TestDeleteUser(t *testing.T)
