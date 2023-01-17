package models_test

import (
	"testing"

	"github.com/kalogsc/ego/models"
)

var user *models.User = &models.User{
	Name:     "Carlos",
	LastName: "Henrique",
	Email:    "carlos@email.com",
	Password: "potato123",
}

func TestCreateUser(t *testing.T) {
	userCopy := *user;
	err := user.Save(serverInstance.DB)
	if err != nil {
		t.Errorf("this is the error creating an user: %v\n", err)
		return
	}

	if userCopy.Name != user.Name {
		t.Errorf("Expected name field be equal %v\n", user.Name)
	}
	if userCopy.LastName != user.LastName {
		t.Errorf("Expected lastname field be equal %v\n", user.LastName)
	}
	if userCopy.Email != user.Email {
		t.Errorf("Expected email field be equal %v\n", user.Email)
	}
	if userCopy.Password != user.Password {
		t.Errorf("Expected password field be equal %v\n", user.Password)
	}

	user = &userCopy
}

func TestGetUserData(t *testing.T) {
	userCopy := *user
	err := user.CollectUserData(serverInstance.DB)
	if err != nil {
		t.Errorf("error getting user data: %v\n", err)
		return
	}

	if userCopy.Name != user.Name {
		t.Errorf("Expected name field be equal %v\n", user.Name)
		return
	}
	if userCopy.LastName != user.LastName {
		t.Errorf("Expected lastname field be equal %v\n", user.LastName)
		return
	}
	if userCopy.Email != user.Email {
		t.Errorf("Expected email field be equal %v\n", user.Email)
		return
	}
	if userCopy.Password != user.Password {
		t.Errorf("Expected password field be equal %v\n", user.Password)
		return
	}
}

func TestListUsers(t *testing.T) {
	users, err := user.FindAllUsers(serverInstance.DB)
	if err != nil {
		t.Errorf("Failed list users %v", err)
		return
	}

	if len(*users) < 1 {
		t.Errorf("Failed list users %v", err)
		return
	}
}

func TestUpdateUser(t *testing.T) {
	userCopy := *user

	user.Name = "João"
	user.LastName = "Caique"
	user.Email = "jao@email.com"
	user.Password = "Tomato123"

	err := user.UpdateUser(serverInstance.DB)
	if err != nil {
		t.Errorf("error updating user")
		return
	}

	if userCopy.Name == user.Name {
		t.Errorf("Expected name field not be equal %v\n", user.Name)
		return
	}
	if userCopy.LastName == user.LastName {
		t.Errorf("Expected lastname field not be equal %v\n", user.LastName)
		return
	}
	if userCopy.Email == user.Email {
		t.Errorf("Expected email field not be equal %v\n", user.Email)
		return
	}
}

func TestDeleteUser(t *testing.T) {
	err := user.Delete(serverInstance.DB)
	if err != nil {
		t.Errorf("error deleting user: %v\n", err)
		return
	}
}
