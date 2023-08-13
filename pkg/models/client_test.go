package models

import (
	"fmt"
	"gomessage/pkg/utils/database"
	"testing"
)

func init() {
	database.DB.Init("sqlite3")
}

func TestGetClientById(t *testing.T) {
	result, _ := GetClientById(1)
	fmt.Println(result)
}
