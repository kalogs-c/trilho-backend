package models

import (
	"errors"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/kalogsc/ego/utils"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint32    `json:"id" gorm:"primary_key;unique;auto_increment"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	LastName  string    `json:"last_name" gorm:"size:100;not null"`
	Email     string    `json:"email" gorm:"size:100;not null;unique"`
	Password  string    `json:"password" gorm:"size:100;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) prepare() {
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.CreatedAt = time.Now()
}

func (u *User) Validate(mode string) error {
	if strings.ToLower(mode) == "add" {
		if u.Name == "" {
			return errors.New("field 'Name' is required")
		}
		if u.LastName == "" {
			return errors.New("field 'LastName' is required")
		}
	}
	if u.Email == "" {
		return errors.New("field 'Email' is required")
	}
	if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("invalid email")
	}
	if u.Password == "" {
		return errors.New("field 'Password' is required")
	}
	if !utils.ValidatePassword(u.Password) {
		return errors.New("field 'Password' must contain at least 6 digits")
	}
	return nil
}

func (u *User) Save(db *gorm.DB) error {
	var err error

	u.prepare()

	err = u.Validate("add")
	if err != nil {
		return err
	}

	err = u.hashPassword()
	if err != nil {
		return err
	}

	err = db.Debug().Create(&u).Error
	if err != nil {
		if strings.Contains(err.Error(), "Error 1062:") {
			return errors.New("email already taken")
		}
		fmt.Println(err)
		return err
	}
	return nil
}

func (u *User) Delete(db *gorm.DB) error {
	err := db.Debug().Delete(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}
	err := db.Debug().Model(&u).Limit(100).Find(&users).Error
	if err != nil {
		return &users, err
	}
	return &users, nil
}

func (u *User) UpdateUser(db *gorm.DB) error {
	err := u.hashPassword()
	if err != nil {
		return err
	}
	db = db.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"name":      u.Name,
			"last_name": u.LastName,
			"email":     u.Email,
			"password":  u.Password,
		},
	)
	if db.Error != nil {
		return db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CollectUserData(db *gorm.DB) error {
	err := db.Debug().Model(&User{}).Where("id = ?", u.ID).Take(&u).Error
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("user not found")
	} else if err != nil {
		return err
	}

	return nil
}
