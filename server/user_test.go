package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kalogsc/trilho/models"
	"github.com/kalogsc/trilho/seed"
)

func TestCreateUser(t *testing.T) {
	users := &[]struct {
		user                 *models.User
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		{
			user: &models.User{
				Name:     "Carlos",
				LastName: "Henrique",
				Username:    "carlos",
				Password: "carlinhos!",
			},
			expectedStatusCode:   http.StatusCreated,
			expectedErrorMessage: "",
		},
		{
			user: &models.User{
				Name:     "Other",
				LastName: "Carlos",
				Username:    "carlos",
				Password: "notcarlinhos!",
			},
			expectedStatusCode:   http.StatusConflict,
			expectedErrorMessage: "username already taken",
		},
		{
			user: &models.User{
				Name:     "",
				LastName: "Santana",
				Username:    "maria",
				Password: "maria!",
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedErrorMessage: "field 'Name' is required",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "",
				Username:    "maria",
				Password: "maria!",
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedErrorMessage: "field 'LastName' is required",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "Santana",
				Username:    "",
				Password: "maria!",
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedErrorMessage: "field 'Username' is required",
		},
		{
			user: &models.User{
				Name:     "Maria",
				LastName: "Santana",
				Username:    "maria",
				Password: "",
			},
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedErrorMessage: "field 'Password' is required",
		},
	}

	for _, v := range *users {
		userJsonFormat, err := json.Marshal(v.user)
		if err != nil {
			t.Errorf("error marshalling json: %v", err)
			return
		}

		req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(userJsonFormat))
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
		if rr.Code != v.expectedStatusCode {
			t.Errorf("error: expected status code be %v but was %v", v.expectedStatusCode, rr.Code)
			return
		}
		if v.expectedStatusCode == http.StatusCreated {
			if responseMap["name"] != v.user.Name {
				t.Errorf("error: expected field 'name' be %v but was %v", responseMap["name"], v.user.Name)
				return
			}
			if responseMap["last_name"] != v.user.LastName {
				t.Errorf("error: expected field 'last_name' be %v but was %v", responseMap["last_name"], v.user.LastName)
				return
			}
			if responseMap["username"] != v.user.Username {
				t.Errorf("error: expected field 'username' be %v but was %v", responseMap["username"], v.user.Username)
				return
			}
			if responseMap["password"] == v.user.Password {
				t.Error("error: expected field 'password' be different")
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
		Username:    "really",
		Password: "trust_me157",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{user})
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
		return
	}

	if userFromResponse.Name != user.Name {
		t.Errorf("Expected name field be equal %v but was %v\n", user.Name, userFromResponse.Name)
		return
	}
	if userFromResponse.LastName != user.LastName {
		t.Errorf("Expected lastname field be equal %v but was %v\n", user.LastName, userFromResponse.LastName)
		return
	}
	if userFromResponse.Username != user.Username {
		t.Errorf("Expected username field be equal %v but was %v\n", user.Username, userFromResponse.Username)
		return
	}
	if userFromResponse.Password != user.Password {
		t.Errorf("Expected password field be equal %v but was %v\n", user.Password, userFromResponse.Password)
		return
	}
}

func TestUpdateUser(t *testing.T) {
	var userForAuth models.User

	customUsers := &[]*models.User{
		{
			Name:     "Babidi",
			LastName: "From DBZ",
			Username:    "babidi",
			Password: "majinboo",
		},
		{
			Name:     "Goku",
			LastName: "Kakarot",
			Username:    "goku",
			Password: "kamehameha",
		},
	}
	seed.LoadCustomUsers(serverInstance.DB, customUsers)

	for _, user := range *customUsers {
		if user.Username == "goku" {
			userForAuth = *user
			userForAuth.Password = "kamehameha"
		}
	}

	token, err := serverInstance.SignIn(&userForAuth)
	if err != nil {
		t.Errorf("error: failed while sign in %e", err)
		return
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	usersUpdate := &[]struct {
		updateUser           models.User
		expectedStatusCode   int
		tokenGiven           string
		expectedErrorMessage string
	}{
		{
			updateUser: models.User{
				ID:       userForAuth.ID,
				Name:     "Vegeta",
				LastName: "King",
				Username:    "vegeta",
				Password: "notkamehameha",
			},
			expectedStatusCode:   http.StatusOK,
			tokenGiven:           tokenString,
			expectedErrorMessage: "",
		},
		{
			// When no token was passed
			updateUser: models.User{
				ID:       userForAuth.ID,
				Name:     "Vegeta",
				LastName: "King",
				Username:    "vegeta",
				Password: "notkamehameha",
			},
			expectedStatusCode:   http.StatusUnauthorized,
			tokenGiven:           "",
			expectedErrorMessage: "Unauthorized",
		},
		{
			// When incorrect token was passed
			updateUser: models.User{
				ID:       userForAuth.ID,
				Name:     "Vegeta",
				LastName: "King",
				Username:    "vegeta",
				Password: "notkamehameha",
			},
			expectedStatusCode:   http.StatusUnauthorized,
			tokenGiven:           "incorrect token",
			expectedErrorMessage: "Unauthorized",
		},
		{
			updateUser: models.User{
				ID:       userForAuth.ID,
				Name:     "Babidi",
				LastName: "From DBZ",
				Username:    "babidi",
				Password: "majinboo",
			},
			expectedStatusCode:   http.StatusConflict,
			tokenGiven:           tokenString,
			expectedErrorMessage: "username already taken",
		},
		{
			// When user 2 is using user 1 token
			updateUser: models.User{
				ID:       1,
				Name:     "Vegeta",
				LastName: "King",
				Username:    "vegeta",
				Password: "notkamehameha",
			},
			tokenGiven:           tokenString,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
	}

	for _, v := range *usersUpdate {
		jsonUpdateUser, err := json.Marshal(v.updateUser)
		if err != nil {
			t.Errorf("error marshalling json: %v", err)
			return
		}

		req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUpdateUser))
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
			return
		}
		req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(v.updateUser.ID))})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.UpdateUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
			return
		}
		if rr.Code != v.expectedStatusCode {
			t.Errorf("error: expected status code be %v but was %v\n Response Body: %v", v.expectedStatusCode, rr.Code, responseMap)
			return
		}
		if v.expectedStatusCode == http.StatusOK {
			if responseMap["name"] != v.updateUser.Name {
				t.Errorf("Expected name field be equal %v but was %v\n", v.updateUser.Name, responseMap["name"])
				return
			}
			if responseMap["last_name"] != v.updateUser.LastName {
				t.Errorf("Expected name field be equal %v but was %v\n", v.updateUser.LastName, responseMap["last_name"])
				return
			}
			if responseMap["username"] != v.updateUser.Username {
				t.Errorf("Expected name field be equal %v but was %v\n", v.updateUser.Name, responseMap["username"])
				return
			}
		} else if v.expectedErrorMessage != "" && v.expectedStatusCode == rr.Code {
			if responseMap["error"] != v.expectedErrorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.expectedErrorMessage, responseMap["error"])
				return
			}
		}
	}
}

