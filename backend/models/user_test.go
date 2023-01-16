package models_test

import (
	"testing"

	"github.com/kalogsc/ego/models"
)

var user models.User = models.User{
	Name:     "Carlos",
	LastName: "Henrique",
	Email:    "carlos@email.com",
	Password: "potato123",
}

func TestCreateUser(t *testing.T) {
	userInstance, err := user.Save(serverInstance.DB)
	if err != nil {
		t.Errorf("this is the error creating an user: %v\n", err)
		return
	}

	if userInstance.Name != user.Name {
		t.Errorf("Expected name field be equal %v\n", user.Name)
	}
	if userInstance.LastName != user.LastName {
		t.Errorf("Expected lastname field be equal %v\n", user.LastName)
	}
	if userInstance.Email != user.Email {
		t.Errorf("Expected email field be equal %v\n", user.Email)
	}
	if userInstance.Password != user.Password {
		t.Errorf("Expected password field be equal %v\n", user.Password)
	}
}

func TestGetUserData(t *testing.T) {
	userInstance, err := user.CollectUserData(serverInstance.DB)
	if err != nil {
		t.Errorf("error getting user data: %v\n", err)
		return
	}

	if userInstance.Name != user.Name {
		t.Errorf("Expected name field be equal %v\n", user.Name)
	}
	if userInstance.LastName != user.LastName {
		t.Errorf("Expected lastname field be equal %v\n", user.LastName)
	}
	if userInstance.Email != user.Email {
		t.Errorf("Expected email field be equal %v\n", user.Email)
	}
	if userInstance.Password != user.Password {
		t.Errorf("Expected password field be equal %v\n", user.Password)
	}
}

func TestDeleteUser(t *testing.T) {
	err := user.Delete(serverInstance.DB)
	if err != nil {
		t.Errorf("error deleting user: %v\n", err)
		return
	}
}
