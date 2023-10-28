package core

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	dbClient "github.com/vincentcent6/product-catalogue/pkg/database"
)

var (
	db *sqlx.DB

	ErrDataNotFound = errors.New("Data not found")
)

// NewProduct will return Product object
func NewProduct() Product {
	err := prepare()
	if err != nil {
		fmt.Printf("New Product - failed to prepare: %s", err.Error())
	}
	return &product{}
}

func prepare() (err error) {
	db, err = dbClient.GetConnection()
	if err != nil {
		return err
	}
	return nil
}
