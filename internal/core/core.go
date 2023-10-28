package core

import (
	"errors"
	"fmt"

	esLib "github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	dbClient "github.com/vincentcent6/product-catalogue/pkg/database"
	esClient "github.com/vincentcent6/product-catalogue/pkg/elasticsearch"
)

var (
	db *sqlx.DB
	es *esLib.TypedClient

	ErrDataNotFound = errors.New("Data not found")
)

const (
	productsIndex = "products"
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

	es, err = esClient.GetESClient()
	if err != nil {
		return err
	}

	return nil
}
