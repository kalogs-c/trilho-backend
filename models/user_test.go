package models_test

import (
	"testing"

	"github.com/kalogsc/ego/models"
	"github.com/kalogsc/ego/seed"
)

func TestCreateUser(t *testing.T) {
	user := &models.User{
		Name:     "Carlos",
		LastName: "Henrique",
		Email:    "carlos@email.com",
		Password: "potato123",
	}
	userCopy := *user
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
	if userCopy.Password == user.Password {
		t.Errorf("Expected password field be equal %v\n", user.Password)
	}

	user = &userCopy
}

func TestGetUserData(t *testing.T) {
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{
		{
			Name:     "Zoro",
			LastName: "Roronoa",
			Email:    "zoro@email.com",
			Password: "sword123",
		},
	})
	if err != nil {
		t.Errorf("failed to seed user: %v\n", err)
		return
	}
	user := &models.User{
		Email:    "zoro@email.com",
	}
	err = user.CollectUserData(serverInstance.DB)
	userCopy := *user
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
	user := &models.User{}
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
	user := &models.User{
		ID:       1,
		Name:     "Monkey",
		LastName: "D. Luffy",
		Email:    "luffy@email.com",
		Password: "gomu123",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{
		user,
	})
	if err != nil {
		t.Errorf("failed to seed user: %v\n", err)
		return
	}
	userCopy := *user

	user.Name = "João"
	user.LastName = "Caique"
	user.Email = "jao@email.com"
	user.Password = "Tomato123"

	err = user.UpdateUser(serverInstance.DB)
	if err != nil {
		t.Errorf("error updating user %v", err)
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
	user := &models.User{
		ID:       1,
		Name:     "Monkey",
		LastName: "D. Luffy",
		Email:    "luffy@email.com",
		Password: "gomu123",
	}
	err := seed.LoadCustomUsers(serverInstance.DB, &[]*models.User{
		user,
	})
	if err != nil {
		t.Errorf("failed to seed user: %v\n", err)
		return
	}
	err = user.Delete(serverInstance.DB)
	if err != nil {
		t.Errorf("error deleting user: %v\n", err)
		return
	}
}
