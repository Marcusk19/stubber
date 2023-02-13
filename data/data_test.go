package data

import (
	"testing"
)

func TestConnect(t *testing.T) {
	db := Connect()
	err := db.Ping()
	if err != nil {
		print(err.Error())
		t.Errorf("Connection to database failed")
	}

}
