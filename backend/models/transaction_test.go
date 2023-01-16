package models_test

import (
	"testing"

	"github.com/kalogsc/ego/models"
)

var transaction models.Transaction = models.Transaction{
	Amount:  250.0,
	Name:    "Bread",
	OwnerId: user.ID,
}

func TestCreateTransaction(t *testing.T) {
	transactionInstance, err := transaction.Save(serverInstance.DB)
	if err != nil {
		t.Errorf("this is the error creating an transaction: %v\n", err)
		return
	}

	if transactionInstance.Name != transaction.Name {
		t.Errorf("Expected name field be equal %v\n", transaction.Name)
	}
	if transactionInstance.Amount != transaction.Amount {
		t.Errorf("Expected lastname field be equal %f\n", transaction.Amount)
	}
	if transactionInstance.OwnerId != transaction.OwnerId {
		t.Errorf("Expected email field be equal %d\n", transaction.OwnerId)
	}
}

func TestCollectUserTransactions(t *testing.T) {
	transactionList, err := transaction.CollectUserTransactions(serverInstance.DB)
	if err != nil {
		t.Errorf("this is the error listing the jokes: %v\n", err)
		return
	}

	if len(*transactionList) == 0 {
		t.Errorf("empty list")
	}
}

func TestDeleteTransaction(t *testing.T) {
	err := transaction.Delete(serverInstance.DB)
	if err != nil {
		t.Errorf("error deleting Transaction: %v\n", err)
		return
	}
}