func TestDeleteUser(t *testing.T) {
	var userForAuth models.User

	customUsers := &[]*models.User{
		{
			Name:     "Majin",
			LastName: "Boo",
			Username:    "boo",
			Password: "chocolate",
		},
		{
			Name:     "Naruto",
			LastName: "Uzumaki",
			Username:    "naruto",
			Password: "ninetail",
		},
	}
	seed.LoadCustomUsers(serverInstance.DB, customUsers)

	for _, user := range *customUsers {
		if user.Username == "boo" {
			userForAuth = *user
			userForAuth.Password = "chocolate"
		}
	}

	token, err := serverInstance.SignIn(&userForAuth)
	if err != nil {
		t.Errorf("error: failed while sign in %e", err)
		return
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	usersToDelete := &[]struct {
		id                   string
		tokenGiven           string
		expectedStatusCode   int
		expectedErrorMessage string
	}{
		{
			id:                   strconv.Itoa(int(userForAuth.ID)),
			tokenGiven:           tokenString,
			expectedStatusCode:   http.StatusOK,
			expectedErrorMessage: "",
		},
		{
			id:                   strconv.Itoa(int(userForAuth.ID)),
			tokenGiven:           "Bearer wrong token",
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
		{
			id:                   "Cannot convert this id to an integer",
			tokenGiven:           tokenString,
			expectedStatusCode:   http.StatusUnprocessableEntity,
			expectedErrorMessage: "cannot convert this id to an integer",
		},
		{
			id:                   strconv.Itoa(2),
			tokenGiven:           tokenString,
			expectedStatusCode:   http.StatusUnauthorized,
			expectedErrorMessage: "Unauthorized",
		},
	}

	for _, v := range *usersToDelete {
		req, err := http.NewRequest("DELETE", "/user", nil)
		if err != nil {
			t.Errorf("This is the error: %v\n", err)
			return
		}
		req = mux.SetURLVars(req, map[string]string{"id": v.id})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(serverInstance.DeleteUser)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert %v to json: %v", rr.Body.String(), err.Error())
			return
		}
		if rr.Code != v.expectedStatusCode {
			t.Errorf("error: expected status code be %v but was %v", v.expectedStatusCode, rr.Code)
			return
		}
		if v.expectedErrorMessage != "" && v.expectedStatusCode == rr.Code {
			if responseMap["error"] != v.expectedErrorMessage {
				t.Errorf("invalid error, expected to be %v but was %v", v.expectedErrorMessage, responseMap["error"])
				return
			}
		}
	}
}
