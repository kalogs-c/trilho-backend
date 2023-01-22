package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Transaction struct {
	ID        uint32    `json:"id" gorm:"primary_key;unique;auto_increment"`
	Amount    float64   `json:"amount" gorm:"not null"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	OwnerId   uint32    `json:"owner_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

func (t *Transaction) prepare() {
	t.Name = html.EscapeString(strings.TrimSpace(t.Name))
	t.CreatedAt = time.Now()
}

func (t *Transaction) validate() error {
	if t.Name == "" {
		return errors.New("field 'Name' is required")
	}
	if t.OwnerId == 0 {
		return errors.New("field 'OwnerId' cannot be zero")
	}
	return nil
}

func (t *Transaction) Save(db *gorm.DB) error {
	t.prepare()

	var err error

	err = t.validate()
	if err != nil {
		return err
	}
	err = db.Debug().Create(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) Delete(db *gorm.DB) error {
	err := db.Debug().Delete(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) UpdateTransaction(db *gorm.DB) error {
	t.prepare()
	err := t.validate()
	if err != nil {
		return err
	}
	db = db.Debug().Model(&Transaction{}).Where("id = ?", t.ID).Take(&Transaction{}).UpdateColumns(
		map[string]interface{}{
			"name":   t.Name,
			"amount": t.Amount,
		},
	)
	if db.Error != nil {
		return db.Error
	}

	err = db.Debug().Model(&Transaction{}).Where("id = ?", t.ID).Take(&t).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *Transaction) CollectUserTransactions(db *gorm.DB) (*[]Transaction, error) {
	var err error

	Transactions := []Transaction{}

	err = db.Debug().Model(&Transaction{}).Where("owner_id = ?", t.OwnerId).Find(&Transactions).Error
	if err != nil {
		return &[]Transaction{}, err
	}

	return &Transactions, nil
}

func (t *Transaction) CollectTransactionData(db *gorm.DB) error {
	err := db.Debug().Model(&Transaction{}).Where("id = ?", t.ID).Take(&t).Error
	if gorm.IsRecordNotFoundError(err) {
		return errors.New("transaction not found")
	} else if err != nil {
		return err
	}

	return nil
}