package models_test

import (
	"testing"

	"github.com/kalogsc/ego/models"
)

var transaction *models.Transaction = &models.Transaction{
	Amount:  250.0,
	Name:    "Bread",
	OwnerId: 1,
}

func TestCreateTransaction(t *testing.T) {
	transactionCopy := *transaction
	err := transaction.Save(serverInstance.DB)
	if err != nil {
		t.Errorf("error creating an transaction: %v\n", err)
		return
	}

	if transactionCopy.Name != transaction.Name {
		t.Errorf("Expected name field be equal %v\n", transaction.Name)
		return
	}
	if transactionCopy.Amount != transaction.Amount {
		t.Errorf("Expected lastname field be equal %f\n", transaction.Amount)
		return
	}
	if transactionCopy.OwnerId != transaction.OwnerId {
		t.Errorf("Expected email field be equal %d\n", transaction.OwnerId)
		return
	}

	transaction = &transactionCopy
}

func TestCollectUserTransactions(t *testing.T) {
	transactionList, err := transaction.CollectUserTransactions(serverInstance.DB)
	if err != nil {
		t.Errorf("error listing the transactions: %v\n", err)
		return
	}

	if len(*transactionList) == 0 {
		t.Errorf("empty list")
		return
	}
}

func TestUpdateTransaction(t *testing.T) {
	transactionCopy := *transaction

	transaction.Amount = 500
	transaction.Name = "Potato"

	err := transaction.UpdateTransaction(serverInstance.DB)
	if err != nil {
		t.Errorf("error updating transaction")
		return
	}

	if transactionCopy.Name == transaction.Name {
		t.Errorf("Expected name field not be equal %v\n", transaction.Name)
		return
	}
	if transactionCopy.Amount == transaction.Amount {
		t.Errorf("Expected lastname field not be equal %f\n", transaction.Amount)
		return
	}
	if transactionCopy.OwnerId != transaction.OwnerId {
		t.Errorf("Expected email field be equal %d\n", transaction.OwnerId)
		return
	}
}

func TestDeleteTransaction(t *testing.T) {
	err := transaction.Delete(serverInstance.DB)
	if err != nil {
		t.Errorf("error deleting Transaction: %v\n", err)
		return
	}
}
