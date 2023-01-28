package utils_test

import (
	"testing"

	"github.com/kalogsc/trilho/utils"
)

func TestValidateName(t *testing.T) {
	names := []struct {
		name   string
		expect bool
	}{
		{name: "Carlos", expect: true},
		{name: "Carlos Henrique", expect: false},
		{name: "1Carlos", expect: false},
		{name: "#arlos", expect: false},
	}

	for _, v := range names {
		result := utils.ValidateName(v.name)
		if result != v.expect {
			t.Errorf("expect the name %v return %v but returned %v", v.name, v.expect, result)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	passwords := []struct {
		password string
		expect   bool
	}{
		{password: "morethan5", expect: true},
		{password: "less", expect: false},
		{password: "Have Space", expect: false},
	}

	for _, v := range passwords {
		result := utils.ValidatePassword(v.password)
		if result != v.expect {
			t.Errorf("expect the password %v return %v but returned %v", v.password, v.expect, result)
		}
	}
}
