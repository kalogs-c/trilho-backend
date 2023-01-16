package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID        uint32    `json:"id" gorm:"primary_key;unique;auto_increment"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	LastName  string    `json:"lastname" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:100;not null"`
	Password  string    `json:"Password" gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (t *User) prepare() {
	t.Name = html.EscapeString(strings.TrimSpace(t.Name))
	t.LastName = html.EscapeString(strings.TrimSpace(t.LastName))
	t.Email = html.EscapeString(strings.TrimSpace(t.Email))
	t.Password = html.EscapeString(strings.TrimSpace(t.Password))
	t.CreatedAt = time.Now()
}

func (t *User) validate() error {
	if t.Name == "" {
		return errors.New("field 'Name' is required")
	}
	if t.LastName == "" {
		return errors.New("field 'LastName' is required")
	}
	if t.Email == "" {
		return errors.New("field 'Email' is required")
	}
	if t.Password == "" {
		return errors.New("field 'Password' is required")
	}
	return nil
}

func (u *User) Save(db *gorm.DB) (*User, error) {
	u.prepare()

	var err error

	err = u.validate()
	if err != nil {
		return &User{}, err
	}
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) Delete(db *gorm.DB) error {
	err := db.Debug().Delete(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CollectUserData(db *gorm.DB) (*User, error) {
	var err error

	UserFromDB := User{}

	err = db.Debug().Model(&User{}).Where("id = ?", u.ID).Find(&UserFromDB).Error
	if err != nil {
		return &User{}, err
	}

	return &UserFromDB, nil
}
